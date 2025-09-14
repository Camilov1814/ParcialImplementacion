package dao

import (
	"devops-chaos-backend/internal/config"
	"devops-chaos-backend/internal/models"
)

type CaptureDAO struct{}

func NewCaptureDAO() *CaptureDAO {
	return &CaptureDAO{}
}

// Crear una nueva captura
func (dao *CaptureDAO) Create(capture *models.Capture) error {
	db := config.GetDatabase()
	return db.Create(capture).Error
}

// Obtener capturas por daemon
func (dao *CaptureDAO) FindByDaemonID(daemonID uint) ([]models.Capture, error) {
	var captures []models.Capture
	db := config.GetDatabase()

	err := db.Preload("Daemon").Preload("Target").
		Where("daemon_id = ?", daemonID).
		Order("capture_date DESC").
		Find(&captures).Error

	return captures, err
}

// Obtener todas las capturas (para Andrei)
func (dao *CaptureDAO) FindAll() ([]models.Capture, error) {
	var captures []models.Capture
	db := config.GetDatabase()

	err := db.Preload("Daemon").Preload("Target").
		Order("capture_date DESC").
		Find(&captures).Error

	return captures, err
}

// Obtener capturas por target (network admin capturado)
func (dao *CaptureDAO) FindByTargetID(targetID uint) ([]models.Capture, error) {
	var captures []models.Capture
	db := config.GetDatabase()

	err := db.Preload("Daemon").Preload("Target").
		Where("target_id = ?", targetID).
		Order("capture_date DESC").
		Find(&captures).Error

	return captures, err
}

// Contar capturas por daemon
func (dao *CaptureDAO) CountByDaemonID(daemonID uint) (int64, error) {
	var count int64
	db := config.GetDatabase()

	err := db.Model(&models.Capture{}).
		Where("daemon_id = ? AND status = ?", daemonID, models.CaptureStatusCaptured).
		Count(&count).Error

	return count, err
}

// Obtener capturas recientes
func (dao *CaptureDAO) FindRecent(limit int) ([]models.Capture, error) {
	var captures []models.Capture
	db := config.GetDatabase()

	err := db.Preload("Daemon").Preload("Target").
		Where("status = ?", models.CaptureStatusCaptured).
		Order("capture_date DESC").
		Limit(limit).
		Find(&captures).Error

	return captures, err
}

// Actualizar status de captura
func (dao *CaptureDAO) UpdateStatus(captureID uint, status string) error {
	db := config.GetDatabase()

	return db.Model(&models.Capture{}).
		Where("id = ?", captureID).
		Update("status", status).Error
}