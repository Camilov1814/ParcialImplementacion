package dto

import "time"

// Request para crear castigo/recompensa
type CreatePunishmentRequest struct {
	TargetID    uint   `json:"target_id" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description" binding:"required"`
	ExpiresAt   string `json:"expires_at,omitempty"` // RFC3339 format
}

// Response de castigo
type PunishmentResponse struct {
	ID          uint       `json:"id"`
	Target      UserInfo   `json:"target"`
	Assigner    UserInfo   `json:"assigner"`
	Type        string     `json:"type"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// Item para lista de castigos
type PunishmentListItem struct {
	ID           uint       `json:"id"`
	TargetName   string     `json:"target_name"`
	AssignerName string     `json:"assigner_name"`
	Type         string     `json:"type"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
}

// Request para actualizar castigo
type UpdatePunishmentRequest struct {
	Status      string `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
	ExpiresAt   string `json:"expires_at,omitempty"`
}
