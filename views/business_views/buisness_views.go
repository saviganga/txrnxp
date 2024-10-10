package business_views

import (
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/business_serializers"
	"txrnxp/utils"
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
		return utils.BadRequestResponse(c, err.Error())
	}

	return utils.CreatedResponse(c, business, "Successfully created business")
}

func GetBusiness(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	privilege := authenticated_user["privilege"]

	businessRepo := utils.NewGenericDB[models.Business](db)

	if privilege == "ADMIN" {
		// define filters based on query parameters
		filters := c.Locals("filters").(map[string]interface{})
		limit := c.Locals("size").(int)
		page := c.Locals("page").(int)
		joins := []string{"LEFT JOIN xusers AS u ON businesses.user_id = u.id"}
        preloads := []string{"User"}
		businesses, err := businessRepo.GetPagedAndFiltered(limit, page, filters, preloads, joins)
        if err != nil {
            return utils.BadRequestResponse(c, "Unable to get businesses")
        }
		serialized_businesses := business_serializers.SerializeReadBusiness(businesses.Data)
		businesses.SerializedData = serialized_businesses
        businesses.Status = "Success"
        businesses.Message = "Successfully fetched businesses"
        businesses.Type = "OK"
		return utils.PaginatedSuccessResponse(c, businesses, "success")
	} else {
		joins := []string{"LEFT JOIN xusers AS u ON businesses.user_id = u.id"}
        preloads := []string{"User"}
		limit := c.Locals("size").(int)
		page := c.Locals("page").(int)
		filters := make(map[string]interface{})
		filters["u__id"] = authenticated_user["id"]
		businesses, err := businessRepo.GetPagedAndFiltered(limit, page, filters, preloads, joins)
        if err != nil {
            return utils.BadRequestResponse(c, "Unable to get businesses")
        }
		serialized_businesses := business_serializers.SerializeReadBusiness(businesses.Data)
		businesses.SerializedData = serialized_businesses
        businesses.Status = "Success"
        businesses.Message = "Successfully fetched businesses"
        businesses.Type = "OK"
		return utils.PaginatedSuccessResponse(c, businesses, "success")
	}

}
