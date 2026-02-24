package utils

import (
	"github.com/gofiber/fiber/v2"
)

// PaginationParams represents pagination query parameters
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// PaginatedResponse wraps paginated data with metadata
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// ParsePagination extracts and validates pagination params from query string
func ParsePagination(c *fiber.Ctx) PaginationParams {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return PaginationParams{Page: page, PageSize: pageSize}
}

// Offset calculates the SQL offset from page and page_size
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// NewPaginatedResponse creates a paginated response from data and total count
func NewPaginatedResponse(data interface{}, total int64, params PaginationParams) PaginatedResponse {
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}
}
