package pow

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"time"
)

const (
	ByteSize     = 8
	NonceSize    = 8
	MaxByteValue = 0xFF
)

type Hashcash struct {
	Difficulty uint8
	random     *rand.Rand
}

func NewProofOfWork(difficulty uint8) *Hashcash {
	source := rand.NewSource(time.Now().UnixNano())
	return &Hashcash{
		Difficulty: difficulty,
		random:     rand.New(source),
	}
}

func (pow *Hashcash) GenerateChallenge() []byte {
	nonce := pow.random.Int63()
	challenge := make([]byte, NonceSize)
	binary.BigEndian.PutUint64(challenge, uint64(nonce))
	return challenge
}

func (pow *Hashcash) ValidateSolution(challenge, solution []byte) bool {
	hash := sha256.New()
	hash.Write(challenge)
	hash.Write(solution)
	hashSum := hash.Sum(nil)

	// check if the hash has the required number of leading zeros
	leadingZeroBytes := int(pow.Difficulty / ByteSize)
	leadingZeroBits := int(pow.Difficulty % ByteSize)

	// check leading full zero bytes
	for i := 0; i < leadingZeroBytes; i++ {
		if hashSum[i] != 0 {
			return false
		}
	}

	// check the remaining bits of the next byte
	if leadingZeroBits > 0 {
		mask := byte(MaxByteValue << (ByteSize - leadingZeroBits)) // create a mask for leading bits
		if hashSum[leadingZeroBytes]&mask != 0 {
			return false
		}
	}

	return true
}

func (pow *Hashcash) Solve(challenge []byte) []byte {
	var solution []byte

	for {
		solution = make([]byte, NonceSize)
		binary.BigEndian.PutUint64(solution, uint64(pow.random.Int63())) // generate a random solution

		if pow.ValidateSolution(challenge, solution) {
			// solution found
			break
		}
	}
	return solution
}
