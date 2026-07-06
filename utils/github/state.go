package github

import (
	"crypto/rand"
	"math/big"
)

func GenerateState() (string, error) {
	result := make([]byte, StateLength)
	alphabetLength := big.NewInt(int64(len(StateAlphabet)))

	for index := range result {
		randomIndex, randomError := rand.Int(rand.Reader, alphabetLength)
		if randomError != nil {
			return "", randomError
		}
		result[index] = StateAlphabet[randomIndex.Int64()]
	}

	return string(result), nil
}
