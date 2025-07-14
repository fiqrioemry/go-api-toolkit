// ==================== pagination/types.go ====================
package pagination

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
	Offset     int `json:"offset"`
}

// DefaultQueryParams for parsing pagination from request
type DefaultQueryParams struct {
	Page      int    `form:"page" json:"page" binding:"omitempty"`
	Limit     int    `form:"limit" json:"limit" binding:"omitempty"`
	SortBy    string `form:"sortBy" json:"sortBy"`
	SortOrder string `form:"sortOrder" json:"sortOrder"`
}

// FlexibleQueryParams for more complex scenarios
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
