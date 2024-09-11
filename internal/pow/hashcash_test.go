package pow

import (
	"bytes"
	"testing"
)

func TestGenerateChallenge(t *testing.T) {
	pow := NewProofOfWork(20)

	challenge1 := pow.GenerateChallenge()
	challenge2 := pow.GenerateChallenge()

	if len(challenge1) != NonceSize || len(challenge2) != NonceSize {
		t.Fatalf("Expected challenge size to be %d, got %d and %d", NonceSize, len(challenge1), len(challenge2))
	}

	if bytes.Equal(challenge1, challenge2) {
		t.Fatal("Expected different challenges, but got the same")
	}
}

func TestValidateSolution_Valid(t *testing.T) {
	pow := NewProofOfWork(20)
	challenge := pow.GenerateChallenge()
	solution := pow.Solve(challenge)

	if !pow.ValidateSolution(challenge, solution) {
		t.Fatal("Expected solution to be valid, but got invalid")
	}
}

func TestValidateSolution_Invalid(t *testing.T) {
	pow := NewProofOfWork(20)
	challenge := pow.GenerateChallenge()

	// Manually create an invalid solution by using a wrong nonce
	invalidSolution := make([]byte, NonceSize)
	for i := range invalidSolution {
		invalidSolution[i] = 255 // invalid random data
	}

	if pow.ValidateSolution(challenge, invalidSolution) {
		t.Fatal("Expected solution to be invalid, but got valid")
	}
}

func TestSolve(t *testing.T) {
	pow := NewProofOfWork(20)
	challenge := pow.GenerateChallenge()

	solution := pow.Solve(challenge)

	if !pow.ValidateSolution(challenge, solution) {
		t.Fatal("Expected solution to be valid, but got invalid")
	}
}
