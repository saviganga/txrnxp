package business_utils

import (
	"errors"
	"fmt"
	"strings"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/business_serializers"
	"txrnxp/serializers/user_serializers"

	"txrnxp/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func CreateBusiness(c *fiber.Ctx) (*business_serializers.ReadCreateBusinessSerializer, error) {

	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	business := new(models.Business)
	privilege := authenticated_user["privilege"]

	if privilege == "ADMIN" {
		return nil, errors.New("oops! this feature is not available for admins")
	}

	parsedUUID, err := utils.ConvertStringToUUID(authenticated_user["id"].(string))

	if err != nil {
		return nil, errors.New("invalid parsed id")
	}

	business.UserId = parsedUUID
	err = c.BodyParser(business)
	if err != nil {
		return nil, errors.New("invalid request body")
	}
	err = db.Create(&business).Error

	if err != nil {
		return nil, errors.New(err.Error())
	}

	business_member_query := models.BusinessMember{UserId: business.UserId, BusinessId: business.Id}
	dbError := db.Create(&business_member_query).Error
	if dbError != nil {
		return nil, errors.New("oops! error creating business member")
	}

	serialized_business, err := business_serializers.SerializeCreateBusiness(*business, c)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return serialized_business, nil
}

func GetBusinessById(c *fiber.Ctx) error {
	// authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	business := models.Business{}
	// privilege := authenticated_user["privilege"]
	err := db.First(&business, "id = ?", c.Params("id")).Error
	if err != nil {
		return utils.BadRequestResponse(c, "Unable to get user")
	}
	// if privilege != "ADMIN" && authenticated_user["id"].(string) != business.UserId.String() {
	// 	return utils.BadRequestResponse(c, "You do not have permission to view this resource")
	// }

	serialized_business, err := business_serializers.SerializeCreateBusiness(business, c)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, serialized_business, "success")

}

func UpdateBusiness(c *fiber.Ctx) error {

	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	businessRepo := utils.NewGenericDB[models.Business](db)
	privilege := authenticated_user["privilege"].(string)
	business_id := c.Params("id")

	if strings.ToUpper(privilege) == "ADMIN" {
		return utils.BadRequestResponse(c, "this feature is not available for admins")
	}

	business, err := businessRepo.UpdateEntity(c, "business", business_id)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	serialized_business, err := business_serializers.SerializeCreateBusiness(business.Data, c)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, serialized_business, "success")

}

func UploadBusinessImage(c *fiber.Ctx) error {

	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	businessRepo := utils.NewGenericDB[models.Business](db)
	privilege := authenticated_user["privilege"].(string)
	business_id := c.Params("id")

	if strings.ToUpper(privilege) == "ADMIN" {
		return utils.BadRequestResponse(c, "this feature is not available for admins")
	}

	business, err := businessRepo.UploadImage(c, "business", business_id)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	serialized_business, err := business_serializers.SerializeCreateBusiness(business.Data, c)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	business.SerializedData = serialized_business
	business.Status = "Success"
	business.Message = "Successfully uploaded business image"
	business.Type = "OK"
	return utils.SuccessResponse(c, serialized_business, "Successfully uploaded business image")

}

func RemoveBusinessKeys(filters map[string]interface{}) {
	for key := range filters {
		if strings.Contains(key, "business__") {
			delete(filters, key)
		}
	}
}

func CreateBusinessMember(c *fiber.Ctx) error {

	db := initialisers.ConnectDb().Db
	business := models.Business{}
	business_reference := c.Get("Business")
	user_request := business_serializers.CreateBusinessMemberSerializer{}
	user_request.Password = "casa1234"
	user := models.Xuser{}

	err := c.BodyParser(&user_request)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid request data")
	}

	// validate the business
	err = db.Model(&models.Business{}).Find(&business, "id = ?", c.Params("id")).Error
	if err != nil {
		return utils.BadRequestResponse(c, fmt.Sprintf("oops! unable to fetch business - reference: %s", business_reference))
	}

	// validate that the user email exists
	serialized_user := user_serializers.UserSerializer{}
	err = db.Find(&user, "email = ?", user_request.Email).Error
	if user.Id == uuid.Nil || err != nil {
		user_query := models.Xuser{FirstName: user_request.FirstName, LastName: user_request.LastName, UserName: user_request.UserName, Email: user_request.Email, PhoneNumber: user_request.PhoneNumber, Password: user_request.Password}
		dbError := db.Create(&user_query).Error
		if dbError != nil {
			return errors.New("oops! error creating user wallet")
		}

		business_member_query := models.BusinessMember{UserId: user_query.Id, BusinessId: business.Id}
		dbError = db.Create(&business_member_query).Error
		if dbError != nil {
			return utils.BadRequestResponse(c, "oops! error creating business member")
		}

	} else {
		serialized_user = user_serializers.SerializeUserSerializer(user)
		business_member_query := models.BusinessMember{UserId: serialized_user.Id, BusinessId: business.Id}
		dbError := db.Create(&business_member_query).Error
		if dbError != nil {
			return utils.BadRequestResponse(c, "oops! error creating business member")
		}
	}


	return utils.CreatedResponse(c, user_request, "Successfully created business member")

}
