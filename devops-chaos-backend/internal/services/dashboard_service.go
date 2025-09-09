package services

import (
	"devops-chaos-backend/internal/dao"
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/models"
	"fmt"
	"time"
)

type DashboardService struct {
	userDAO       *dao.UserDAO
	reportDAO     *dao.ReportDAO
	punishmentDAO *dao.PunishmentDAO
	statisticDAO  *dao.StatisticDAO
}

func NewDashboardService() *DashboardService {
	return &DashboardService{
		userDAO:       dao.NewUserDAO(),
		reportDAO:     dao.NewReportDAO(),
		punishmentDAO: dao.NewPunishmentDAO(),
		statisticDAO:  dao.NewStatisticDAO(),
	}
}

// Dashboard de Andrei
func (s *DashboardService) GetAndreiDashboard() (*dto.AndreiDashboardResponse, error) {
	// Estadísticas del sistema
	systemStats, err := s.getSystemStats()
	if err != nil {
		return nil, err
	}

	// Reportes recientes
	recentReports, err := s.reportDAO.FindRecent(5)
	if err != nil {
		return nil, err
	}

	reportItems := make([]dto.ReportListItem, len(recentReports))
	for i, report := range recentReports {
		authorName := ""
		if report.Author != nil {
			authorName = report.Author.Username
		}

		reportItems[i] = dto.ReportListItem{
			ID:        report.ID,
			Title:     report.Title,
			Type:      report.Type,
			Status:    report.Status,
			AuthorID:  report.AuthorID,
			Author:    authorName,
			CreatedAt: report.CreatedAt,
		}
	}

	// Top daemons
	topDaemons, err := s.statisticDAO.FindTopDaemons(5)
	if err != nil {
		return nil, err
	}

	topPerformers := make([]dto.TopPerformerResponse, len(topDaemons))
	for i, stat := range topDaemons {
		topPerformers[i] = dto.TopPerformerResponse{
			Username:      stat.User.Username,
			Points:        stat.Points,
			CapturesCount: stat.CapturesCount,
			ReportsCount:  stat.ReportsCount,
		}
	}

	// Actividad reciente (simplificada)
	recentActivity := s.generateRecentActivity()

	response := &dto.AndreiDashboardResponse{
		WelcomeMessage: "Welcome back, Andrei. The chaos continues...",
		SystemStats:    *systemStats,
		RecentReports:  reportItems,
		TopDaemons:     topPerformers,
		RecentActivity: recentActivity,
	}

	return response, nil
}

// Dashboard de Daemon
func (s *DashboardService) GetDaemonDashboard(userID uint) (*dto.DaemonDashboardResponse, error) {
	// Estadísticas del usuario
	userStats, err := s.statisticDAO.FindByUserIDWithUser(userID)
	if err != nil {
		return nil, err
	}

	daemonStats := dto.DaemonStats{
		CapturesCount: userStats.CapturesCount,
		ReportsCount:  userStats.ReportsCount,
		Points:        userStats.Points,
		Ranking:       userStats.Ranking,
		Status:        userStats.User.Status,
	}

	// Leaderboard
	topDaemons, err := s.statisticDAO.FindTopDaemons(10)
	if err != nil {
		return nil, err
	}

	leaderboard := make([]dto.RankingItem, len(topDaemons))
	for i, stat := range topDaemons {
		leaderboard[i] = dto.RankingItem{
			Position:      stat.Ranking,
			Username:      stat.User.Username,
			CapturesCount: stat.CapturesCount,
			ReportsCount:  stat.ReportsCount,
			Points:        stat.Points,
			Status:        stat.User.Status,
		}
	}

	// Castigos activos
	activePunishments, err := s.punishmentDAO.FindActiveByTarget(userID)
	if err != nil {
		return nil, err
	}

	punishmentItems := make([]dto.PunishmentListItem, len(activePunishments))
	for i, punishment := range activePunishments {
		punishmentItems[i] = dto.PunishmentListItem{
			ID:           punishment.ID,
			TargetName:   punishment.Target.Username,
			AssignerName: punishment.Assigner.Username,
			Type:         punishment.Type,
			Status:       punishment.Status,
			CreatedAt:    punishment.CreatedAt,
			ExpiresAt:    punishment.ExpiresAt,
		}
	}

	// Misiones activas (hardcoded para el ejemplo)
	activeMissions := s.generateActiveMissions()

	// Actividad reciente de caos
	recentChaos := s.generateRecentChaos()

	response := &dto.DaemonDashboardResponse{
		WelcomeMessage:    fmt.Sprintf("Welcome back, %s. Ready to spread more chaos?", userStats.User.Username),
		UserStats:         daemonStats,
		ActiveMissions:    activeMissions,
		Leaderboard:       leaderboard,
		RecentCaos:        recentChaos,
		ActivePunishments: punishmentItems,
	}

	return response, nil
}

