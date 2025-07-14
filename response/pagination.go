// ==================== response/pagination.go ====================
package response

import "math"

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
	Offset     int `json:"offset"`
}

// DefaultQueryParams for parsing pagination from request
// Uses omitempty without min validation to allow 0 values that will be converted to defaults
type DefaultQueryParams struct {
	Page      int    `form:"page" json:"page" binding:"omitempty"`
	Limit     int    `form:"limit" json:"limit" binding:"omitempty"`
	SortBy    string `form:"sortBy" json:"sortBy"`
	SortOrder string `form:"sortOrder" json:"sortOrder"`
}

// FlexibleQueryParams for more complex scenarios
// Allows 0 values and converts them to defaults automatically
type FlexibleQueryParams struct {
	Page       int      `form:"page" json:"page" binding:"omitempty"`
	Limit      int      `form:"limit" json:"limit" binding:"omitempty"`
	Search     string   `form:"search" json:"search" binding:"omitempty,max=100"`
	CategoryID string   `form:"categoryId" json:"categoryId" binding:"omitempty,uuid"`
	LocationID string   `form:"locationId" json:"locationId" binding:"omitempty,uuid"`
	Condition  string   `form:"condition" json:"condition" binding:"omitempty,oneof=new good fair poor"`
	MinPrice   *float64 `form:"minPrice" json:"minPrice" binding:"omitempty,min=0"`
	MaxPrice   *float64 `form:"maxPrice" json:"maxPrice" binding:"omitempty,min=0"`
	SortBy     string   `form:"sortBy" json:"sortBy" binding:"omitempty,oneof=name price createdAt purchaseDate"`
	SortOrder  string   `form:"sortOrder" json:"sortOrder" binding:"omitempty,oneof=asc desc"`
}

// BuildPagination creates pagination from parameters with smart defaults
func BuildPagination(page, limit, total int) *Pagination {
	// Smart defaults - handle all edge cases
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Prevent abuse
	}
	if total < 0 {
		total = 0
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages < 1 && total > 0 {
		totalPages = 1
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	return &Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Offset:     offset,
	}
}

// Optional helper methods
func (p *Pagination) HasNext() bool {
	return p.Page < p.TotalPages
}

func (p *Pagination) HasPrev() bool {
	return p.Page > 1
}

func (p *Pagination) NextPage() int {
	if p.HasNext() {
		return p.Page + 1
	}
	return p.Page
}

func (p *Pagination) PrevPage() int {
	if p.HasPrev() {
		return p.Page - 1
	}
	return p.Page
}

// SetDefaults applies default values to query params with smart validation
func (q *DefaultQueryParams) SetDefaults() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100 // Prevent abuse
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	if q.SortOrder == "" {
		q.SortOrder = "desc"
	}
	// Normalize sort order
	if q.SortOrder != "asc" && q.SortOrder != "desc" {
		q.SortOrder = "desc"
	}
}

// SetDefaults for FlexibleQueryParams
func (q *FlexibleQueryParams) SetDefaults() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	if q.SortOrder == "" {
		q.SortOrder = "desc"
	}
	if q.SortOrder != "asc" && q.SortOrder != "desc" {
		q.SortOrder = "desc"
	}
}

// Validate checks business logic constraints (not binding constraints)
func (q *FlexibleQueryParams) Validate() error {
	if q.MinPrice != nil && q.MaxPrice != nil && *q.MinPrice > *q.MaxPrice {
		return BadRequest("Min price cannot be greater than max price")
	}
	return nil
}

// GetOffset calculates offset for database queries
func (q *DefaultQueryParams) GetOffset() int {
	return (q.Page - 1) * q.Limit
}

// GetOffset for FlexibleQueryParams
func (q *FlexibleQueryParams) GetOffset() int {
	return (q.Page - 1) * q.Limit
}
