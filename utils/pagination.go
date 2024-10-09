package utils

import (
	"math"
	"strconv"

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

func (r *GenericDBStruct[T]) GetPagedAndFiltered(limit, page int) (PaginationResponse[T], error) {
	var results []T
	var total int64

	// Get total count of records
	err := r.db.Model(new(T)).Count(&total).Error
	if err != nil {
		return PaginationResponse[T]{}, err
	}

	// Apply pagination and get the results
	err = r.db.Scopes(paginate(limit, page)).Find(&results).Error
	if err != nil {
		return PaginationResponse[T]{}, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Determine next and previous pages
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

	// Return paginated response with metadata
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

func ValidateRequestLimitAndPage(c *fiber.Ctx) error {

	// Extract query parameters for limit and page.
	limitStr := c.Query("size", "10")
	pageStr := c.Query("page", "1")

	// Convert query parameters to integers.
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
