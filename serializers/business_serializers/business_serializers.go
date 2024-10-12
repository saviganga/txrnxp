package business_serializers

import (
	"errors"
	"time"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/user_serializers"
	"txrnxp/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UpdateBusinessSerializer struct {
	Name string `json:"name"`
}

type ReadBusinessSerializer struct {
	Id        uuid.UUID                       `json:"id" validate:"required"`
	User      user_serializers.UserSerializer `json:"user" validate:"required"`
	Reference string                          `json:"reference" validate:"required"`
	Name      string                          `json:"name" validate:"required"`
	Image     string                          `json:"image" validate:"required"`
	Country   string                          `json:"country" validate:"required"`
	CreatedAt time.Time                       `json:"created_at" validate:"required"`
	UpdatedAt time.Time                       `json:"updated_at" validate:"required"`
}

type ReadCreateBusinessSerializer struct {
	Id        uuid.UUID `json:"id" validate:"required"`
	Reference string    `json:"reference" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Image     string    `json:"image" validate:"required"`
	Country   string    `json:"country" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}

func SerializeCreateBusiness(business models.Business, c *fiber.Ctx) (*ReadCreateBusinessSerializer, error) {

	db := initialisers.ConnectDb().Db
	businessRepo := utils.NewGenericDB[models.Business](db)
	serialized_business := new(ReadCreateBusinessSerializer)
	var imageUrl string
	var err error

	if business.Image != "" {
		imageUrl, err = businessRepo.GetSignedUrl(c, "business", business.Id.String())
		if err != nil {
			return nil, errors.New(err.Error())
		}
	} else {
		imageUrl = ""
	}

	serialized_business.Id = business.Id
	serialized_business.Reference = business.Reference
	serialized_business.Name = business.Name
	serialized_business.Image = imageUrl
	serialized_business.Country = business.Country
	serialized_business.CreatedAt = business.CreatedAt
	serialized_business.UpdatedAt = business.UpdatedAt

	return serialized_business, nil

}

func SerializeReadBusiness(businesses []models.Business, c *fiber.Ctx) ([]ReadBusinessSerializer, error) {

	serialized_user := new(user_serializers.UserSerializer)
	serialized_business := new(ReadBusinessSerializer)
	serialized_businesses := []ReadBusinessSerializer{}
	db := initialisers.ConnectDb().Db
	businessRepo := utils.NewGenericDB[models.Business](db)
	var imageUrl string
	var err error

	for _, business := range businesses {

		serialized_user.Id = business.User.Id
		serialized_user.Email = business.User.Email
		serialized_user.UserName = business.User.UserName
		serialized_user.FirstName = business.User.FirstName
		serialized_user.LastName = business.User.LastName
		serialized_user.PhoneNumber = business.User.PhoneNumber
		serialized_user.IsActive = business.User.IsActive
		serialized_user.IsBusiness = business.User.IsBusiness
		serialized_user.LastLogin = business.User.LastLogin
		serialized_user.CreatedAt = business.User.CreatedAt
		serialized_user.UpdatedAt = business.User.UpdatedAt

		if business.Image != "" {
			imageUrl, err = businessRepo.GetSignedUrl(c, "business", business.Id.String())
			if err != nil {
				return nil, errors.New(err.Error())
			}
		} else {
			imageUrl = ""
		}

		serialized_business.Id = business.Id
		serialized_business.User = *serialized_user
		serialized_business.Reference = business.Reference
		serialized_business.Name = business.Name
		serialized_business.Image = imageUrl
		serialized_business.Country = business.Country
		serialized_business.CreatedAt = business.CreatedAt
		serialized_business.UpdatedAt = business.UpdatedAt

		serialized_businesses = append(serialized_businesses, *serialized_business)
	}

	return serialized_businesses, nil

}
