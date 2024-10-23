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

type CreateBusinessMemberSerializer struct {
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserName    string `json:"username"`
	PhoneNumber string `json:"phone"`
	Password    string `json:"password"`
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

type ReadBusinessMemberSerializer struct {
	Id        uuid.UUID                       `json:"id" validate:"required"`
	User      user_serializers.UserSerializer `json:"user" validate:"required"`
	Business  ReadCreateBusinessSerializer          `json:"business" validate:"required"`
	CreatedAt time.Time                       `json:"created_at" validate:"required"`
	UpdatedAt time.Time                       `json:"updated_at" validate:"required"`
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

func SerializeReadBusinessMembers(business_members []models.BusinessMember, c *fiber.Ctx) ([]ReadBusinessMemberSerializer, error) {

	serialized_user := new(user_serializers.UserSerializer)
	serialized_business := new(ReadCreateBusinessSerializer)
	serialized_business_member := new(ReadBusinessMemberSerializer)
	serialized_business_members := []ReadBusinessMemberSerializer{}
	db := initialisers.ConnectDb().Db
	businessRepo := utils.NewGenericDB[models.BusinessMember](db)
	var imageUrl string
	var err error

	for _, business_member := range business_members {

		serialized_user.Id = business_member.User.Id
		serialized_user.Email = business_member.User.Email
		serialized_user.UserName = business_member.User.UserName
		serialized_user.FirstName = business_member.User.FirstName
		serialized_user.LastName = business_member.User.LastName
		serialized_user.PhoneNumber = business_member.User.PhoneNumber
		serialized_user.IsActive = business_member.User.IsActive
		serialized_user.IsBusiness = business_member.User.IsBusiness
		serialized_user.LastLogin = business_member.User.LastLogin
		serialized_user.CreatedAt = business_member.User.CreatedAt
		serialized_user.UpdatedAt = business_member.User.UpdatedAt

		if business_member.Business.Image != "" {
			imageUrl, err = businessRepo.GetSignedUrl(c, "business", business_member.Business.Id.String())
			if err != nil {
				return nil, errors.New(err.Error())
			}
		} else {
			imageUrl = ""
		}

		serialized_business.Id = business_member.Business.Id
		serialized_business.Reference = business_member.Business.Reference
		serialized_business.Name = business_member.Business.Name
		serialized_business.Image = imageUrl
		serialized_business.Country = business_member.Business.Country
		serialized_business.CreatedAt = business_member.Business.CreatedAt
		serialized_business.UpdatedAt = business_member.Business.UpdatedAt

		serialized_business_member.Id = business_member.Id
		serialized_business_member.User = *serialized_user
		serialized_business_member.Business = *serialized_business
		serialized_business_member.CreatedAt = business_member.CreatedAt
		serialized_business_member.UpdatedAt = business_member.UpdatedAt

		serialized_business_members = append(serialized_business_members, *serialized_business_member)
	}

	return serialized_business_members, nil

}