// Helper methods

func (s *DashboardService) getSystemStats() (*dto.SystemStats, error) {
	totalUsers, err := s.userDAO.CountByRole("")
	if err != nil {
		return nil, err
	}

	totalDaemons, err := s.userDAO.CountByRole(models.RoleDaemon)
	if err != nil {
		return nil, err
	}

	totalNetAdmins, err := s.userDAO.CountByRole(models.RoleNetworkAdmin)
	if err != nil {
		return nil, err
	}

	capturedAdmins, err := s.userDAO.CountByStatus(models.StatusCaptured)
	if err != nil {
		return nil, err
	}

	punishedDaemons, err := s.userDAO.CountByStatus(models.StatusPunished)
	if err != nil {
		return nil, err
	}

	pendingReports, err := s.reportDAO.CountByStatus(models.ReportStatusPending)
	if err != nil {
		return nil, err
	}

	totalReports, err := s.reportDAO.CountByStatus("")
	if err != nil {
		return nil, err
	}

	return &dto.SystemStats{
		TotalUsers:      int(totalUsers),
		TotalDaemons:    int(totalDaemons),
		TotalNetAdmins:  int(totalNetAdmins),
		CapturedAdmins:  int(capturedAdmins),
		PunishedDaemons: int(punishedDaemons),
		PendingReports:  int(pendingReports),
		TotalReports:    int(totalReports),
	}, nil
}

func (s *DashboardService) generateRecentActivity() []dto.ActivityItem {
	// En una implementación real, esto vendría de la base de datos
	return []dto.ActivityItem{
		{
			ID:        1,
			Type:      "capture",
			Message:   "Network Admin 'alice' was captured by daemon 'bob'",
			Timestamp: time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			UserID:    2,
			Username:  "bob",
		},
		{
			ID:        2,
			Type:      "report",
			Message:   "New resistance report submitted anonymously",
			Timestamp: time.Now().Add(-4 * time.Hour).Format(time.RFC3339),
			UserID:    0,
			Username:  "anonymous",
		},
		{
			ID:        3,
			Type:      "punishment",
			Message:   "Daemon 'charlie' received timeout punishment",
			Timestamp: time.Now().Add(-6 * time.Hour).Format(time.RFC3339),
			UserID:    3,
			Username:  "charlie",
		},
	}
}

func (s *DashboardService) generateActiveMissions() []dto.MissionItem {
	return []dto.MissionItem{
		{
			ID:          1,
			Title:       "Infiltrate Server Room",
			Description: "Gain physical access to the main server room",
			Difficulty:  "high",
			Points:      25,
			Status:      "active",
		},
		{
			ID:          2,
			Title:       "Social Engineering",
			Description: "Extract credentials through social manipulation",
			Difficulty:  "medium",
			Points:      15,
			Status:      "active",
		},
		{
			ID:          3,
			Title:       "Network Reconnaissance",
			Description: "Map the internal network topology",
			Difficulty:  "low",
			Points:      10,
			Status:      "available",
		},
	}
}

func (s *DashboardService) generateRecentChaos() []dto.ChaosActivityItem {
	return []dto.ChaosActivityItem{
		{
			ID:          1,
			Description: "Successfully compromised DNS server",
			Type:        "technical",
			Timestamp:   time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			Points:      20,
		},
		{
			ID:          2,
			Description: "Deployed network monitoring backdoor",
			Type:        "surveillance",
			Timestamp:   time.Now().Add(-3 * time.Hour).Format(time.RFC3339),
			Points:      15,
		},
		{
			ID:          3,
			Description: "Intercepted admin communications",
			Type:        "intelligence",
			Timestamp:   time.Now().Add(-5 * time.Hour).Format(time.RFC3339),
			Points:      10,
		},
	}
}
