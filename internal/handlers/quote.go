package handlers

import (
	"net"

	"github.com/denis-sukhoverkhov/word-of-wisdom/internal/repository"
	"go.uber.org/zap"
)

func HandleQuote(conn net.Conn, repo *repository.GlobalRepository, logger *zap.Logger) {
	quote, err := repo.QuoteRepo.GetRandomQuote()
	if err != nil {
		_, writeErr := conn.Write([]byte("Failed to get quote: " + err.Error() + "\n"))
		if writeErr != nil {
			logger.Error("Error writing to connection", zap.Error(writeErr))
		}
		return
	}

	_, writeErr := conn.Write([]byte("Here is your quote: '" + quote.Text + "' - " + quote.Author + "\n"))
	if writeErr != nil {
		logger.Error("Error writing to connection", zap.Error(writeErr))
	}
}
