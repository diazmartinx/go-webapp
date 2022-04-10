package helpers

import (
	"crypto/rand"
	"math/big"
)

var chars = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"q", "w", "e", "r", "t", "y", "u", "i", "o", "p",
	"a", "s", "d", "f", "g", "h", "j", "k", "l",
	"z", "x", "c", "v", "b", "n", "m",
	"Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P",
	"A", "S", "D", "F", "G", "H", "J", "K", "L",
	"Z", "X", "C", "V", "B", "N", "M",
}

func randint() int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(62))
	if err != nil {
		panic(err)
	}
	return nBig.Int64()
}

func GenerateRandomString(length int) string {
	var s string
	for i := 0; i < length; i++ {
		s += chars[randint()]
	}
	return s
}
