package pagination

import "math"

type Meta struct {
	Total      int64 `json:"total"`
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	TotalPages int64 `json:"total_pages"`
}

func NewMeta(total, page, limit int64) Meta {
	totalPages := int64(math.Ceil(float64(total) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	return Meta{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}
