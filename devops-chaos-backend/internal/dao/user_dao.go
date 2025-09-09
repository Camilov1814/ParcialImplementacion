package dao

import (
	"devops-chaos-backend/internal/config"
	"devops-chaos-backend/internal/models"
)

type UserDAO struct{}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

// Crear usuario
func (dao *UserDAO) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

// Buscar usuario por ID
func (dao *UserDAO) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Buscar usuario por username
func (dao *UserDAO) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Buscar usuario por email
func (dao *UserDAO) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Listar todos los usuarios con paginación
func (dao *UserDAO) FindAll(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Contar total
	config.DB.Model(&models.User{}).Count(&total)

	// Obtener usuarios con paginación
	offset := (page - 1) * limit
	err := config.DB.Offset(offset).Limit(limit).Find(&users).Error

	return users, total, err
}

// Listar usuarios por rol
func (dao *UserDAO) FindByRole(role string, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := config.DB.Where("role = ?", role)

	// Contar total
	query.Model(&models.User{}).Count(&total)

	// Obtener usuarios con paginación
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&users).Error

	return users, total, err
}

// Actualizar usuario
func (dao *UserDAO) Update(user *models.User) error {
	return config.DB.Save(user).Error
}

// Eliminar usuario (soft delete)
func (dao *UserDAO) Delete(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}

// Contar usuarios por status
func (dao *UserDAO) CountByStatus(status string) (int64, error) {
	var count int64
	err := config.DB.Model(&models.User{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// Contar usuarios por rol
func (dao *UserDAO) CountByRole(role string) (int64, error) {
	var count int64
	err := config.DB.Model(&models.User{}).Where("role = ?", role).Count(&count).Error
	return count, err
}

// Verificar si username ya existe
func (dao *UserDAO) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := config.DB.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// Verificar si email ya existe
func (dao *UserDAO) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := config.DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
