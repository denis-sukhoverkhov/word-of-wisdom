package pow

type PoWAlgorithm interface {
	GenerateChallenge() []byte
	ValidateSolution(challenge []byte, nonce []byte) bool
	Solve(challenge []byte) []byte
}
