package main

import (
	"devops-chaos-backend/internal/config"
	"devops-chaos-backend/internal/controllers"
	"devops-chaos-backend/internal/middleware"
	"devops-chaos-backend/internal/models"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found: %v", err)
	}

	// DEBUG: Verificar que las variables se cargan
	log.Printf("DB_HOST: '%s'", os.Getenv("DB_HOST"))
	log.Printf("DB_USER: '%s'", os.Getenv("DB_USER"))
	log.Printf("DB_NAME: '%s'", os.Getenv("DB_NAME"))

	// Conectar a la base de datos
	config.ConnectDatabase()

	// Auto-migrar modelos
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Report{},
		&models.Punishment{},
		&models.Statistic{},
		&models.Capture{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Configurar Gin
	r := gin.Default()

	// Middleware global
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Inicializar controllers
	authController := controllers.NewAuthController()
	userController := controllers.NewUserController()
	reportController := controllers.NewReportController()
	punishmentController := controllers.NewPunishmentController()
	statisticsController := controllers.NewStatisticsController()
	dashboardController := controllers.NewDashboardController()
	resistanceController := controllers.NewResistanceController()

	// Rutas públicas
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Andrei's Chaos Backend is running!",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "healthy",
			"database": "connected",
			"version":  "1.0.0",
		})
	})

	// Rutas de autenticación (públicas)
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", authController.Login)
	}

	// Rutas protegidas
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())

	// Auth routes (protegidas)
	authProtected := api.Group("/auth")
	{
		authProtected.GET("/me", authController.GetProfile)
		authProtected.POST("/change-password", authController.ChangePassword)
		authProtected.POST("/register", middleware.AndreiOnlyMiddleware(), authController.Register)
	}

	// User routes
	users := api.Group("/users")
	{
		users.GET("", userController.GetUsers)
		users.GET("/:id", userController.GetUserByID)
		users.PUT("/:id", middleware.AndreiOnlyMiddleware(), userController.UpdateUser)
		users.DELETE("/:id", middleware.AndreiOnlyMiddleware(), userController.DeleteUser)
		users.GET("/stats", middleware.AndreiOnlyMiddleware(), userController.GetUserStats)
		users.POST("/:id/capture", middleware.DaemonOnlyMiddleware(), userController.CaptureNetworkAdmin)
		users.GET("/:id/captures", middleware.AndreiAndDaemonMiddleware(), userController.GetDaemonCaptures)
		users.GET("/:id/punishments/active", punishmentController.GetActivePunishments)
	}

	// Capture routes
	captures := api.Group("/captures")
	{
		captures.GET("/all", middleware.AndreiOnlyMiddleware(), userController.GetAllCaptures)
	}

	// Network Admin capture endpoint
	api.GET("/network-admins/capture", middleware.DaemonOnlyMiddleware(), userController.GetNetworkAdminsForCapture)

	// Report routes
	reports := api.Group("/reports")
	{
		reports.POST("", reportController.CreateReport)
		reports.GET("", reportController.GetReports)
		reports.GET("/:id", reportController.GetReportByID)
		reports.PUT("/:id/status", middleware.AndreiOnlyMiddleware(), reportController.UpdateReportStatus)
		reports.DELETE("/:id", reportController.DeleteReport)
		reports.GET("/recent", middleware.AndreiOnlyMiddleware(), reportController.GetRecentReports)
	}

	// Punishment routes
	punishments := api.Group("/punishments")
	{
		punishments.POST("", middleware.AndreiOnlyMiddleware(), punishmentController.CreatePunishment)
		punishments.GET("", punishmentController.GetPunishments)
		punishments.GET("/:id", punishmentController.GetPunishmentByID)
		punishments.PUT("/:id", middleware.AndreiOnlyMiddleware(), punishmentController.UpdatePunishment)
		punishments.DELETE("/:id", middleware.AndreiOnlyMiddleware(), punishmentController.DeletePunishment)
	}

	// Statistics routes
	statistics := api.Group("/statistics")
	{
		statistics.GET("/:user_id", statisticsController.GetUserStatistics)
		statistics.PUT("/:user_id", middleware.AndreiOnlyMiddleware(), statisticsController.UpdateStatistics)
		statistics.GET("/leaderboard", middleware.AndreiAndDaemonMiddleware(), statisticsController.GetLeaderboard)
		statistics.POST("/recalculate-rankings", middleware.AndreiOnlyMiddleware(), statisticsController.RecalculateRankings)
	}

	// Dashboard routes
	dashboard := api.Group("/dashboard")
	{
		dashboard.GET("/andrei", middleware.AndreiOnlyMiddleware(), dashboardController.GetAndreiDashboard)
		dashboard.GET("/daemon", middleware.DaemonOnlyMiddleware(), dashboardController.GetDaemonDashboard)
	}

	// Resistance routes (para Network Admins)
	resistance := api.Group("/resistance")
	{
		resistance.GET("", middleware.NetworkAdminOnlyMiddleware(), resistanceController.GetResistancePage)
		resistance.POST("/report", resistanceController.ReportSuspiciousActivity) // Anónimo
	}

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET  /ping")
	log.Printf("  GET  /health")
	log.Printf("  POST /api/auth/login")
	log.Printf("  GET  /api/auth/me")
	log.Printf("  GET  /api/users")
	log.Printf("  POST /api/users/:id/capture")
	log.Printf("  GET  /api/network-admins/capture")
	log.Printf("  GET  /api/reports")
	log.Printf("  GET  /api/dashboard/andrei")
	log.Printf("  GET  /api/dashboard/daemon")
	log.Printf("  ... and many more!")

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
