package main

import (
	"devops-chaos-backend/internal/config"
	"devops-chaos-backend/internal/models"
	"devops-chaos-backend/internal/seeders"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

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

	// Ejecutar seeders
	seeder := seeders.NewSeeder()
	if err := seeder.RunAllSeeders(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	log.Println("Seeding completed successfully!")
}
