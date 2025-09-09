package dto

import "time"

// Request para crear reporte
type CreateReportRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Type        string `json:"type" binding:"required"` // "resistance", "capture", "anonymous"
}

// Response de reporte
type ReportResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Author      *UserInfo `json:"author,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Item para lista de reportes
type ReportListItem struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	AuthorID  *uint     `json:"author_id"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

// Request para actualizar reporte
type UpdateReportRequest struct {
	Status string `json:"status" binding:"required"`
}
