package pagination

import "math"

// Build creates pagination from parameters with smart defaults
func Build(page, limit, total int) *Pagination {
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

	offset := max((page-1)*limit, 0)

	return &Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Offset:     offset,
	}
}

// Quick creates pagination with automatic defaults
func Quick(params DefaultQueryParams, total int) *Pagination {
	params.SetDefaults()
	return Build(params.Page, params.Limit, total)
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

// GetOffset calculates offset for database queries
func (q *DefaultQueryParams) GetOffset() int {
	return (q.Page - 1) * q.Limit
}
