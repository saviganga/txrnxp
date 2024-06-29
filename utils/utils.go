package utils

import (
	"errors"
	"math/rand"

	"github.com/google/uuid"
)

func GenerateRandomString(length int) string {
	characters := "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	s := make([]byte, length)
	for i := range s {
		s[i] = characters[rand.Intn(len(characters))]
	}

	return string(s)
}

func ConvertStringToUUID(s string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, errors.New("invalid parsed id")
	}
	return parsedUUID, nil
}
