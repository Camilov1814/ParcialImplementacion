package services

import (
	"devops-chaos-backend/internal/dao"
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/models"
	"errors"
	"math"
	"time"
)

type PunishmentService struct {
	punishmentDAO *dao.PunishmentDAO
	userDAO       *dao.UserDAO
	statisticDAO  *dao.StatisticDAO
}

func NewPunishmentService() *PunishmentService {
	return &PunishmentService{
		punishmentDAO: dao.NewPunishmentDAO(),
		userDAO:       dao.NewUserDAO(),
		statisticDAO:  dao.NewStatisticDAO(),
	}
}

// Crear castigo/recompensa (solo Andrei)
func (s *PunishmentService) CreatePunishment(req *dto.CreatePunishmentRequest, assignerID uint, assignerRole string) (*dto.PunishmentResponse, error) {
	if assignerRole != models.RoleAndrei {
		return nil, errors.New("unauthorized: only Andrei can assign punishments")
	}

	// Verificar que el objetivo existe
	target, err := s.userDAO.FindByID(req.TargetID)
	if err != nil {
		return nil, errors.New("target user not found")
	}

	// Solo se puede castigar/recompensar a Daemons
	if target.Role != models.RoleDaemon {
		return nil, errors.New("can only punish/reward daemons")
	}

	// Validar tipo de castigo
	if !s.isValidPunishmentType(req.Type) {
		return nil, errors.New("invalid punishment type")
	}

	// Parsear fecha de expiración si se proporciona
	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			return nil, errors.New("invalid expires_at format, use RFC3339 (2024-01-15T10:00:00Z)")
		}
		expiresAt = &parsedTime
	}

	// Crear castigo
	punishment := &models.Punishment{
		TargetID:    req.TargetID,
		AssignedBy:  assignerID,
		Type:        req.Type,
		Description: req.Description,
		Status:      models.PunishmentStatusActive,
		ExpiresAt:   expiresAt,
	}

	err = s.punishmentDAO.Create(punishment)
	if err != nil {
		return nil, err
	}

	// Actualizar status del usuario si es castigo
	if s.isPunishmentType(req.Type) {
		target.Status = models.StatusPunished
		s.userDAO.Update(target)
	}

	// Si es recompensa, dar puntos adicionales
	if req.Type == models.PunishmentTypeReward {
		s.statisticDAO.AddPoints(req.TargetID, 20) // 20 puntos por recompensa
	}

	// Cargar el castigo completo
	createdPunishment, err := s.punishmentDAO.FindByID(punishment.ID)
	if err != nil {
		return nil, err
	}

	return s.convertToPunishmentResponse(createdPunishment), nil
}

// Obtener castigo por ID
func (s *PunishmentService) GetPunishmentByID(id uint, userID uint, userRole string) (*dto.PunishmentResponse, error) {
	punishment, err := s.punishmentDAO.FindByID(id)
	if err != nil {
		return nil, errors.New("punishment not found")
	}

	// Control de acceso
	if !s.canAccessPunishment(punishment, userID, userRole) {
		return nil, errors.New("unauthorized to access this punishment")
	}

	return s.convertToPunishmentResponse(punishment), nil
}

