package models

import (
	"time"

	"gorm.io/gorm"
)

type Punishment struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TargetID    uint           `json:"target_id"`
	Target      User           `json:"target" gorm:"foreignKey:TargetID"`
	AssignedBy  uint           `json:"assigned_by"`
	Assigner    User           `json:"assigner" gorm:"foreignKey:AssignedBy"`
	Type        string         `json:"type"` // "timeout", "demotion", "extra_tasks"
	Description string         `json:"description"`
	Status      string         `json:"status" gorm:"default:'active'"` // "active", "completed", "cancelled"
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	ExpiresAt   *time.Time     `json:"expires_at"`
}

// Constantes para tipos de castigo
const (
	PunishmentTypeTimeout    = "timeout"
	PunishmentTypeDemotion   = "demotion"
	PunishmentTypeExtraTasks = "extra_tasks"
	PunishmentTypeReward     = "reward"
)

// Constantes para status de castigo
const (
	PunishmentStatusActive    = "active"
	PunishmentStatusCompleted = "completed"
	PunishmentStatusCancelled = "cancelled"
)
