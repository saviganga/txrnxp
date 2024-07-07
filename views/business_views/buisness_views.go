package business_views

import (
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/utils/business_utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to the United States of Ganga business biaaattccchhhhhh!")
}

func CreateBusiness(c *fiber.Ctx) error {

	business, err := business_utils.CreateBusiness(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(business)
}

func GetBusiness(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	businesses := []models.Business{}
	privilege := authenticated_user["privilege"]
	if privilege == "ADMIN" {
		db.Model(&models.Business{}).Joins("User").Find(&businesses).Order("businesses.created_at DESC")
	} else {
		db.Model(&models.Business{}).Joins("User").First(&businesses, "businesses.user_id = ?", authenticated_user["id"]).Order("businesses.created_at DESC")
	}
	return c.Status(200).JSON(businesses)

}
