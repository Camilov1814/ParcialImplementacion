package dto

import "time"

// Item para lista de usuarios
type UserListItem struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Request para actualizar usuario
type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role,omitempty"`
	Status   string `json:"status,omitempty"`
}

// Request para capturar network admin
type CaptureNetworkAdminRequest struct {
	TargetID uint   `json:"target_id" binding:"required"`
	Notes    string `json:"notes,omitempty"`
}

// Response para captura exitosa
type CaptureResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	PointsGiven int    `json:"points_given"`
}
