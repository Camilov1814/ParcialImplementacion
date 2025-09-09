package main

import (
	"devops-chaos-backend/internal/config"
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
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Configurar Gin
	r := gin.Default()

	// Middleware básico
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Rutas básicas
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Andrei's Chaos Backend is running!",
		})
	})

	// Iniciar servidor
	port := "8080"
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
