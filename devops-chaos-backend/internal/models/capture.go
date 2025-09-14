package models

import (
	"time"

	"gorm.io/gorm"
)

type Capture struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	DaemonID     uint           `json:"daemon_id"`
	Daemon       User           `json:"daemon" gorm:"foreignKey:DaemonID"`
	TargetID     uint           `json:"target_id"`
	Target       User           `json:"target" gorm:"foreignKey:TargetID"`
	CaptureDate  time.Time      `json:"capture_date"`
	Status       string         `json:"status" gorm:"default:'captured'"` // "captured", "released", "escaped"
	Method       string         `json:"method"`       // Método de captura usado
	Points       int            `json:"points"`       // Puntos otorgados por la captura
	Difficulty   string         `json:"difficulty"`   // "easy", "medium", "hard"
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// Constantes para status de captura
const (
	CaptureStatusCaptured = "captured"
	CaptureStatusReleased = "released"
	CaptureStatusEscaped  = "escaped"
)

// Constantes para dificultad
const (
	CaptureDifficultyEasy   = "easy"
	CaptureDifficultyMedium = "medium"
	CaptureDifficultyHard   = "hard"
)