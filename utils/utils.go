package utils

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"strconv"

	"context"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")

	if awsAccessKeyID == "" || awsSecretAccessKey == "" || awsRegion == "" {
		log.Fatal("AWS credentials or region not set in .env file")
	}

	staticCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
        awsAccessKeyID,
        awsSecretAccessKey,
        ""))

    // load AWS config with static credentials
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion(awsRegion),
        config.WithCredentialsProvider(staticCreds),
    )
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

    // create an S3 client
    s3Client = s3.NewFromConfig(cfg)

}

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

func notInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return false
		}
	}
	return true
}
