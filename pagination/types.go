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
