package services

import (
	"devops-chaos-backend/internal/dao"
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/models"
	"errors"
	"math"
)

type ReportService struct {
	reportDAO    *dao.ReportDAO
	statisticDAO *dao.StatisticDAO
}

func NewReportService() *ReportService {
	return &ReportService{
		reportDAO:    dao.NewReportDAO(),
		statisticDAO: dao.NewStatisticDAO(),
	}
}

// Crear reporte
func (s *ReportService) CreateReport(req *dto.CreateReportRequest, userID uint, userRole string) (*dto.ReportResponse, error) {
	// Validar tipo de reporte según rol
	if !s.canCreateReportType(req.Type, userRole) {
		return nil, errors.New("unauthorized to create this type of report")
	}

	var authorID *uint
	if req.Type != models.ReportTypeAnonymous {
		authorID = &userID
	}

	report := &models.Report{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		AuthorID:    authorID,
		Status:      models.ReportStatusPending,
	}

	err := s.reportDAO.Create(report)
	if err != nil {
		return nil, err
	}

	// Incrementar contador de reportes si no es anónimo
	if authorID != nil {
		s.statisticDAO.IncrementReports(*authorID)
	}

	// Cargar el reporte con autor
	createdReport, err := s.reportDAO.FindByID(report.ID)
	if err != nil {
		return nil, err
	}

	return s.convertToReportResponse(createdReport), nil
}

// Obtener reporte por ID
func (s *ReportService) GetReportByID(id uint, userID uint, userRole string) (*dto.ReportResponse, error) {
	report, err := s.reportDAO.FindByID(id)
	if err != nil {
		return nil, errors.New("report not found")
	}

	// Control de acceso
	if !s.canAccessReport(report, userID, userRole) {
		return nil, errors.New("unauthorized to access this report")
	}

	return s.convertToReportResponse(report), nil
}

