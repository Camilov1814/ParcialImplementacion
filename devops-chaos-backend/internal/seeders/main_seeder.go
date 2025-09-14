package seeders

import (
	"devops-chaos-backend/internal/config"
	"devops-chaos-backend/internal/models"
	"devops-chaos-backend/internal/utils"
	"log"
	"time"
)

type Seeder struct{}

func NewSeeder() *Seeder {
	return &Seeder{}
}

// RunAllSeeders ejecuta todos los seeders
func (s *Seeder) RunAllSeeders() error {
	log.Println("Starting database seeding...")

	if err := s.SeedUsers(); err != nil {
		return err
	}

	if err := s.SeedReports(); err != nil {
		return err
	}

	if err := s.SeedStatistics(); err != nil {
		return err
	}

	if err := s.SeedPunishments(); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

// SeedUsers crea usuarios iniciales
func (s *Seeder) SeedUsers() error {
	log.Println("Seeding users...")

	users := []models.User{
		{
			Username: "andrei",
			Email:    "andrei@chaos.dev",
			Password: s.hashPassword("AndreI2024!"),
			Role:     models.RoleAndrei,
			Status:   models.StatusActive,
		},
		// Daemons
		{
			Username: "daemon_alpha",
			Email:    "alpha@daemon.chaos",
			Password: s.hashPassword("DaemonAlpha123!"),
			Role:     models.RoleDaemon,
			Status:   models.StatusActive,
		},
		{
			Username: "daemon_beta",
			Email:    "beta@daemon.chaos",
			Password: s.hashPassword("DaemonBeta123!"),
			Role:     models.RoleDaemon,
			Status:   models.StatusActive,
		},
		{
			Username: "daemon_gamma",
			Email:    "gamma@daemon.chaos",
			Password: s.hashPassword("DaemonGamma123!"),
			Role:     models.RoleDaemon,
			Status:   models.StatusPunished,
		},
		{
			Username: "daemon_delta",
			Email:    "delta@daemon.chaos",
			Password: s.hashPassword("DaemonDelta123!"),
			Role:     models.RoleDaemon,
			Status:   models.StatusActive,
		},
		{
			Username: "daemon_omega",
			Email:    "omega@daemon.chaos",
			Password: s.hashPassword("DaemonOmega123!"),
			Role:     models.RoleDaemon,
			Status:   models.StatusActive,
		},
		// Network Admins
		{
			Username: "alice",
			Email:    "alice@resistance.net",
			Password: s.hashPassword("Alice123!"),
			Role:     models.RoleNetworkAdmin,
			Status:   models.StatusActive,
		},
		{
			Username: "bob",
			Email:    "bob@resistance.net",
			Password: s.hashPassword("Bob123!"),
			Role:     models.RoleNetworkAdmin,
			Status:   models.StatusActive,
		},
		{
			Username: "chalier",
			Email:    "chalier@resistance.net",
			Password: s.hashPassword("Chalier123!"),
			Role:     models.RoleNetworkAdmin,
			Status:   models.StatusActive,
		},
		{
			Username: "diana",
			Email:    "diana@resistance.net",
			Password: s.hashPassword("Diana123!"),
			Role:     models.RoleNetworkAdmin,
			Status:   models.StatusActive,
		},
	}

	for _, user := range users {
		// Verificar si el usuario ya existe
		var existingUser models.User
		result := config.DB.Where("username = ?", user.Username).First(&existingUser)
		if result.RowsAffected == 0 {
			if err := config.DB.Create(&user).Error; err != nil {
				log.Printf("Error creating user %s: %v", user.Username, err)
				return err
			}
			log.Printf("Created user: %s (%s)", user.Username, user.Role)
		} else {
			log.Printf("User %s already exists, skipping...", user.Username)
		}
	}

	return nil
}

// SeedReports crea reportes de ejemplo
func (s *Seeder) SeedReports() error {
	log.Println("Seeding reports...")

	// Obtener usuarios para referenciar
	var users []models.User
	config.DB.Find(&users)

	userMap := make(map[string]uint)
	for _, user := range users {
		userMap[user.Username] = user.ID
	}

	reports := []models.Report{
		{
			Title:       "Successful Server Infiltration",
			Description: "Successfully gained root access to main database server through SQL injection vulnerability in login form. Extracted user credentials and system configurations.",
			Type:        models.ReportTypeCapture,
			AuthorID:    &[]uint{userMap["daemon_alpha"]}[0],
			Status:      models.ReportStatusApproved,
			CreatedAt:   time.Now().Add(-72 * time.Hour),
		},
		{
			Title:       "Network Admin Behavioral Analysis",
			Description: "Observed suspicious activity from network admin 'alice'. She's been working late hours and accessing unusual network segments. Recommend increased surveillance.",
			Type:        models.ReportTypeCapture,
			AuthorID:    &[]uint{userMap["daemon_beta"]}[0],
			Status:      models.ReportStatusPending,
			CreatedAt:   time.Now().Add(-48 * time.Hour),
		},
		{
			Title:       "Social Engineering Success",
			Description: "Successfully extracted VPN credentials from junior network admin through phishing campaign. Target provided credentials thinking it was legitimate IT security test.",
			Type:        models.ReportTypeCapture,
			AuthorID:    &[]uint{userMap["daemon_delta"]}[0],
			Status:      models.ReportStatusApproved,
			CreatedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			Title:       "Firewall Configuration Changes",
			Description: "Detected unusual modifications to firewall rules during off-hours. Rules were changed to allow suspicious outbound traffic on non-standard ports.",
			Type:        models.ReportTypeResistance,
			AuthorID:    &[]uint{userMap["daemon_alpha"]}[0],
			Status:      models.ReportStatusApproved,
			CreatedAt:   time.Now().Add(-36 * time.Hour),
		},
		{
			Title:       "Suspicious Daemon Activity Detected",
			Description: "Anonymous report: I've noticed someone accessing server rooms during night shifts when no maintenance is scheduled. They seem to know their way around our systems too well.",
			Type:        models.ReportTypeAnonymous,
			AuthorID:    nil, // Reporte anónimo
			Status:      models.ReportStatusPending,
			CreatedAt:   time.Now().Add(-12 * time.Hour),
		},
		{
			Title:       "Potential Network Compromise",
			Description: "Anonymous report: Our internal monitoring shows data exfiltration during off-peak hours. The traffic patterns don't match normal business operations.",
			Type:        models.ReportTypeAnonymous,
			AuthorID:    nil,
			Status:      models.ReportStatusPending,
			CreatedAt:   time.Now().Add(-6 * time.Hour),
		},
		{
			Title:       "Physical Security Breach",
			Description: "Successfully bypassed keycard security and gained physical access to main server room. Planted network monitoring devices on critical infrastructure.",
			Type:        models.ReportTypeCapture,
			AuthorID:    &[]uint{userMap["daemon_omega"]}[0],
			Status:      models.ReportStatusApproved,
			CreatedAt:   time.Now().Add(-96 * time.Hour),
		},
		{
			Title:       "Emergency: Admin Credentials Compromised",
			Description: "Anonymous report: I believe my credentials have been compromised. I've seen login attempts from locations I've never been to. Please investigate immediately.",
			Type:        models.ReportTypeAnonymous,
			AuthorID:    nil,
			Status:      models.ReportStatusApproved,
			CreatedAt:   time.Now().Add(-18 * time.Hour),
		},
	}

	for _, report := range reports {
		var existingReport models.Report
		result := config.DB.Where("title = ?", report.Title).First(&existingReport)
		if result.RowsAffected == 0 {
			if err := config.DB.Create(&report).Error; err != nil {
				log.Printf("Error creating report %s: %v", report.Title, err)
				return err
			}
			log.Printf("Created report: %s", report.Title)
		}
	}

	return nil
}

// SeedStatistics crea estadísticas para los daemons
func (s *Seeder) SeedStatistics() error {
	log.Println("Seeding statistics...")

	// Obtener todos los daemons
	var daemons []models.User
	config.DB.Where("role = ?", models.RoleDaemon).Find(&daemons)

	statisticsData := map[string]struct {
		captures int
		reports  int
		points   int
	}{
		"daemon_alpha": {captures: 8, reports: 12, points: 120},
		"daemon_beta":  {captures: 5, reports: 8, points: 85},
		"daemon_gamma": {captures: 3, reports: 4, points: 45}, // Este está castigado
		"daemon_delta": {captures: 6, reports: 9, points: 95},
		"daemon_omega": {captures: 10, reports: 15, points: 150}, // El mejor
	}

	for _, daemon := range daemons {
		var existingStats models.Statistic
		result := config.DB.Where("user_id = ?", daemon.ID).First(&existingStats)

		if result.RowsAffected == 0 {
			data := statisticsData[daemon.Username]
			stats := models.Statistic{
				UserID:        daemon.ID,
				CapturesCount: data.captures,
				ReportsCount:  data.reports,
				Points:        data.points,
				Ranking:       0, // Se calculará después
			}

			if err := config.DB.Create(&stats).Error; err != nil {
				log.Printf("Error creating statistics for %s: %v", daemon.Username, err)
				return err
			}
			log.Printf("Created statistics for: %s", daemon.Username)
		}
	}

	// Recalcular rankings
	s.recalculateRankings()

	return nil
}

// SeedPunishments crea castigos de ejemplo
func (s *Seeder) SeedPunishments() error {
	log.Println("Seeding punishments...")

	// Obtener IDs necesarios
	var andrei, gammaUser models.User
	config.DB.Where("username = ?", "andrei").First(&andrei)
	config.DB.Where("username = ?", "daemon_gamma").First(&gammaUser)

	expiresAt := time.Now().Add(48 * time.Hour)
	punishments := []models.Punishment{
		{
			TargetID:    gammaUser.ID,
			AssignedBy:  andrei.ID,
			Type:        models.PunishmentTypeTimeout,
			Description: "Excessive system access during unauthorized hours. 48-hour network access suspension for security review.",
			Status:      models.PunishmentStatusActive,
			ExpiresAt:   &expiresAt,
			CreatedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			TargetID:    gammaUser.ID,
			AssignedBy:  andrei.ID,
			Type:        models.PunishmentTypeExtraTasks,
			Description: "Additional security protocol documentation required. Must complete comprehensive security audit report within 72 hours.",
			Status:      models.PunishmentStatusActive,
			CreatedAt:   time.Now().Add(-12 * time.Hour),
		},
	}

	for _, punishment := range punishments {
		var existingPunishment models.Punishment
		result := config.DB.Where("target_id = ? AND description = ?", punishment.TargetID, punishment.Description).First(&existingPunishment)
		if result.RowsAffected == 0 {
			if err := config.DB.Create(&punishment).Error; err != nil {
				log.Printf("Error creating punishment: %v", err)
				return err
			}
			log.Printf("Created punishment for daemon_gamma: %s", punishment.Type)
		}
	}

	return nil
}

// Helper methods

func (s *Seeder) hashPassword(password string) string {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}
	return hashed
}

func (s *Seeder) recalculateRankings() {
	var statistics []models.Statistic
	config.DB.Joins("JOIN users ON users.id = statistics.user_id").
		Where("users.role = ?", models.RoleDaemon).
		Order("statistics.points DESC, statistics.captures_count DESC").
		Find(&statistics)

	for i, stat := range statistics {
		stat.Ranking = i + 1
		config.DB.Save(&stat)
	}
}
