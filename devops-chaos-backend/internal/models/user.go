package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Role      string         `json:"role" gorm:"not null"`           // "andrei", "daemon", "network_admin"
	Status    string         `json:"status" gorm:"default:'active'"` // "active", "captured", "punished"
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	Reports     []Report     `json:"reports,omitempty" gorm:"foreignKey:AuthorID"`
	Statistics  *Statistic   `json:"statistics,omitempty" gorm:"foreignKey:UserID"`
	Punishments []Punishment `json:"punishments,omitempty" gorm:"foreignKey:TargetID"`
}

// Constantes para roles
const (
	RoleAndrei       = "andrei"
	RoleDaemon       = "daemon"
	RoleNetworkAdmin = "network_admin"
)

// Constantes para status
const (
	StatusActive   = "active"
	StatusCaptured = "captured"
	StatusPunished = "punished"
)
