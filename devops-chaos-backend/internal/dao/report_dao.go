package dao

import (
	"devops-chaos-backend/internal/config"
	"devops-chaos-backend/internal/models"
)

type ReportDAO struct{}

func NewReportDAO() *ReportDAO {
	return &ReportDAO{}
}

// Crear reporte
func (dao *ReportDAO) Create(report *models.Report) error {
	return config.DB.Create(report).Error
}

// Buscar reporte por ID
func (dao *ReportDAO) FindByID(id uint) (*models.Report, error) {
	var report models.Report
	err := config.DB.Preload("Author").First(&report, id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// Listar todos los reportes con paginación
func (dao *ReportDAO) FindAll(page, limit int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	// Contar total
	config.DB.Model(&models.Report{}).Count(&total)

	// Obtener reportes con autor
	offset := (page - 1) * limit
	err := config.DB.Preload("Author").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&reports).Error

	return reports, total, err
}

// Listar reportes por autor
func (dao *ReportDAO) FindByAuthor(authorID uint, page, limit int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	query := config.DB.Where("author_id = ?", authorID)

	// Contar total
	query.Model(&models.Report{}).Count(&total)

	// Obtener reportes
	offset := (page - 1) * limit
	err := query.Preload("Author").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&reports).Error

	return reports, total, err
}

// Listar reportes por tipo
func (dao *ReportDAO) FindByType(reportType string, page, limit int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	query := config.DB.Where("type = ?", reportType)

	// Contar total
	query.Model(&models.Report{}).Count(&total)

	// Obtener reportes
	offset := (page - 1) * limit
	err := query.Preload("Author").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&reports).Error

	return reports, total, err
}

// Listar reportes por status
func (dao *ReportDAO) FindByStatus(status string, page, limit int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	query := config.DB.Where("status = ?", status)

	// Contar total
	query.Model(&models.Report{}).Count(&total)

	// Obtener reportes
	offset := (page - 1) * limit
	err := query.Preload("Author").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&reports).Error

	return reports, total, err
}

// Listar reportes anónimos
func (dao *ReportDAO) FindAnonymous(page, limit int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	query := config.DB.Where("author_id IS NULL")

	// Contar total
	query.Model(&models.Report{}).Count(&total)

	// Obtener reportes
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&reports).Error

	return reports, total, err
}

// Actualizar reporte
func (dao *ReportDAO) Update(report *models.Report) error {
	return config.DB.Save(report).Error
}

// Eliminar reporte
func (dao *ReportDAO) Delete(id uint) error {
	return config.DB.Delete(&models.Report{}, id).Error
}

// Contar reportes por status
func (dao *ReportDAO) CountByStatus(status string) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Report{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// Contar reportes por autor
func (dao *ReportDAO) CountByAuthor(authorID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Report{}).Where("author_id = ?", authorID).Count(&count).Error
	return count, err
}

// Reportes recientes para dashboard
func (dao *ReportDAO) FindRecent(limit int) ([]models.Report, error) {
	var reports []models.Report
	err := config.DB.Preload("Author").
		Order("created_at DESC").
		Limit(limit).
		Find(&reports).Error
	return reports, err
}
