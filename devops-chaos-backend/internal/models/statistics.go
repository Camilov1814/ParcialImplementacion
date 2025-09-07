package models

import (
	"time"

	"gorm.io/gorm"
)

type Statistic struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"unique"`
	User          User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CapturesCount int            `json:"captures_count" gorm:"default:0"`
	ReportsCount  int            `json:"reports_count" gorm:"default:0"`
	Ranking       int            `json:"ranking" gorm:"default:0"`
	Points        int            `json:"points" gorm:"default:0"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
