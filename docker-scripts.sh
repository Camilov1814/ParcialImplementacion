#!/bin/bash

# DevOps Chaos - Docker Management Scripts

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Help function
show_help() {
    echo "DevOps Chaos - Docker Management Scripts"
    echo ""
    echo "Usage: ./docker-scripts.sh [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  build-prod       Build production images"
    echo "  build-dev        Build development images"
    echo "  start-prod       Start production environment"
    echo "  start-dev        Start development environment"
    echo "  stop             Stop all services"
    echo "  restart          Restart all services"
    echo "  logs             Show logs from all services"
    echo "  logs-backend     Show backend logs"
    echo "  logs-frontend    Show frontend logs"
    echo "  logs-db          Show database logs"
    echo "  clean            Clean up containers and images"
    echo "  clean-all        Clean everything (containers, images, volumes)"
    echo "  db-backup        Backup database"
    echo "  db-restore       Restore database from backup"
    echo "  status           Show status of all services"
    echo "  shell-backend    Open shell in backend container"
    echo "  shell-db         Open database shell"
    echo "  help             Show this help message"
}

# Build production images
build_prod() {
    print_status "Building production images..."
    docker-compose build --no-cache
    if [ $? -eq 0 ]; then
        print_success "Production images built successfully"
    else
        print_error "Failed to build production images"
        exit 1
    fi
}

# Build development images
build_dev() {
    print_status "Building development images..."
    docker-compose -f docker-compose.dev.yml build --no-cache
    if [ $? -eq 0 ]; then
        print_success "Development images built successfully"
    else
        print_error "Failed to build development images"
        exit 1
    fi
}

# Start production environment
start_prod() {
    print_status "Starting production environment..."
    docker-compose up -d
    if [ $? -eq 0 ]; then
        print_success "Production environment started"
        print_status "Frontend: http://localhost:3000"
        print_status "Backend: http://localhost:8080"
        print_status "Database: localhost:5432"
    else
        print_error "Failed to start production environment"
        exit 1
    fi
}

# Start development environment
start_dev() {
    print_status "Starting development environment..."
    docker-compose -f docker-compose.dev.yml up -d
    if [ $? -eq 0 ]; then
        print_success "Development environment started"
        print_status "Frontend: http://localhost:5173 (with hot reload)"
        print_status "Backend: http://localhost:8081"
        print_status "Database: localhost:5433"
    else
        print_error "Failed to start development environment"
        exit 1
    fi
}

# Stop all services
stop_services() {
    print_status "Stopping all services..."
    docker-compose down
    docker-compose -f docker-compose.dev.yml down
    print_success "All services stopped"
}

# Restart services
restart_services() {
    print_status "Restarting services..."
    stop_services
    start_prod
}

# Show logs
show_logs() {
    print_status "Showing logs from all services..."
    docker-compose logs -f
}

# Show backend logs
show_backend_logs() {
    print_status "Showing backend logs..."
    docker-compose logs -f backend
}

# Show frontend logs
show_frontend_logs() {
    print_status "Showing frontend logs..."
    docker-compose logs -f frontend
}

# Show database logs
show_db_logs() {
    print_status "Showing database logs..."
    docker-compose logs -f postgres
}

# Clean up
clean_up() {
    print_status "Cleaning up containers and images..."
    docker-compose down --rmi local
    docker system prune -f
    print_success "Cleanup completed"
}

# Clean everything
clean_all() {
    print_warning "This will remove ALL containers, images, and volumes!"
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Removing everything..."
        docker-compose down -v --rmi all
        docker system prune -a -f --volumes
        print_success "Everything cleaned"
    else
        print_status "Cleanup cancelled"
    fi
}

# Backup database
backup_db() {
    print_status "Creating database backup..."
    timestamp=$(date +"%Y%m%d_%H%M%S")
    backup_file="backup_${timestamp}.sql"

    docker exec devops-chaos-db pg_dump -U postgres devops_chaos > "$backup_file"
    if [ $? -eq 0 ]; then
        print_success "Database backed up to $backup_file"
    else
        print_error "Failed to create database backup"
        exit 1
    fi
}

# Restore database
restore_db() {
    if [ -z "$2" ]; then
        print_error "Please specify backup file: ./docker-scripts.sh db-restore backup_file.sql"
        exit 1
    fi

    backup_file="$2"
    if [ ! -f "$backup_file" ]; then
        print_error "Backup file $backup_file not found"
        exit 1
    fi

    print_status "Restoring database from $backup_file..."
    docker exec -i devops-chaos-db psql -U postgres devops_chaos < "$backup_file"
    if [ $? -eq 0 ]; then
        print_success "Database restored successfully"
    else
        print_error "Failed to restore database"
        exit 1
    fi
}

# Show status
show_status() {
    print_status "Service Status:"
    docker-compose ps
    echo ""
    print_status "Container Health:"
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
}

# Open backend shell
backend_shell() {
    print_status "Opening shell in backend container..."
    docker exec -it devops-chaos-backend sh
}

# Open database shell
db_shell() {
    print_status "Opening database shell..."
    docker exec -it devops-chaos-db psql -U postgres devops_chaos
}

# Main script logic
case "$1" in
    "build-prod")
        build_prod
        ;;
    "build-dev")
        build_dev
        ;;
    "start-prod")
        start_prod
        ;;
    "start-dev")
        start_dev
        ;;
    "stop")
        stop_services
        ;;
    "restart")
        restart_services
        ;;
    "logs")
        show_logs
        ;;
    "logs-backend")
        show_backend_logs
        ;;
    "logs-frontend")
        show_frontend_logs
        ;;
    "logs-db")
        show_db_logs
        ;;
    "clean")
        clean_up
        ;;
    "clean-all")
        clean_all
        ;;
    "db-backup")
        backup_db
        ;;
    "db-restore")
        restore_db "$@"
        ;;
    "status")
        show_status
        ;;
    "shell-backend")
        backend_shell
        ;;
    "shell-db")
        db_shell
        ;;
    "help"|*)
        show_help
        ;;
esac