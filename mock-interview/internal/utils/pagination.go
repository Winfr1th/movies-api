package utils

import (
	"net/http"
	"strconv"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// PagedResponse represents a paginated API response
type PagedResponse struct {
	Data     interface{} `json:"data"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Total    int         `json:"total"`
}

// ParsePaginationParams extracts and validates pagination parameters from request
func ParsePaginationParams(r *http.Request) (page, pageSize int, err error) {
	// Parse page parameter
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		page = DefaultPage
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return 0, 0, &PaginationError{Code: "INVALID_PAGE", Message: "page must be a positive integer"}
		}
	}

	// Parse page_size parameter
	pageSizeStr := r.URL.Query().Get("page_size")
	if pageSizeStr == "" {
		pageSize = DefaultPageSize
	} else {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			return 0, 0, &PaginationError{Code: "INVALID_PAGE_SIZE", Message: "page_size must be a positive integer"}
		}
	}

	// Validate page_size doesn't exceed maximum
	if err := ValidatePageSize(pageSize); err != nil {
		return 0, 0, err
	}

	return page, pageSize, nil
}

// ValidatePageSize ensures page_size doesn't exceed the maximum allowed
func ValidatePageSize(pageSize int) error {
	if pageSize > MaxPageSize {
		return &PaginationError{Code: "PAGE_SIZE_TOO_LARGE", Message: "page_size exceeds maximum allowed value of 100"}
	}
	return nil
}

// PaginationError represents a pagination validation error
type PaginationError struct {
	Code    string
	Message string
}

func (e *PaginationError) Error() string {
	return e.Message
}

// CreatePagedResponse creates a paginated response structure
func CreatePagedResponse(data interface{}, total, page, pageSize int) PagedResponse {
	return PagedResponse{
		Data:     data,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
}

// CalculateOffset calculates the SQL OFFSET value for pagination
func CalculateOffset(page, pageSize int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * pageSize
}
