package utils


import (
	"math/rand"
)

func GenerateRandomString(length int) string {
	characters := "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	s := make([]byte, length)
	for i := range s {
		s[i] = characters[rand.Intn(len(characters))]
	}

	return string(s)
}