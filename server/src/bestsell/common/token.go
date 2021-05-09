package common

import (
	"math/rand"
)

var tokenChars = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*")

func GenToken(r *rand.Rand, length int) []byte {
	chars := len(tokenChars)
	token := make([]byte, length)
	for i := 0; i < len(token); i++ {
		token[i] = tokenChars[r.Intn(chars)]
	}
	return token
}

var tokenChars2 = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenToken2(r *rand.Rand, length int) []byte {
	chars := len(tokenChars2)
	token := make([]byte, length)
	for i := 0; i < len(token); i++ {
		token[i] = tokenChars2[r.Intn(chars)]
	}
	return token
}
