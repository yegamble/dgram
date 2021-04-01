package wallet

import (
	"crypto/rand"
	"math/big"
)

type IotaWallet struct {
	seed       string `json:"seed"`
	privateKey string `json:"privateKey"`
}

var version = "undefined"

const letters = "9ABCDEFGHIJKLMNOPQRSTUVWXYZ" //pool of letters to generate IOTA seed

func GenerateNewWallet() string {
	newWallet := new(IotaWallet)
	newWallet.seed, _ = GenerateRandomSeed()
	//encodedSeed, err := encodeArgon(newWallet.seed) //save seed instead

	//if err != nil {
	//	return err.Error()
	//}

	return newWallet.seed
}

func generateRandomInts(n int) ([]int64, error) {
	ints := make([]int64, n)

	for i := range ints {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(27))

		if err != nil {
			return nil, err
		}

		ints[i] = randomInt.Int64()
	}

	return ints, nil
}

func intToCharByte(i int64) byte {
	return byte(letters[i])
}

func GenerateRandomSeed() (string, error) {
	ints, err := generateRandomInts(81)

	if err != nil {
		return "", err
	}

	token := make([]byte, 81)

	for i, x := range ints {
		token[i] = intToCharByte(x)
	}

	return string(token), nil
}
