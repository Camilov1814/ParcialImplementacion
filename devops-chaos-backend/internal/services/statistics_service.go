package services

import (
	"devops-chaos-backend/internal/dao"
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/models"
	"errors"
)

type StatisticsService struct {
	statisticDAO *dao.StatisticDAO
	userDAO      *dao.UserDAO
}

func NewStatisticsService() *StatisticsService {
	return &StatisticsService{
		statisticDAO: dao.NewStatisticDAO(),
		userDAO:      dao.NewUserDAO(),
	}
}

// Obtener estadísticas de un usuario
func (s *StatisticsService) GetUserStatistics(targetUserID, currentUserID uint, currentUserRole string) (*dto.StatisticResponse, error) {
	// Solo Andrei o el mismo usuario pueden ver las estadísticas
	if currentUserRole != models.RoleAndrei && currentUserID != targetUserID {
		return nil, errors.New("unauthorized to view these statistics")
	}

	// Verificar que el usuario objetivo sea daemon
	targetUser, err := s.userDAO.FindByID(targetUserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if targetUser.Role != models.RoleDaemon {
		return nil, errors.New("statistics only available for daemons")
	}

	// Obtener estadísticas
	stats, err := s.statisticDAO.FindByUserIDWithUser(targetUserID)
	if err != nil {
		return nil, errors.New("statistics not found")
	}

	response := &dto.StatisticResponse{
		ID:            stats.ID,
		UserID:        stats.UserID,
		Username:      stats.User.Username,
		CapturesCount: stats.CapturesCount,
		ReportsCount:  stats.ReportsCount,
		Ranking:       stats.Ranking,
		Points:        stats.Points,
		UpdatedAt:     stats.UpdatedAt,
	}

	return response, nil
}

// Obtener leaderboard
func (s *StatisticsService) GetLeaderboard(limit int, userRole string) ([]dto.RankingItem, error) {
	if limit <= 0 {
		limit = 10
	}

	// Solo Andrei y Daemons pueden ver el leaderboard
	if userRole != models.RoleAndrei && userRole != models.RoleDaemon {
		return nil, errors.New("unauthorized to view leaderboard")
	}

	topDaemons, err := s.statisticDAO.FindTopDaemons(limit)
	if err != nil {
		return nil, err
	}

	leaderboard := make([]dto.RankingItem, len(topDaemons))
	for i, stat := range topDaemons {
		leaderboard[i] = dto.RankingItem{
			Position:      i + 1, // Posición basada en orden
			Username:      stat.User.Username,
			CapturesCount: stat.CapturesCount,
			ReportsCount:  stat.ReportsCount,
			Points:        stat.Points,
			Status:        stat.User.Status,
		}
	}

	return leaderboard, nil
}

// Actualizar estadísticas manualmente (solo Andrei)
func (s *StatisticsService) UpdateStatistics(targetUserID uint, req *dto.UpdateStatisticRequest, currentUserRole string) error {
	if currentUserRole != models.RoleAndrei {
		return errors.New("unauthorized: only Andrei can update statistics")
	}

	// Verificar que el usuario objetivo sea daemon
	targetUser, err := s.userDAO.FindByID(targetUserID)
	if err != nil {
		return errors.New("user not found")
	}

	if targetUser.Role != models.RoleDaemon {
		return errors.New("can only update daemon statistics")
	}

	// Obtener estadísticas actuales
	stats, err := s.statisticDAO.FindByUserID(targetUserID)
	if err != nil {
		// Si no existen, crear nuevas
		stats, err = s.statisticDAO.CreateOrUpdate(targetUserID)
		if err != nil {
			return err
		}
	}

	// Actualizar campos proporcionados
	if req.CapturesCount != nil {
		stats.CapturesCount = *req.CapturesCount
	}

	if req.ReportsCount != nil {
		stats.ReportsCount = *req.ReportsCount
	}

	if req.Points != nil {
		stats.Points = *req.Points
	}

	err = s.statisticDAO.Update(stats)
	if err != nil {
		return err
	}

	// Recalcular rankings después de la actualización
	return s.statisticDAO.RecalculateRankings()
}

// Recalcular rankings (solo Andrei)
func (s *StatisticsService) RecalculateRankings(currentUserRole string) error {
	if currentUserRole != models.RoleAndrei {
		return errors.New("unauthorized: only Andrei can recalculate rankings")
	}

	return s.statisticDAO.RecalculateRankings()
}
