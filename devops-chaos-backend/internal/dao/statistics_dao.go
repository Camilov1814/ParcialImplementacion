package dao

import (
	"devops-chaos-backend/internal/config"
	"devops-chaos-backend/internal/models"

	"gorm.io/gorm"
)

type StatisticDAO struct{}

func NewStatisticDAO() *StatisticDAO {
	return &StatisticDAO{}
}

// Crear estadística
func (dao *StatisticDAO) Create(statistic *models.Statistic) error {
	return config.DB.Create(statistic).Error
}

// Buscar estadística por UserID
func (dao *StatisticDAO) FindByUserID(userID uint) (*models.Statistic, error) {
	var statistic models.Statistic
	err := config.DB.Where("user_id = ?", userID).First(&statistic).Error
	if err != nil {
		return nil, err
	}
	return &statistic, nil
}

// Buscar estadística por UserID con información del usuario
func (dao *StatisticDAO) FindByUserIDWithUser(userID uint) (*models.Statistic, error) {
	var statistic models.Statistic
	err := config.DB.Preload("User").Where("user_id = ?", userID).First(&statistic).Error
	if err != nil {
		return nil, err
	}
	return &statistic, nil
}

// Actualizar estadística
func (dao *StatisticDAO) Update(statistic *models.Statistic) error {
	return config.DB.Save(statistic).Error
}

// Incrementar capturas
func (dao *StatisticDAO) IncrementCaptures(userID uint) error {
	return config.DB.Model(&models.Statistic{}).
		Where("user_id = ?", userID).
		UpdateColumn("captures_count", gorm.Expr("captures_count + ?", 1)).Error
}

// Incrementar reportes
func (dao *StatisticDAO) IncrementReports(userID uint) error {
	return config.DB.Model(&models.Statistic{}).
		Where("user_id = ?", userID).
		UpdateColumn("reports_count", gorm.Expr("reports_count + ?", 1)).Error
}

// Agregar puntos
func (dao *StatisticDAO) AddPoints(userID uint, points int) error {
	return config.DB.Model(&models.Statistic{}).
		Where("user_id = ?", userID).
		UpdateColumn("points", gorm.Expr("points + ?", points)).Error
}

// Top daemons por ranking
func (dao *StatisticDAO) FindTopDaemons(limit int) ([]models.Statistic, error) {
	var statistics []models.Statistic
	err := config.DB.Preload("User").
		Joins("JOIN users ON users.id = statistics.user_id").
		Where("users.role = ?", models.RoleDaemon).
		Order("statistics.points DESC, statistics.captures_count DESC").
		Limit(limit).
		Find(&statistics).Error
	return statistics, err
}

// Listar todas las estadísticas con paginación
func (dao *StatisticDAO) FindAll(page, limit int) ([]models.Statistic, int64, error) {
	var statistics []models.Statistic
	var total int64

	// Contar total
	config.DB.Model(&models.Statistic{}).Count(&total)

	// Obtener estadísticas
	offset := (page - 1) * limit
	err := config.DB.Preload("User").
		Order("points DESC").
		Offset(offset).Limit(limit).
		Find(&statistics).Error

	return statistics, total, err
}

// Crear o actualizar estadística
func (dao *StatisticDAO) CreateOrUpdate(userID uint) (*models.Statistic, error) {
	var statistic models.Statistic

	// Intentar encontrar existente
	err := config.DB.Where("user_id = ?", userID).First(&statistic).Error

	if err != nil {
		// No existe, crear nueva
		statistic = models.Statistic{
			UserID:        userID,
			CapturesCount: 0,
			ReportsCount:  0,
			Ranking:       0,
			Points:        0,
		}
		err = dao.Create(&statistic)
		if err != nil {
			return nil, err
		}
	}

	return &statistic, nil
}

// Recalcular rankings de todos los daemons
func (dao *StatisticDAO) RecalculateRankings() error {
	var statistics []models.Statistic

	// Obtener todas las estadísticas ordenadas por puntos
	err := config.DB.Joins("JOIN users ON users.id = statistics.user_id").
		Where("users.role = ?", models.RoleDaemon).
		Order("statistics.points DESC, statistics.captures_count DESC").
		Find(&statistics).Error

	if err != nil {
		return err
	}

	// Actualizar rankings
	for i, stat := range statistics {
		stat.Ranking = i + 1
		config.DB.Save(&stat)
	}

	return nil
}
