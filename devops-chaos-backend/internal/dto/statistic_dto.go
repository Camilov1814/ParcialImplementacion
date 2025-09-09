package dto

import "time"

// Response de estadísticas de daemon
type StatisticResponse struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	Username      string    `json:"username"`
	CapturesCount int       `json:"captures_count"`
	ReportsCount  int       `json:"reports_count"`
	Ranking       int       `json:"ranking"`
	Points        int       `json:"points"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Request para actualizar estadísticas manualmente (solo Andrei)
type UpdateStatisticRequest struct {
	CapturesCount *int `json:"captures_count,omitempty"`
	ReportsCount  *int `json:"reports_count,omitempty"`
	Points        *int `json:"points,omitempty"`
}

// Ranking de daemons para leaderboard
type RankingItem struct {
	Position      int    `json:"position"`
	Username      string `json:"username"`
	CapturesCount int    `json:"captures_count"`
	ReportsCount  int    `json:"reports_count"`
	Points        int    `json:"points"`
	Status        string `json:"status"`
}

// Top performers para dashboard
type TopPerformerResponse struct {
	Username      string `json:"username"`
	Points        int    `json:"points"`
	CapturesCount int    `json:"captures_count"`
	ReportsCount  int    `json:"reports_count"`
}
