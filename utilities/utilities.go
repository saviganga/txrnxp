package utilities

import (
	"fmt"
	"txrnxp/initialisers"
	"txrnxp/models"
)

func GetEventOrganiser(organiser_id string, reference string, is_business bool) (map[string]interface{}, error) {

	db := initialisers.ConnectDb().Db
	organiser_user := []models.Xuser{}
	organiser_business := []models.Business{}
	organiser_details := make(map[string]interface{})

	if is_business {
		err := db.Model(&models.Business{}).First(&organiser_business, "id = ?", organiser_id).Error
		if err != nil {
			return nil, fmt.Errorf("oops! unable to fetch events - organiser: %s", reference)
		}
		organiser_details["name"] = organiser_business[0].Name
		organiser_details["is_business"] = true
		organiser_details["id"] = organiser_id
		organiser_details["business_user"] = organiser_business[0].UserId.String()
	} else {
		err := db.Model(&models.Xuser{}).First(&organiser_user, "id = ?", organiser_id).Error
		if err != nil {
			return nil, fmt.Errorf("oops! unable to fetch events - organiser: %s", reference)
		}
		organiser_details["name"] = organiser_user[0].UserName
		organiser_details["is_business"] = false
		organiser_details["id"] = organiser_id
		organiser_details["business_user"] = ""
	}

	return organiser_details, nil
}
