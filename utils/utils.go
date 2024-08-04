package utils

import (
	"errors"
	"math/rand"
	"strconv"

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

func CreateEventReference() string {
	return GenerateRandomString((6))
}

func ConvertStringToFloat(amount string) (float64, error) {
	// convert amount to float
	amount_float, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0.0, errors.New("error converting event ticket price")
	}

	return amount_float, nil

}
