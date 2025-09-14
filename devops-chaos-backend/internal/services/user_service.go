package services

import (
	"devops-chaos-backend/internal/dao"
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/models"
	"errors"
	"math"
	"time"
)

type UserService struct {
	userDAO      *dao.UserDAO
	statisticDAO *dao.StatisticDAO
	captureDAO   *dao.CaptureDAO
}

func NewUserService() *UserService {
	return &UserService{
		userDAO:      dao.NewUserDAO(),
		statisticDAO: dao.NewStatisticDAO(),
		captureDAO:   dao.NewCaptureDAO(),
	}
}

// Obtener usuario por ID
func (s *UserService) GetByID(id uint, currentUserRole string) (*dto.UserInfo, error) {
	user, err := s.userDAO.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Solo Andrei puede ver información completa de cualquier usuario
	if currentUserRole != models.RoleAndrei && user.ID != id {
		return nil, errors.New("unauthorized")
	}

	userInfo := &dto.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Status:   user.Status,
	}

	return userInfo, nil
}

// Listar usuarios con paginación (solo Andrei)
func (s *UserService) GetUsers(page, limit int, role, currentUserRole string) (*dto.PaginatedResponse, error) {
	if currentUserRole != models.RoleAndrei {
		return nil, errors.New("unauthorized: only Andrei can list users")
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var users []models.User
	var total int64
	var err error

	if role != "" {
		users, total, err = s.userDAO.FindByRole(role, page, limit)
	} else {
		users, total, err = s.userDAO.FindAll(page, limit)
	}

	if err != nil {
		return nil, err
	}

	// Convertir a DTO
	userItems := make([]dto.UserListItem, len(users))
	for i, user := range users {
		userItems[i] = dto.UserListItem{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.PaginatedResponse{
		Success:    true,
		Message:    "Users retrieved successfully",
		Data:       userItems,
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

// Actualizar usuario (solo Andrei)
func (s *UserService) UpdateUser(userID uint, req *dto.UpdateUserRequest, currentUserRole string) error {
	if currentUserRole != models.RoleAndrei {
		return errors.New("unauthorized: only Andrei can update users")
	}

	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Actualizar campos si se proporcionan
	if req.Username != "" {
		// Verificar que no exista otro usuario con ese username
		existing, _ := s.userDAO.FindByUsername(req.Username)
		if existing != nil && existing.ID != userID {
			return errors.New("username already exists")
		}
		user.Username = req.Username
	}

	if req.Email != "" {
		// Verificar que no exista otro usuario con ese email
		existing, _ := s.userDAO.FindByEmail(req.Email)
		if existing != nil && existing.ID != userID {
			return errors.New("email already exists")
		}
		user.Email = req.Email
	}

	if req.Role != "" {
		user.Role = req.Role
	}

	if req.Status != "" {
		user.Status = req.Status
	}

	return s.userDAO.Update(user)
}

// Eliminar usuario (solo Andrei)
func (s *UserService) DeleteUser(userID uint, currentUserRole string) error {
	if currentUserRole != models.RoleAndrei {
		return errors.New("unauthorized: only Andrei can delete users")
	}

	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// No permitir eliminar a Andrei
	if user.Role == models.RoleAndrei {
		return errors.New("cannot delete Andrei")
	}

	return s.userDAO.Delete(userID)
}

// Obtener estadísticas generales de usuarios (solo Andrei)
func (s *UserService) GetUserStats(currentUserRole string) (map[string]interface{}, error) {
	if currentUserRole != models.RoleAndrei {
		return nil, errors.New("unauthorized: only Andrei can view user stats")
	}

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

	stats := map[string]interface{}{
		"total_users":      totalUsers,
		"total_daemons":    totalDaemons,
		"total_net_admins": totalNetAdmins,
		"captured_admins":  capturedAdmins,
		"punished_daemons": punishedDaemons,
	}

	return stats, nil
}

// Capturar Network Admin (para Daemons)
func (s *UserService) CaptureNetworkAdmin(targetID uint, daemonID uint, daemonRole string) error {
	if daemonRole != models.RoleDaemon {
		return errors.New("unauthorized: only daemons can capture network admins")
	}

	// Verificar que el objetivo sea un network admin
	target, err := s.userDAO.FindByID(targetID)
	if err != nil {
		return errors.New("target user not found")
	}

	if target.Role != models.RoleNetworkAdmin {
		return errors.New("can only capture network administrators")
	}

	if target.Status == models.StatusCaptured {
		return errors.New("network admin is already captured")
	}

	// Calcular dificultad y puntos basado en el ID del target
	difficulty := s.calculateCaptureDifficulty(targetID)
	points := s.calculateCapturePoints(difficulty)

	// Actualizar status del network admin
	target.Status = models.StatusCaptured
	err = s.userDAO.Update(target)
	if err != nil {
		return err
	}

	// Crear registro de captura
	capture := &models.Capture{
		DaemonID:    daemonID,
		TargetID:    targetID,
		CaptureDate: time.Now(),
		Status:      models.CaptureStatusCaptured,
		Method:      "network_infiltration",
		Points:      points,
		Difficulty:  difficulty,
	}

	err = s.captureDAO.Create(capture)
	if err != nil {
		return err
	}

	// Incrementar capturas del daemon
	err = s.statisticDAO.IncrementCaptures(daemonID)
	if err != nil {
		// Log error pero no fallar la operación
	}

	// Dar puntos al daemon
	err = s.statisticDAO.AddPoints(daemonID, points)
	if err != nil {
		// Log error pero no fallar la operación
	}

	return nil
}

// Función auxiliar para calcular la dificultad de captura
func (s *UserService) calculateCaptureDifficulty(targetID uint) string {
	// Simulamos dificultad basada en el ID
	difficulty := (targetID % 3) + 1
	switch difficulty {
	case 1:
		return models.CaptureDifficultyEasy
	case 2:
		return models.CaptureDifficultyMedium
	case 3:
		return models.CaptureDifficultyHard
	default:
		return models.CaptureDifficultyMedium
	}
}

// Función auxiliar para calcular los puntos de captura
func (s *UserService) calculateCapturePoints(difficulty string) int {
	switch difficulty {
	case models.CaptureDifficultyEasy:
		return 100
	case models.CaptureDifficultyMedium:
		return 250
	case models.CaptureDifficultyHard:
		return 500
	default:
		return 100
	}
}

// Obtener capturas de un daemon
func (s *UserService) GetDaemonCaptures(daemonID uint, currentUserRole string, currentUserID uint) ([]models.Capture, error) {
	// Los daemons solo pueden ver sus propias capturas, Andrei puede ver todas
	if currentUserRole == models.RoleDaemon && currentUserID != daemonID {
		return nil, errors.New("unauthorized: can only view your own captures")
	}

	if currentUserRole != models.RoleDaemon && currentUserRole != models.RoleAndrei {
		return nil, errors.New("unauthorized: only daemons and Andrei can view captures")
	}

	return s.captureDAO.FindByDaemonID(daemonID)
}

// Obtener todas las capturas (solo Andrei)
func (s *UserService) GetAllCaptures(currentUserRole string) ([]models.Capture, error) {
	if currentUserRole != models.RoleAndrei {
		return nil, errors.New("unauthorized: only Andrei can view all captures")
	}

	return s.captureDAO.FindAll()
}

// Obtener network admins disponibles para captura (para Daemons)
func (s *UserService) GetNetworkAdminsForCapture(currentUserRole string) ([]models.User, error) {
	if currentUserRole != models.RoleDaemon {
		return nil, errors.New("unauthorized: only daemons can view capture targets")
	}

	// Obtener todos los network admins activos (no capturados)
	users, _, err := s.userDAO.FindByRole(models.RoleNetworkAdmin, 1, 100)
	if err != nil {
		return nil, err
	}

	// Filtrar solo los que están activos (no capturados)
	var availableAdmins []models.User
	for _, user := range users {
		if user.Status == models.StatusActive {
			availableAdmins = append(availableAdmins, user)
		}
	}

	return availableAdmins, nil
}
