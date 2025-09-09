package dto

// Response para la página de resistencia (Network Admins)
type ResistancePageResponse struct {
	WelcomeMessage       string             `json:"welcome_message"`
	UserStatus           string             `json:"user_status"`
	SurvivalTips         []SurvivalTip      `json:"survival_tips"`
	ResistanceMemes      []ResistanceMeme   `json:"resistance_memes"`
	AnonymousReportCount int                `json:"anonymous_reports_sent"`
	ResistanceStats      ResistanceStats    `json:"resistance_stats"`
	EmergencyContacts    []EmergencyContact `json:"emergency_contacts"`
}

// Consejos de supervivencia
type SurvivalTip struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"` // "low", "medium", "high", "critical"
	Category    string `json:"category"` // "technical", "social", "security"
}

// Memes de resistencia
type ResistanceMeme struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	Upvotes     int    `json:"upvotes"`
}

// Estadísticas de la resistencia
type ResistanceStats struct {
	TotalNetworkAdmins int `json:"total_network_admins"`
	CapturedAdmins     int `json:"captured_admins"`
	FreeAdmins         int `json:"free_admins"`
	AnonymousReports   int `json:"anonymous_reports_today"`
}

// Contactos de emergencia
type EmergencyContact struct {
	Name        string `json:"name"`
	Role        string `json:"role"`
	Contact     string `json:"contact"`
	Description string `json:"description"`
}

// Request para reportar actividad sospechosa
type ReportSuspiciousActivityRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Location    string `json:"location,omitempty"`
	Severity    string `json:"severity" binding:"required"` // "low", "medium", "high", "critical"
}
