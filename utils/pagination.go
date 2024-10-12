package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GenericDBStruct[T any] struct {
	db *gorm.DB
}

type PaginationResponse[T any] struct {
	Status         string      `json:"status"`
	Message        string      `json:"message"`
	Type           string      `json:"type"`
	Total          int64       `json:"total"`
	Page           int         `json:"page"`
	Limit          int         `json:"size"`
	NextPage       *int        `json:"nextPage,omitempty"`
	PreviousPage   *int        `json:"previousPage,omitempty"`
	TotalPages     int         `json:"totalPages"`
	Data           []T         `json:"-"`
	SerializedData interface{} `json:"data"`
}

func isStringField(field string) bool {
	stringFields := map[string]bool{
		"name":         true,
		"email":        true,
		"reference":    true,
		"country":      true,
		"first_name":   true,
		"last_name":    true,
		"event_type":   true,
		"description":  true,
		"address":      true,
		"category":     true,
		"duration":     true,
		"entry_type":   true,
		"ticket_type":  true,
		"is_validated": true,
	}

	return stringFields[field]
}

func NewGenericDB[T any](db *gorm.DB) *GenericDBStruct[T] {
	return &GenericDBStruct[T]{db: db}
}

func DBpaginate(db *gorm.DB, limit, page int) *gorm.DB {
	offset := (page - 1) * limit
	return db.Offset(offset).Limit(limit)
}

func paginate(limit, page int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Order("created_at desc").Offset(offset).Limit(limit)
	}
}

func ValidateRequestLimitAndPage(c *fiber.Ctx) error {

	// extract query parameters for limit and page.
	limitStr := c.Query("size", "10")
	pageStr := c.Query("page", "1")

	// convert query parameters to integers.
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	c.Locals("size", limit)
	c.Locals("page", page)

	return c.Next()
}

func (r *GenericDBStruct[T]) GetPagedAndFiltered(limit, page int, filters map[string]interface{}, preloads []string, joins []string) (PaginationResponse[T], error) {

	var results []T
	var total int64

	// start building the query
	query := r.db.Model(new(T))

	// Apply joins if any are provided
	if len(joins) > 0 {
		for _, join := range joins {
			if join != "" {
				query = query.Joins(join)
			}
		}
	}

	// apply preloads if any are provided
	if len(preloads) > 0 {
		for _, preload := range preloads {
			if preload != "" {
				query = query.Preload(preload)
			}
		}
	}

	// apply filters to the query
	for key, value := range filters {
		normalizedKey := strings.Replace(key, "__", ".", -1)
		parts := strings.Split(normalizedKey, ".")
		if isStringField(parts[len(parts)-1]) {
			query = query.Where(fmt.Sprintf("LOWER(%s) LIKE ?", normalizedKey), fmt.Sprintf("%%%s%%", value))
		} else {
			query = query.Where(fmt.Sprintf("%s = ?", normalizedKey), fmt.Sprintf("%s", value))
		}

	}

	// get total count of records after applying filters
	err := query.Count(&total).Error
	if err != nil {
		return PaginationResponse[T]{}, err
	}

	// Apply pagination scope and get the results with filters applied
	err = query.Scopes(paginate(limit, page)).Find(&results).Error
	if err != nil {
		return PaginationResponse[T]{}, err
	}

	// calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// determine next and previous pages
	var nextPage *int
	var previousPage *int
	if page < totalPages {
		n := page + 1
		nextPage = &n
	}
	if page > 1 {
		p := page - 1
		previousPage = &p
	}

	// return paginated response with metadata
	return PaginationResponse[T]{
		Total:        total,
		Page:         page,
		Limit:        limit,
		NextPage:     nextPage,
		PreviousPage: previousPage,
		TotalPages:   totalPages,
		Data:         results,
	}, nil
}
