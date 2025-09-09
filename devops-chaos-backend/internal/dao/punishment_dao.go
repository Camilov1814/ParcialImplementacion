package dao

import (
	"devops-chaos-backend/internal/config"
	"devops-chaos-backend/internal/models"
)

type PunishmentDAO struct{}

func NewPunishmentDAO() *PunishmentDAO {
	return &PunishmentDAO{}
}

// Crear castigo
func (dao *PunishmentDAO) Create(punishment *models.Punishment) error {
	return config.DB.Create(punishment).Error
}

// Buscar castigo por ID
func (dao *PunishmentDAO) FindByID(id uint) (*models.Punishment, error) {
	var punishment models.Punishment
	err := config.DB.Preload("Target").Preload("Assigner").First(&punishment, id).Error
	if err != nil {
		return nil, err
	}
	return &punishment, nil
}

// Listar todos los castigos con paginación
func (dao *PunishmentDAO) FindAll(page, limit int) ([]models.Punishment, int64, error) {
	var punishments []models.Punishment
	var total int64

	// Contar total
	config.DB.Model(&models.Punishment{}).Count(&total)

	// Obtener castigos
	offset := (page - 1) * limit
	err := config.DB.Preload("Target").Preload("Assigner").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&punishments).Error

	return punishments, total, err
}

// Listar castigos por objetivo (target)
func (dao *PunishmentDAO) FindByTarget(targetID uint, page, limit int) ([]models.Punishment, int64, error) {
	var punishments []models.Punishment
	var total int64

	query := config.DB.Where("target_id = ?", targetID)

	// Contar total
	query.Model(&models.Punishment{}).Count(&total)

	// Obtener castigos
	offset := (page - 1) * limit
	err := query.Preload("Target").Preload("Assigner").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&punishments).Error

	return punishments, total, err
}

// Listar castigos por quien los asignó
func (dao *PunishmentDAO) FindByAssigner(assignerID uint, page, limit int) ([]models.Punishment, int64, error) {
	var punishments []models.Punishment
	var total int64

	query := config.DB.Where("assigned_by = ?", assignerID)

	// Contar total
	query.Model(&models.Punishment{}).Count(&total)

	// Obtener castigos
	offset := (page - 1) * limit
	err := query.Preload("Target").Preload("Assigner").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&punishments).Error

	return punishments, total, err
}

// Listar castigos por status
func (dao *PunishmentDAO) FindByStatus(status string, page, limit int) ([]models.Punishment, int64, error) {
	var punishments []models.Punishment
	var total int64

	query := config.DB.Where("status = ?", status)

	// Contar total
	query.Model(&models.Punishment{}).Count(&total)

	// Obtener castigos
	offset := (page - 1) * limit
	err := query.Preload("Target").Preload("Assigner").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&punishments).Error

	return punishments, total, err
}

// Actualizar castigo
func (dao *PunishmentDAO) Update(punishment *models.Punishment) error {
	return config.DB.Save(punishment).Error
}

// Eliminar castigo
func (dao *PunishmentDAO) Delete(id uint) error {
	return config.DB.Delete(&models.Punishment{}, id).Error
}

// Contar castigos por status
func (dao *PunishmentDAO) CountByStatus(status string) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Punishment{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// Castigos activos para un usuario
func (dao *PunishmentDAO) FindActiveByTarget(targetID uint) ([]models.Punishment, error) {
	var punishments []models.Punishment
	err := config.DB.Where("target_id = ? AND status = ?", targetID, models.PunishmentStatusActive).
		Preload("Assigner").
		Find(&punishments).Error
	return punishments, err
}
