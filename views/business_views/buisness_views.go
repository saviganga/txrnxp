package business_views

import (
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/business_serializers"
	"txrnxp/serializers/user_serializers"
	"txrnxp/utils"
	"txrnxp/utils/business_utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to the United States of Ganga business biaaattccchhhhhh!")
}

func CreateBusiness(c *fiber.Ctx) error {

	serialized_business := new(business_serializers.ReadBusinessSerializer)
	serialized_user := new(user_serializers.UserSerializer)
	business, err := business_utils.CreateBusiness(c)

	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
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

	serialized_business.Id = business.Id
	serialized_business.User = *serialized_user
	serialized_business.Reference = business.Reference
	serialized_business.Name = business.Name
	serialized_business.Country = business.Country
	serialized_business.CreatedAt = business.CreatedAt
	serialized_business.UpdatedAt = business.UpdatedAt

	return utils.CreatedResponse(c, serialized_business, "Successfully created business")
}

func GetBusiness(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	businesses := []models.Business{}
	serialized_business := new(business_serializers.ReadBusinessSerializer)
	serialized_businesses := []business_serializers.ReadBusinessSerializer{}
	serialized_user := new(user_serializers.UserSerializer)
	privilege := authenticated_user["privilege"]
	if privilege == "ADMIN" {
		db.Model(&models.Business{}).Joins("User").Order("created_at desc").Find(&businesses).Order("businesses.created_at DESC")
	} else {
		db.Model(&models.Business{}).Joins("User").Order("created_at desc").First(&businesses, "businesses.user_id = ?", authenticated_user["id"]).Order("businesses.created_at DESC")
	}
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

		serialized_business.Id = business.Id
		serialized_business.User = *serialized_user
		serialized_business.Reference = business.Reference
		serialized_business.Name = business.Name
		serialized_business.Country = business.Country
		serialized_business.CreatedAt = business.CreatedAt
		serialized_business.UpdatedAt = business.UpdatedAt

		serialized_businesses = append(serialized_businesses, *serialized_business)
	}
	return utils.SuccessResponse(c, serialized_businesses, "Successfully fetched businesses")

}
