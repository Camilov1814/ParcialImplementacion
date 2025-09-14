package dto

import "time"

// Response del dashboard de Andrei
type AndreiDashboardResponse struct {
	WelcomeMessage string                 `json:"welcome_message"`
	SystemStats    SystemStats            `json:"system_stats"`
	RecentReports  []ReportListItem       `json:"recent_reports"`
	TopDaemons     []TopPerformerResponse `json:"top_daemons"`
	RecentActivity []ActivityItem         `json:"recent_activity"`
	AllCaptures    []CaptureDetailItem    `json:"all_captures"`
}

// Response del dashboard de Daemon
type DaemonDashboardResponse struct {
	WelcomeMessage    string               `json:"welcome_message"`
	UserStats         DaemonStats          `json:"user_stats"`
	ActiveMissions    []MissionItem        `json:"active_missions"`
	Leaderboard       []RankingItem        `json:"leaderboard"`
	RecentCaos        []ChaosActivityItem  `json:"recent_chaos"`
	ActivePunishments []PunishmentListItem `json:"active_punishments"`
	RecentCaptures    []CaptureListItem    `json:"recent_captures"`
}

// Estadísticas del sistema
type SystemStats struct {
	TotalUsers      int `json:"total_users"`
	TotalDaemons    int `json:"total_daemons"`
	TotalNetAdmins  int `json:"total_network_admins"`
	CapturedAdmins  int `json:"captured_admins"`
	PunishedDaemons int `json:"punished_daemons"`
	PendingReports  int `json:"pending_reports"`
	TotalReports    int `json:"total_reports"`
}

// Estadísticas del daemon
type DaemonStats struct {
	CapturesCount int    `json:"captures_count"`
	ReportsCount  int    `json:"reports_count"`
	Points        int    `json:"points"`
	Ranking       int    `json:"ranking"`
	Status        string `json:"status"`
}

// Item de actividad
type ActivityItem struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"` // "capture", "report", "punishment"
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
}

// Item de misión
type MissionItem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	Points      int    `json:"points"`
	Status      string `json:"status"`
}

// Item de actividad de caos
type ChaosActivityItem struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Timestamp   string `json:"timestamp"`
	Points      int    `json:"points"`
}

// Item de captura para el dashboard
type CaptureListItem struct {
	ID          uint      `json:"id"`
	TargetName  string    `json:"target_name"`
	CaptureDate time.Time `json:"capture_date"`
	Status      string    `json:"status"`
	Points      int       `json:"points"`
	Difficulty  string    `json:"difficulty"`
	Method      string    `json:"method"`
}

// Item de captura detallado para Andrei
type CaptureDetailItem struct {
	ID          uint      `json:"id"`
	DaemonName  string    `json:"daemon_name"`
	TargetName  string    `json:"target_name"`
	CaptureDate time.Time `json:"capture_date"`
	Status      string    `json:"status"`
	Points      int       `json:"points"`
	Difficulty  string    `json:"difficulty"`
	Method      string    `json:"method"`
}