// Listar castigos con paginación
func (s *PunishmentService) GetPunishments(page, limit int, status string, userID uint, userRole string) (*dto.PaginatedResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var punishments []models.Punishment
	var total int64
	var err error

	switch userRole {
	case models.RoleAndrei:
		// Andrei puede ver todos los castigos
		if status != "" {
			punishments, total, err = s.punishmentDAO.FindByStatus(status, page, limit)
		} else {
			punishments, total, err = s.punishmentDAO.FindAll(page, limit)
		}

	case models.RoleDaemon:
		// Daemon solo ve sus propios castigos
		punishments, total, err = s.punishmentDAO.FindByTarget(userID, page, limit)

	default:
		return nil, errors.New("unauthorized")
	}

	if err != nil {
		return nil, err
	}

	// Convertir a DTO
	punishmentItems := make([]dto.PunishmentListItem, len(punishments))
	for i, punishment := range punishments {
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

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.PaginatedResponse{
		Success:    true,
		Message:    "Punishments retrieved successfully",
		Data:       punishmentItems,
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

// Actualizar castigo (solo Andrei)
func (s *PunishmentService) UpdatePunishment(punishmentID uint, req *dto.UpdatePunishmentRequest, userRole string) error {
	if userRole != models.RoleAndrei {
		return errors.New("unauthorized: only Andrei can update punishments")
	}

	punishment, err := s.punishmentDAO.FindByID(punishmentID)
	if err != nil {
		return errors.New("punishment not found")
	}

	// Actualizar campos si se proporcionan
	if req.Status != "" {
		if !s.isValidPunishmentStatus(req.Status) {
			return errors.New("invalid status")
		}

		oldStatus := punishment.Status
		punishment.Status = req.Status

		// Si se completa o cancela un castigo, actualizar status del usuario
		if oldStatus == models.PunishmentStatusActive &&
			(req.Status == models.PunishmentStatusCompleted || req.Status == models.PunishmentStatusCancelled) {
			target, err := s.userDAO.FindByID(punishment.TargetID)
			if err == nil && target.Status == models.StatusPunished {
				target.Status = models.StatusActive
				s.userDAO.Update(target)
			}
		}
	}

	if req.Description != "" {
		punishment.Description = req.Description
	}

	if req.ExpiresAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			return errors.New("invalid expires_at format")
		}
		punishment.ExpiresAt = &parsedTime
	}

	return s.punishmentDAO.Update(punishment)
}

// Eliminar castigo (solo Andrei)
func (s *PunishmentService) DeletePunishment(punishmentID uint, userRole string) error {
	if userRole != models.RoleAndrei {
		return errors.New("unauthorized: only Andrei can delete punishments")
	}

	punishment, err := s.punishmentDAO.FindByID(punishmentID)
	if err != nil {
		return errors.New("punishment not found")
	}

	// Si el castigo está activo, actualizar status del usuario
	if punishment.Status == models.PunishmentStatusActive {
		target, err := s.userDAO.FindByID(punishment.TargetID)
		if err == nil && target.Status == models.StatusPunished {
			target.Status = models.StatusActive
			s.userDAO.Update(target)
		}
	}

	return s.punishmentDAO.Delete(punishmentID)
}

// Obtener castigos activos de un daemon
func (s *PunishmentService) GetActivePunishments(targetID uint, userID uint, userRole string) ([]dto.PunishmentListItem, error) {
	// Solo el daemon afectado o Andrei pueden ver castigos activos
	if userRole != models.RoleAndrei && userID != targetID {
		return nil, errors.New("unauthorized")
	}

	punishments, err := s.punishmentDAO.FindActiveByTarget(targetID)
	if err != nil {
		return nil, err
	}

	items := make([]dto.PunishmentListItem, len(punishments))
	for i, punishment := range punishments {
		items[i] = dto.PunishmentListItem{
			ID:           punishment.ID,
			TargetName:   punishment.Target.Username,
			AssignerName: punishment.Assigner.Username,
			Type:         punishment.Type,
			Status:       punishment.Status,
			CreatedAt:    punishment.CreatedAt,
			ExpiresAt:    punishment.ExpiresAt,
		}
	}

	return items, nil
}

// Helper methods

func (s *PunishmentService) canAccessPunishment(punishment *models.Punishment, userID uint, userRole string) bool {
	switch userRole {
	case models.RoleAndrei:
		return true // Andrei puede ver todo
	case models.RoleDaemon:
		return punishment.TargetID == userID // Solo sus propios castigos
	default:
		return false
	}
}

func (s *PunishmentService) isValidPunishmentType(punishmentType string) bool {
	validTypes := []string{
		models.PunishmentTypeTimeout,
		models.PunishmentTypeDemotion,
		models.PunishmentTypeExtraTasks,
		models.PunishmentTypeReward,
	}
	for _, validType := range validTypes {
		if punishmentType == validType {
			return true
		}
	}
	return false
}

func (s *PunishmentService) isValidPunishmentStatus(status string) bool {
	validStatuses := []string{
		models.PunishmentStatusActive,
		models.PunishmentStatusCompleted,
		models.PunishmentStatusCancelled,
	}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func (s *PunishmentService) isPunishmentType(punishmentType string) bool {
	punishmentTypes := []string{
		models.PunishmentTypeTimeout,
		models.PunishmentTypeDemotion,
		models.PunishmentTypeExtraTasks,
	}
	for _, pType := range punishmentTypes {
		if punishmentType == pType {
			return true
		}
	}
	return false
}

func (s *PunishmentService) convertToPunishmentResponse(punishment *models.Punishment) *dto.PunishmentResponse {
	return &dto.PunishmentResponse{
		ID: punishment.ID,
		Target: dto.UserInfo{
			ID:       punishment.Target.ID,
			Username: punishment.Target.Username,
			Email:    punishment.Target.Email,
			Role:     punishment.Target.Role,
			Status:   punishment.Target.Status,
		},
		Assigner: dto.UserInfo{
			ID:       punishment.Assigner.ID,
			Username: punishment.Assigner.Username,
			Email:    punishment.Assigner.Email,
			Role:     punishment.Assigner.Role,
			Status:   punishment.Assigner.Status,
		},
		Type:        punishment.Type,
		Description: punishment.Description,
		Status:      punishment.Status,
		CreatedAt:   punishment.CreatedAt,
		UpdatedAt:   punishment.UpdatedAt,
		ExpiresAt:   punishment.ExpiresAt,
	}
}
