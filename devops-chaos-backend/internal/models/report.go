package models

import (
	"time"

	"gorm.io/gorm"
)

type Report struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Type        string         `json:"type"`      // "resistance", "capture", "anonymous"
	AuthorID    *uint          `json:"author_id"` // nil for anonymous reports
	Author      *User          `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Status      string         `json:"status" gorm:"default:'pending'"` // "pending", "approved", "rejected"
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Constantes para tipos de reporte
const (
	ReportTypeResistance = "resistance"
	ReportTypeCapture    = "capture"
	ReportTypeAnonymous  = "anonymous"
)

// Constantes para status de reporte
const (
	ReportStatusPending  = "pending"
	ReportStatusApproved = "approved"
	ReportStatusRejected = "rejected"
)