// Listar reportes con paginación
func (s *ReportService) GetReports(page, limit int, reportType, status string, userID uint, userRole string) (*dto.PaginatedResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var reports []models.Report
	var total int64
	var err error

	switch userRole {
	case models.RoleAndrei:
		// Andrei puede ver todos los reportes
		if reportType != "" && status != "" {
			// Filtrar por ambos
			reports, total, err = s.getReportsByTypeAndStatus(reportType, status, page, limit)
		} else if reportType != "" {
			reports, total, err = s.reportDAO.FindByType(reportType, page, limit)
		} else if status != "" {
			reports, total, err = s.reportDAO.FindByStatus(status, page, limit)
		} else {
			reports, total, err = s.reportDAO.FindAll(page, limit)
		}

	case models.RoleDaemon:
		// Daemon solo ve sus propios reportes
		reports, total, err = s.reportDAO.FindByAuthor(userID, page, limit)

	case models.RoleNetworkAdmin:
		// Network admin solo ve reportes anónimos (los suyos)
		reports, total, err = s.reportDAO.FindAnonymous(page, limit)

	default:
		return nil, errors.New("unauthorized")
	}

	if err != nil {
		return nil, err
	}

	// Convertir a DTO
	reportItems := make([]dto.ReportListItem, len(reports))
	for i, report := range reports {
		authorName := ""
		if report.Author != nil {
			authorName = report.Author.Username
		}

		reportItems[i] = dto.ReportListItem{
			ID:        report.ID,
			Title:     report.Title,
			Type:      report.Type,
			Status:    report.Status,
			AuthorID:  report.AuthorID,
			Author:    authorName,
			CreatedAt: report.CreatedAt,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.PaginatedResponse{
		Success:    true,
		Message:    "Reports retrieved successfully",
		Data:       reportItems,
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

// Actualizar status de reporte (solo Andrei)
func (s *ReportService) UpdateReportStatus(reportID uint, req *dto.UpdateReportRequest, userRole string) error {
	if userRole != models.RoleAndrei {
		return errors.New("unauthorized: only Andrei can update report status")
	}

	report, err := s.reportDAO.FindByID(reportID)
	if err != nil {
		return errors.New("report not found")
	}

	// Validar status
	if !s.isValidStatus(req.Status) {
		return errors.New("invalid status")
	}

	oldStatus := report.Status
	report.Status = req.Status

	err = s.reportDAO.Update(report)
	if err != nil {
		return err
	}

	// Si se aprueba un reporte, dar puntos al autor
	if oldStatus == models.ReportStatusPending && req.Status == models.ReportStatusApproved && report.AuthorID != nil {
		s.statisticDAO.AddPoints(*report.AuthorID, 5) // 5 puntos por reporte aprobado
	}

	return nil
}

// Eliminar reporte
func (s *ReportService) DeleteReport(reportID uint, userID uint, userRole string) error {
	report, err := s.reportDAO.FindByID(reportID)
	if err != nil {
		return errors.New("report not found")
	}

	// Solo Andrei o el autor pueden eliminar
	if userRole != models.RoleAndrei && (report.AuthorID == nil || *report.AuthorID != userID) {
		return errors.New("unauthorized to delete this report")
	}

	return s.reportDAO.Delete(reportID)
}

// Obtener reportes recientes para dashboard
func (s *ReportService) GetRecentReports(limit int, userRole string) ([]dto.ReportListItem, error) {
	if userRole != models.RoleAndrei {
		return nil, errors.New("unauthorized")
	}

	reports, err := s.reportDAO.FindRecent(limit)
	if err != nil {
		return nil, err
	}

	reportItems := make([]dto.ReportListItem, len(reports))
	for i, report := range reports {
		authorName := ""
		if report.Author != nil {
			authorName = report.Author.Username
		}

		reportItems[i] = dto.ReportListItem{
			ID:        report.ID,
			Title:     report.Title,
			Type:      report.Type,
			Status:    report.Status,
			AuthorID:  report.AuthorID,
			Author:    authorName,
			CreatedAt: report.CreatedAt,
		}
	}

	return reportItems, nil
}

// Helper methods

func (s *ReportService) canCreateReportType(reportType, userRole string) bool {
	switch userRole {
	case models.RoleDaemon:
		return reportType == models.ReportTypeCapture || reportType == models.ReportTypeResistance
	case models.RoleNetworkAdmin:
		return reportType == models.ReportTypeAnonymous
	case models.RoleAndrei:
		return true // Andrei puede crear cualquier tipo
	default:
		return false
	}
}

func (s *ReportService) canAccessReport(report *models.Report, userID uint, userRole string) bool {
	switch userRole {
	case models.RoleAndrei:
		return true // Andrei puede ver todo
	case models.RoleDaemon:
		return report.AuthorID != nil && *report.AuthorID == userID // Solo sus reportes
	case models.RoleNetworkAdmin:
		return report.Type == models.ReportTypeAnonymous // Solo anónimos
	default:
		return false
	}
}

func (s *ReportService) isValidStatus(status string) bool {
	validStatuses := []string{
		models.ReportStatusPending,
		models.ReportStatusApproved,
		models.ReportStatusRejected,
	}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func (s *ReportService) convertToReportResponse(report *models.Report) *dto.ReportResponse {
	response := &dto.ReportResponse{
		ID:          report.ID,
		Title:       report.Title,
		Description: report.Description,
		Type:        report.Type,
		Status:      report.Status,
		CreatedAt:   report.CreatedAt,
		UpdatedAt:   report.UpdatedAt,
	}

	if report.Author != nil {
		response.Author = &dto.UserInfo{
			ID:       report.Author.ID,
			Username: report.Author.Username,
			Email:    report.Author.Email,
			Role:     report.Author.Role,
			Status:   report.Author.Status,
		}
	}

	return response
}

func (s *ReportService) getReportsByTypeAndStatus(reportType, status string, page, limit int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	query := s.reportDAO.reportDAO.DB.Model(&models.Report{}).Where("type = ? AND status = ?", reportType, status)

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("Author").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&reports).Error

	return reports, total, err
}
