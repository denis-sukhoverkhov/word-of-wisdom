package repository

import "github.com/denis-sukhoverkhov/word-of-wisdom/internal/models"

type QuoteRepository interface {
	GetRandomQuote() (models.Quote, error)
}

type GlobalRepository struct {
	QuoteRepo QuoteRepository
}

func NewGlobalRepository(quoteRepo QuoteRepository) *GlobalRepository {
	return &GlobalRepository{
		QuoteRepo: quoteRepo,
	}
}
