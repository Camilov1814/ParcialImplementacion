# Variables
DB_HOST := localhost
DB_PORT := 5432
DB_USER := chaos_user
DB_PASSWORD := chaos_password
DB_NAME := devops_chaos

.PHONY: help dev build seed clean-db reset-db

help: ## Mostrar ayuda
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

dev: ## Ejecutar en modo desarrollo
	go run cmd/main.go

build: ## Construir la aplicación
	go build -o bin/devops-chaos-backend cmd/main.go

seed: ## Ejecutar seeders de la base de datos
	@echo "🌱 Ejecutando seeders..."
	go run cmd/seeder/main.go

clean-db: ## Limpiar la base de datos
	@echo "🧹 Limpiando base de datos..."
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -c "TRUNCATE TABLE users, reports, punishments, statistics RESTART IDENTITY CASCADE;"

reset-db: clean-db seed ## Limpiar y re-popular la base de datos
	@echo "🔄 Base de datos reseteada con datos iniciales"

test: ## Ejecutar tests
	go test -v ./...

format: ## Formatear código
	go fmt ./...

tidy: ## Limpiar módulos
	go mod tidy

install-deps: ## Instalar dependencias
	go mod download