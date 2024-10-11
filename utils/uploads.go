package utils

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type UploadImageSerializer struct {
	Image string `json:"image"`
}

type UploadImageResponse[T any] struct {
	Status         string      `json:"status"`
	Message        string      `json:"message"`
	Type           string      `json:"type"`
	Data           T           `json:"-"`
	SerializedData interface{} `json:"data"`
}

var s3Client *s3.Client
var presignClient *s3.PresignClient

func (r *GenericDBStruct[T]) UploadImage(c *fiber.Ctx, table string) (UploadImageResponse[T], error) {

	table = strings.ToLower(table)
	validTables := []string{"xuser", "business", "event", "event_ticket", "user_ticket"}
	if notInList(table, validTables) {
		return UploadImageResponse[T]{}, errors.New("invalid model")
	}

	authenticated_user := c.Locals("user").(jwt.MapClaims)

	if table == "xuser" {

		// validate db model
		var model T
		if err := r.db.First(&model, "id = ?", authenticated_user["id"].(string)).Error; err != nil {
			return UploadImageResponse[T]{}, errors.New("model not found")
		}

		modelValue := reflect.ValueOf(&model).Elem()
		imageField := modelValue.FieldByName("Image")
		if !imageField.IsValid() {
			return UploadImageResponse[T]{}, errors.New("model does not have an Image field")
		}

		// validate the request body
		body := new(UploadImageSerializer)

		if err := c.BodyParser(&body); err != nil {
			return UploadImageResponse[T]{}, err
		}

		// validate the base64 encoding
		if !strings.Contains(body.Image, "data:image") {
			return UploadImageResponse[T]{}, errors.New("invalid image format")
		}

		// get image data from base 64 string
		parts := strings.Split(body.Image, ",")
		if len(parts) < 2 {
			return UploadImageResponse[T]{}, errors.New("invalid image data")
		}

		imageData, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return UploadImageResponse[T]{}, errors.New("unable to decode image")
		}

		// generate filename and save
		emailField := modelValue.FieldByName("Email")
		if !emailField.IsValid() {
			return UploadImageResponse[T]{}, errors.New("model does not have an Email field")
		}

		emailValue := emailField.Interface().(string)
		fileName := emailValue + ".png"
		bucketName := "txrnxp"
		objectKey := "users/" + fileName

		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(objectKey),
			Body:        strings.NewReader(string(imageData)),
			ContentType: aws.String("image/png"),
		})
		if err != nil {
			return UploadImageResponse[T]{}, fmt.Errorf("failed to upload image to S3: %w", err)
		}

		s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, objectKey)

		// set the Image field with the new file path
		imageField.SetString(s3URL)

		if err := r.db.Save(&model).Error; err != nil {
			return UploadImageResponse[T]{}, errors.New("unable to update model")
		}

		return UploadImageResponse[T]{
			Data: model,
		}, nil
	} else {
		return UploadImageResponse[T]{}, errors.New("workflow not ready. BE PATIENT NIGGGGAAAAAA")
	}

}

func GeneratePresignedURL(bucketName string, objectKey string) (string, error) {
	presignedURL, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.URL, nil
}
