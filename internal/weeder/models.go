package weeder

import "time"

// Document example
type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Embedding []float32 `json:"embedding"`
	Score     float32   `json:"score"`
}

// SearchParams example
type SearchParams struct {
	Limit      int       `json:"limit"`
	DaysOffset int       `json:"days_offset"`
	Vector     []float32 `json:"vector"`
}
