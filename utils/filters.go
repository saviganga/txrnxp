package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ValidateRequestFilters(getTableName func() string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get the table name dynamically (from the closure)
		table := getTableName()

		// Convert table name to lowercase
		table = strings.ToLower(table)

		// Validate the model table exists
		validTables := []string{"xuser"}
		if notInList(table, validTables) {
			return c.Next()
		}

		// Validate the table fields
		filters := map[string]interface{}{}
		if table == "xuser" {
			validFields := []string{"email", "first_name", "last_name", "username", "is_active", "is_verified", "is_business"}
			c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {

				keyStr := string(key)
				valueStr := string(value)

				// Skip empty values
				if len(valueStr) == 0 {
					return
				}

				// Skip invalid fields
				if notInList(keyStr, validFields) {
					return
				}

				// Add valid filters to the filters map
				filters[keyStr] = valueStr
			})
		}

		// Store filters in locals
		c.Locals("filters", filters)

		return c.Next()
	}
}
