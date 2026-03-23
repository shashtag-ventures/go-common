package dto

import "time"

// PaginationParams contains the limit and offset for paginated queries.
type PaginationParams struct {
	Limit  int `json:"limit" query:"limit"`
	Offset int `json:"offset" query:"offset"`
}

// PaginatedResponse is a generic wrapper for paginated API responses.
type PaginatedResponse[T any] struct {
	Data       []T        `json:"data"`
	TotalCount int        `json:"total_count"`
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset,omitempty"`
	NextCursor *time.Time `json:"next_cursor,omitempty"`
}
