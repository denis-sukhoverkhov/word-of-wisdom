package repository

import "testing"

func TestGetRandomQuote_ReturnsValidQuote(t *testing.T) {
	repo := NewStaticQuoteRepository()

	quote, err := repo.GetRandomQuote()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if quote.Text == "" {
		t.Errorf("Expected non-empty quote text, got empty")
	}

	if quote.Author == "" {
		t.Errorf("Expected non-empty quote author, got empty")
	}
}
