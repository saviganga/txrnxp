package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ValidateRequestFilters(getTableName func() string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		filters := make(map[string]interface{})
		validFields := []string{""}

		// get the table name dynamically
		table := getTableName()

		// Convert table name to lowercase
		table = strings.ToLower(table)

		// validate the model table exists
		validTables := []string{"xuser", "wallets", "business", "event", "wallet_tx", "admin_wallet_tx"}
		if notInList(table, validTables) {
			c.Locals("filters", filters)
			return c.Next()
		}

		// validate the table fields
		if table == "xuser" {
			validFields = []string{"email", "first_name", "last_name", "is_active", "is_verified", "is_business"}
		} else if table == "wallets" {
			validFields = []string{"u__email", "u__first_name", "u__last_name"}
		} else if table == "business" {
			validFields = []string{"u__email", "u__first_name", "u__last_name", "name", "reference"}
		} else if table == "event" {
			validFields = []string{"name", "reference", "event_type", "is_business", "reference", "description", "address", "category", "duration"}
		} else if table == "wallet_tx" {
			validFields = []string{"u__email", "u__first_name", "u__last_name", "reference", "entry_type", "description"}
		} else if table == "admin_wallet_tx" {
			validFields = []string{"reference", "entry_type", "description"}
		}

		c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {

			keyStr := string(key)
			valueStr := string(value)

			// skip empty values
			if len(valueStr) == 0 {
				return
			}

			// skip invalid fields
			if notInList(keyStr, validFields) {
				return
			}

			// sdd valid filters to the filters map
			filters[keyStr] = strings.ToLower(valueStr)
		})

		// store filters in locals
		c.Locals("filters", filters)

		return c.Next()
	}
}
