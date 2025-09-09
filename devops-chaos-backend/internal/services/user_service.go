package services

import (
	"devops-chaos-backend/internal/dao"
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/models"
	"errors"
	"math"
)

type UserService struct {
	userDAO      *dao.UserDAO
	statisticDAO *dao.StatisticDAO
}

func NewUserService() *UserService {
	return &UserService{
		userDAO:      dao.NewUserDAO(),
		statisticDAO: dao.NewStatisticDAO(),
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

	// Actualizar status del network admin
	target.Status = models.StatusCaptured
	err = s.userDAO.Update(target)
	if err != nil {
		return err
	}

	// Incrementar capturas del daemon
	err = s.statisticDAO.IncrementCaptures(daemonID)
	if err != nil {
		// Log error pero no fallar la operación
	}

	// Dar puntos al daemon
	err = s.statisticDAO.AddPoints(daemonID, 10) // 10 puntos por captura
	if err != nil {
		// Log error pero no fallar la operación
	}

	return nil
}
