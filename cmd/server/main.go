package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/denis-sukhoverkhov/word-of-wisdom/internal/handlers"
	"github.com/denis-sukhoverkhov/word-of-wisdom/internal/pow"
	"github.com/denis-sukhoverkhov/word-of-wisdom/internal/repository"
	"github.com/denis-sukhoverkhov/word-of-wisdom/internal/server"
	"github.com/denis-sukhoverkhov/word-of-wisdom/pkg/config"
	"go.uber.org/zap"
)

func main() {

	// Initialize zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		logger.Fatal("Error initializing logger", zap.Error(err))
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("Error flushing logger", zap.Error(err))
		}
	}()

	// Flag to specify the path to the config file
	configPathDir := flag.String("config_dir", "./configs", "Path to the configuration file")
	flag.Parse()

	// Load server configuration
	serverConfig, err := config.LoadServerConfig(*configPathDir)
	if err != nil {
		logger.Fatal("Failed to load server configuration", zap.Error(err))
	}

	logger.Info("Server configuration loaded", zap.Any("config", serverConfig))

	// Create a parent context for the server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	router := handlers.NewRouter()
	router.AddRoute(handlers.HandlerQuote, handlers.HandleQuote)

	// Create the PoW server
	srv := server.NewPoWServer(
		ctx,
		serverConfig,
		pow.NewProofOfWork(serverConfig.Difficulty),
		repository.NewGlobalRepository(repository.NewStaticQuoteRepository()),
		logger,
		router,
	)

	// Create a channel to listen for interrupt or terminate signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("Server is running...")

	// Wait for termination signal
	<-stopChan
	logger.Info("Shutdown signal received")

	// Gracefully shut down the server
	if err := srv.Shutdown(); err != nil {
		logger.Error("Error during server shutdown", zap.Error(err))
	} else {
		logger.Info("Server shutdown completed")
	}
}
