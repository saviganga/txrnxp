package utils

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type UpdateEntityResponse[T any] struct {
	Data           T           `json:"-"`
	SerializedData interface{} `json:"data"`
}

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

	// create an S3 clients
	s3Client = s3.NewFromConfig(cfg)
	presignClient = s3.NewPresignClient(s3Client)

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

func (r *GenericDBStruct[T]) UpdateEntity(c *fiber.Ctx, table string, id string) (UpdateEntityResponse[T], error) {

	// validate db model
	var model T
	var modelField reflect.Value
	var modelValue reflect.Value

	authenticated_user := c.Locals("user").(jwt.MapClaims)

	table = strings.ToLower(table)
	validTables := []string{"xuser", "business", "event", "event_ticket", "user_ticket"}
	if notInList(table, validTables) {
		return UpdateEntityResponse[T]{}, errors.New("invalid model")
	}

	fields := c.Locals("body")

	val := reflect.ValueOf(fields)
	typ := reflect.TypeOf(fields)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if table == "xuser" {

		if authenticated_user["id"].(string) != id {
			return UpdateEntityResponse[T]{}, errors.New("you do not have permission to perform this action")
		}

		if err := r.db.First(&model, "id = ?", id).Error; err != nil {
			return UpdateEntityResponse[T]{}, errors.New("model not found")
		}

	} else if table == "business" {

		if err := r.db.First(&model, "id = ?", id).Error; err != nil {
			return UpdateEntityResponse[T]{}, errors.New("model not found")
		}

		modelValue = reflect.ValueOf(&model).Elem()

		// validate business owner
		UserIdField := modelValue.FieldByName("UserId")
		if !UserIdField.IsValid() {
			return UpdateEntityResponse[T]{}, errors.New("model does not have UserId field")
		}
		UserIdValue := UserIdField.Interface().(uuid.UUID).String()

		if authenticated_user["id"].(string) != UserIdValue {
			return UpdateEntityResponse[T]{}, errors.New("you do not have permission to perform this action")
		}

	} else if table == "event" {

		if err := r.db.First(&model, "id = ?", id).Error; err != nil {
			return UpdateEntityResponse[T]{}, errors.New("model not found")
		}

	} else {
		return UpdateEntityResponse[T]{}, errors.New("chill out for the workflow to be ready my guy")
	}

	for i := 0; i < val.NumField(); i++ {

		// get the field name
		fieldName := typ.Field(i).Name

		// get the field value
		fieldValue := val.Field(i).Interface()

		// get the values in the fields
		modelValue = reflect.ValueOf(&model).Elem()
		modelField = modelValue.FieldByName(fieldName)
		if !modelField.IsValid() {
			return UpdateEntityResponse[T]{}, errors.New("model does not have " + fieldName + "field")
		}

		// update the field value
		modelField.SetString(fieldValue.(string))

	}

	if err := r.db.Save(&model).Error; err != nil {
		return UpdateEntityResponse[T]{}, errors.New("unable to update model")
	}

	return UpdateEntityResponse[T]{
		Data: model,
	}, nil

}
