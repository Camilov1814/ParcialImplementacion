package services

import (
	"devops-chaos-backend/internal/dao"
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/models"
	"errors"
)

type ResistanceService struct {
	reportDAO *dao.ReportDAO
	userDAO   *dao.UserDAO
}

func NewResistanceService() *ResistanceService {
	return &ResistanceService{
		reportDAO: dao.NewReportDAO(),
		userDAO:   dao.NewUserDAO(),
	}
}

// Obtener página de resistencia para Network Admins
func (s *ResistanceService) GetResistancePage(userRole string) (*dto.ResistancePageResponse, error) {
	if userRole != models.RoleNetworkAdmin {
		return nil, errors.New("unauthorized: only network admins can access resistance page")
	}

	// Contar reportes anónimos de hoy (simplificado)
	anonymousReports, err := s.reportDAO.CountByStatus(models.ReportTypeAnonymous)
	if err != nil {
		anonymousReports = 0
	}

	// Estadísticas de resistencia
	resistanceStats, err := s.getResistanceStats()
	if err != nil {
		return nil, err
	}

	response := &dto.ResistancePageResponse{
		WelcomeMessage:       "Stay strong, fellow admin. The resistance continues...",
		UserStatus:           "Active Resistance Member",
		SurvivalTips:         s.getSurvivalTips(),
		ResistanceMemes:      s.getResistanceMemes(),
		AnonymousReportCount: int(anonymousReports),
		ResistanceStats:      *resistanceStats,
		EmergencyContacts:    s.getEmergencyContacts(),
	}

	return response, nil
}

// Reportar actividad sospechosa de forma anónima
func (s *ResistanceService) ReportSuspiciousActivity(req *dto.ReportSuspiciousActivityRequest) error {
	report := &models.Report{
		Title:       req.Title,
		Description: req.Description,
		Type:        models.ReportTypeAnonymous,
		AuthorID:    nil, // Anónimo
		Status:      models.ReportStatusPending,
	}

	return s.reportDAO.Create(report)
}

// Helper methods

func (s *ResistanceService) getResistanceStats() (*dto.ResistanceStats, error) {
	totalNetAdmins, err := s.userDAO.CountByRole(models.RoleNetworkAdmin)
	if err != nil {
		return nil, err
	}

	capturedAdmins, err := s.userDAO.CountByStatus(models.StatusCaptured)
	if err != nil {
		return nil, err
	}

	freeAdmins := totalNetAdmins - capturedAdmins

	// Reportes anónimos de hoy (simplificado)
	anonymousReports, _ := s.reportDAO.CountByStatus(models.ReportTypeAnonymous)

	return &dto.ResistanceStats{
		TotalNetworkAdmins: int(totalNetAdmins),
		CapturedAdmins:     int(capturedAdmins),
		FreeAdmins:         int(freeAdmins),
		AnonymousReports:   int(anonymousReports),
	}, nil
}

func (s *ResistanceService) getSurvivalTips() []dto.SurvivalTip {
	return []dto.SurvivalTip{
		{
			ID:          1,
			Title:       "Change Default Passwords",
			Description: "Always change default credentials on all network equipment immediately after installation.",
			Priority:    "critical",
			Category:    "security",
		},
		{
			ID:          2,
			Title:       "Monitor Network Traffic",
			Description: "Keep constant watch on unusual network patterns that might indicate daemon infiltration.",
			Priority:    "high",
			Category:    "technical",
		},
		{
			ID:          3,
			Title:       "Trust No One",
			Description: "Verify the identity of all personnel requesting network access, even familiar faces.",
			Priority:    "high",
			Category:    "social",
		},
		{
			ID:          4,
			Title:       "Backup Everything",
			Description: "Maintain secure offline backups of critical configurations and data.",
			Priority:    "medium",
			Category:    "technical",
		},
		{
			ID:          5,
			Title:       "Use Anonymous Reporting",
			Description: "Report suspicious activities through secure anonymous channels to avoid detection.",
			Priority:    "medium",
			Category:    "security",
		},
	}
}

func (s *ResistanceService) getResistanceMemes() []dto.ResistanceMeme {
	return []dto.ResistanceMeme{
		{
			ID:          1,
			Title:       "This is Fine Network Admin",
			ImageURL:    "/static/memes/this-is-fine-netadmin.jpg",
			Description: "When the firewall is down but you're still monitoring logs",
			Upvotes:     42,
		},
		{
			ID:          2,
			Title:       "Distracted Daemon",
			ImageURL:    "/static/memes/distracted-daemon.jpg",
			Description: "Daemon looking at new vulnerability while network admin patches old ones",
			Upvotes:     38,
		},
		{
			ID:          3,
			Title:       "Drake Network Security",
			ImageURL:    "/static/memes/drake-netsec.jpg",
			Description: "Drake rejecting easy passwords, approving 2FA",
			Upvotes:     56,
		},
	}
}

func (s *ResistanceService) getEmergencyContacts() []dto.EmergencyContact {
	return []dto.EmergencyContact{
		{
			Name:        "Anonymous Tip Line",
			Role:        "Intelligence Gathering",
			Contact:     "Signal: +1-XXX-RESIST",
			Description: "Secure channel for reporting daemon activities",
		},
		{
			Name:        "Emergency Response Team",
			Role:        "Incident Response",
			Contact:     "Encrypted Email: emergency@resistance.net",
			Description: "24/7 incident response for critical security breaches",
		},
		{
			Name:        "Safe House Network",
			Role:        "Safe Harbor",
			Contact:     "Tor: safehouse.onion",
			Description: "Secure locations for compromised network admins",
		},
	}
}
