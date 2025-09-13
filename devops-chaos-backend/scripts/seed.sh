#!/bin/bash

echo "🌱 Starting database seeding..."

# Verificar que el archivo .env existe
if [ ! -f .env ]; then
    echo "❌ Error: .env file not found!"
    echo "Please create a .env file with database configuration"
    exit 1
fi

# Ejecutar el seeder
go run cmd/seeder/main.go

if [ $? -eq 0 ]; then
    echo "✅ Database seeding completed successfully!"
else
    echo "❌ Database seeding failed!"
    exit 1
fi