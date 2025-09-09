package services

import (
	"devops-chaos-backend/internal/dao"
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/models"
	"devops-chaos-backend/internal/utils"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type AuthService struct {
	userDAO      *dao.UserDAO
	statisticDAO *dao.StatisticDAO
}

func NewAuthService() *AuthService {
	return &AuthService{
		userDAO:      dao.NewUserDAO(),
		statisticDAO: dao.NewStatisticDAO(),
	}
}

// Login del usuario
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Buscar usuario por username
	user, err := s.userDAO.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Verificar contraseña
	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generar token JWT
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("could not generate token")
	}

	// Preparar respuesta
	response := &dto.LoginResponse{
		Token: token,
		User: dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Status:   user.Status,
		},
	}

	return response, nil
}

// Registrar nuevo usuario (solo Andrei puede hacerlo)
func (s *AuthService) Register(req *dto.RegisterRequest, currentUserRole string) (*dto.UserInfo, error) {
	// Solo Andrei puede registrar usuarios
	if currentUserRole != models.RoleAndrei {
		return nil, errors.New("unauthorized: only Andrei can register users")
	}

	// Validar rol
	if !s.isValidRole(req.Role) {
		return nil, errors.New("invalid role")
	}

	// Verificar que username no exista
	exists, err := s.userDAO.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Verificar que email no exista
	exists, err = s.userDAO.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Hash de la contraseña
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("could not hash password")
	}

	// Crear usuario
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
		Status:   models.StatusActive,
	}

	err = s.userDAO.Create(user)
	if err != nil {
		return nil, err
	}

	// Si es daemon, crear estadísticas
	if user.Role == models.RoleDaemon {
		_, err = s.statisticDAO.CreateOrUpdate(user.ID)
		if err != nil {
			// Log el error pero no falla el registro
			// En producción, usar un logger apropiado
		}
	}

	// Preparar respuesta
	userInfo := &dto.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Status:   user.Status,
	}

	return userInfo, nil
}

// Cambiar contraseña
func (s *AuthService) ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	// Buscar usuario
	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verificar contraseña actual
	if !utils.CheckPassword(user.Password, req.OldPassword) {
		return errors.New("current password is incorrect")
	}

	// Hash nueva contraseña
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("could not hash password")
	}

	// Actualizar contraseña
	user.Password = hashedPassword
	return s.userDAO.Update(user)
}

// Validar token y obtener usuario
func (s *AuthService) ValidateTokenAndGetUser(tokenString string) (*models.User, error) {
	// Limpiar token
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	// Validar token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Buscar usuario
	user, err := s.userDAO.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// Validar roles
func (s *AuthService) isValidRole(role string) bool {
	validRoles := []string{models.RoleAndrei, models.RoleDaemon, models.RoleNetworkAdmin}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}
