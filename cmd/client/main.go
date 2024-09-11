package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/denis-sukhoverkhov/word-of-wisdom/internal/pow"
	"github.com/denis-sukhoverkhov/word-of-wisdom/pkg/config"
	"go.uber.org/zap"
)

const (
	HandlerQuote byte = 0x01 // Constant for the "quote" handler
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		logger.Fatal("Error initializing logger", zap.Error(err))
	}

	configPathDir := flag.String("config_dir", "./configs", "Путь к папке с конфигами")
	flag.Parse()

	clientConfig, err := config.LoadClientConfig(*configPathDir)
	if err != nil {
		logger.Fatal("Error of loading cfg", zap.Error(err))
	}
	logger.Info("Client configuration loaded", zap.Any("config", clientConfig))

	// create a ticker to control the request rate
	ticker := time.NewTicker(time.Second / time.Duration(clientConfig.RPS))
	defer ticker.Stop()

	// restrict the number of concurrent requests
	for i := 0; i < clientConfig.TotalRequests; i++ {
		<-ticker.C

		go func(requestNum int) {
			err := sendRequest(clientConfig.ServerAddr, HandlerQuote, logger)
			if err != nil {
				logger.Error("Request failed", zap.Int("request_num", requestNum), zap.Error(err))
			} else {
				logger.Info("Request succeeded", zap.Int("request_num", requestNum))
			}
		}(i)
	}

	// wait for all requests to finish
	time.Sleep(time.Duration(clientConfig.TotalRequests/clientConfig.RPS) * time.Second)
}

func sendRequest(serverAddr string, handlerID byte, logger *zap.Logger) error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	buffer := make([]byte, 9) // 8 bytes for the challenge and 1 byte for the difficulty
	_, err = conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read challenge and difficulty: %w", err)
	}

	challenge := buffer[:8] // first 8 bytes are the challenge
	difficulty := buffer[8] // last byte is the difficulty

	powAlgo := pow.NewProofOfWork(difficulty)
	solution := powAlgo.Solve(challenge)

	var requestBuffer bytes.Buffer
	requestBuffer.Write(solution)

	// Write the handlerID (1 byte) to the buffer
	err = requestBuffer.WriteByte(handlerID)
	if err != nil {
		return fmt.Errorf("failed to write handler ID: %w", err)
	}

	// Send the binary request to the server
	_, err = conn.Write(requestBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	// Step 5: Read the response from the server
	response, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	logger.Info("Server response", zap.String("response", string(response)))
	return nil
}
