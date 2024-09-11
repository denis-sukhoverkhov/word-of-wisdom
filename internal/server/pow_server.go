package server

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/denis-sukhoverkhov/word-of-wisdom/internal/pow"
	"github.com/denis-sukhoverkhov/word-of-wisdom/internal/repository"
	"github.com/denis-sukhoverkhov/word-of-wisdom/pkg/config"
	"go.uber.org/zap"
)

type PoWServer struct {
	ctx       context.Context
	cancel    context.CancelFunc
	config    *config.ServerConfig
	repo      *repository.GlobalRepository
	powAlgo   pow.PoWAlgorithm
	connQueue chan net.Conn
	wg        sync.WaitGroup
	logger    *zap.Logger
	router    Router
}

func NewPoWServer(
	ctx context.Context,
	config *config.ServerConfig,
	pow pow.PoWAlgorithm,
	repo *repository.GlobalRepository,
	logger *zap.Logger,
	router Router,
) *PoWServer {

	ctx, cancel := context.WithCancel(ctx)

	server := &PoWServer{
		ctx:       ctx,
		cancel:    cancel,
		config:    config,
		powAlgo:   pow,
		repo:      repo,
		connQueue: make(chan net.Conn, 100),
		logger:    logger,
		router:    router,
	}

	return server
}

func (s *PoWServer) Start() error {
	listener, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		s.logger.Error("Failed to start listening", zap.Error(err))
		return fmt.Errorf("failed to listen on %s: %w", s.config.Addr, err)
	}
	defer listener.Close()
	s.logger.Info("Server started", zap.String("addr", s.config.Addr))

	// run workers
	for i := 0; i < s.config.WorkerCount; i++ {
		s.wg.Add(1)
		go s.worker(i)
		s.logger.Info("Worker started", zap.Int("worker_id", i))
	}

	// Accept incoming connections
	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("Server context cancelled, stopping accepting new connections")
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				s.logger.Warn("Failed to accept connection", zap.Error(err))
				continue
			}
			s.logger.Info("Connection accepted", zap.String("remote_addr", conn.RemoteAddr().String()))
			s.connQueue <- conn
		}
	}
}

func (s *PoWServer) worker(workerID int) {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case conn, ok := <-s.connQueue:
			if !ok {
				s.logger.Info("Connection queue closed, stopping worker", zap.Int("worker_id", workerID))
				return
			}
			s.logger.Info("Worker processing connection", zap.Int("worker_id", workerID))
			s.handleConnection(conn)
		}
	}
}

func (s *PoWServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	timeout := time.Duration(s.config.PowTimeout) * time.Second
	done := make(chan struct{})
	defer close(done)

	// Step 1: Send PoW challenge to client
	challenge := s.powAlgo.GenerateChallenge()

	// Create a buffer for challenge + difficulty
	message := append(challenge, byte(s.config.Difficulty))
	_, err := conn.Write(message)
	if err != nil {
		s.logger.Error("Failed to send challenge", zap.Error(err))
		return
	}
	s.logger.Info("Challenge and difficulty sent", zap.Binary("challenge", challenge), zap.Uint8("difficulty", s.config.Difficulty))

	// Step 2: Read client solution (binary) and handler_id (as binary data)
	go func() {
		solution := make([]byte, 8) // Expecting 8 bytes for the solution
		_, err = conn.Read(solution)
		if err != nil {
			s.logger.Error("Failed to read solution", zap.Error(err))
			return
		}
		s.logger.Info("Solution received", zap.Binary("solution", solution))

		// Read handler_id as binary (assuming it's a single byte or a small number)
		handlerIDBuf := make([]byte, 1)
		_, err = conn.Read(handlerIDBuf)
		if err != nil {
			s.logger.Error("Failed to read handler ID", zap.Error(err))
			return
		}
		handlerID := handlerIDBuf[0]
		s.logger.Info("Handler ID received", zap.Uint8("handler_id", handlerID))

		// Step 3: Validate PoW solution
		if s.powAlgo.ValidateSolution(challenge, solution) {
			s.logger.Info("PoW solution validated", zap.Binary("solution", solution))

			// Step 4: Route to the appropriate handler using byte key
			handler, exists := s.router.GetRoute(handlerID)

			if !exists {
				_, err := conn.Write([]byte("Handler not found\n"))
				if err != nil {
					s.logger.Error("Failed to write to connection", zap.Error(err))
					return
				}
				s.logger.Warn("Handler not found", zap.Uint8("handler_id", handlerID))
				return
			}

			// Call the handler
			handler(conn, s.repo, s.logger)
			s.logger.Info("Handler executed", zap.Uint8("handler_id", handlerID))
		} else {
			_, err := conn.Write([]byte("Invalid PoW solution\n"))
			if err != nil {
				s.logger.Error("Failed to write to connection", zap.Error(err))
				return
			}
			s.logger.Warn("Invalid PoW solution", zap.Binary("solution", solution))
		}

		done <- struct{}{}
	}()

	// Step 6: handle timeout
	select {
	case <-done:
		// finished processing
	case <-time.After(timeout):
		s.logger.Warn("Timeout waiting for PoW solution")
		_, err := conn.Write([]byte("Timeout waiting for PoW solution\n"))
		if err != nil {
			s.logger.Error("Failed to write timeout message", zap.Error(err))
		}
	}
}

func (s *PoWServer) Shutdown() error {
	s.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.ShutdownTimeout)*time.Second)
	defer cancel()

	// Cancel the server's context to stop workers and other operations
	s.cancel()

	// Close connection queue
	close(s.connQueue)

	// Wait for workers to finish processing connections
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	// Wait for workers or context timeout
	select {
	case <-done:
		s.logger.Info("All workers have finished")
	case <-ctx.Done():
		s.logger.Info("Context timeout exceeded during shutdown")
	}

	return nil
}
