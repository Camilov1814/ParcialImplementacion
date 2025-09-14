# DevOps Chaos - Docker Setup

Este proyecto ha sido completamente dockerizado para facilitar el despliegue y desarrollo.

## Estructura del Proyecto

```
.
├── devops-chaos-backend/     # Backend en Go
│   ├── Dockerfile           # Dockerfile para el backend
│   └── .dockerignore       # Archivos a ignorar en Docker
├── devops-chaos-frontend/   # Frontend en React
│   ├── Dockerfile          # Dockerfile para el frontend
│   ├── nginx.conf          # Configuración de Nginx
│   └── .dockerignore      # Archivos a ignorar en Docker
├── docker-compose.yml      # Orquestación para producción
├── docker-compose.dev.yml  # Orquestación para desarrollo
└── .env.example           # Variables de entorno de ejemplo
```

## Comandos Principales

### Producción

```bash
# Construir y levantar todos los servicios
docker-compose up --build

# Levantar en segundo plano
docker-compose up -d

# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down

# Detener y eliminar volúmenes
docker-compose down -v
```

### Desarrollo

```bash
# Usar el archivo de desarrollo
docker-compose -f docker-compose.dev.yml up --build

# Desarrollo con hot reload
docker-compose -f docker-compose.dev.yml up -d
```

## Servicios

### 🗄️ Base de Datos (PostgreSQL)
- **Puerto**: 5432 (producción), 5433 (desarrollo)
- **Usuario**: postgres
- **Contraseña**: password123 (producción), devpass123 (desarrollo)
- **Base de datos**: devops_chaos (producción), devops_chaos_dev (desarrollo)

### 🚀 Backend (Go)
- **Puerto**: 8080 (producción), 8081 (desarrollo)
- **Framework**: Gin
- **Endpoint**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

### 🎨 Frontend (React + Vite)
- **Puerto**: 3000 (producción con Nginx), 5173 (desarrollo con Vite)
- **Framework**: React con TypeScript
- **URL**: http://localhost:3000 (producción), http://localhost:5173 (desarrollo)

## Configuración

### Variables de Entorno

1. Copia el archivo de ejemplo:
   ```bash
   cp .env.example .env
   ```

2. Modifica las variables según tu entorno:
   ```bash
   # Para producción
   DATABASE_URL=postgres://postgres:password123@postgres:5432/devops_chaos?sslmode=disable
   JWT_SECRET=change_this_secret_key

   # Para desarrollo
   DATABASE_URL=postgres://postgres:devpass123@postgres:5432/devops_chaos_dev?sslmode=disable
   JWT_SECRET=dev_secret_key
   ```

### Primer Uso

1. **Instalar Docker y Docker Compose**
   - [Docker Desktop](https://www.docker.com/products/docker-desktop)

2. **Clonar el repositorio**
   ```bash
   git clone <repository-url>
   cd devops-chaos
   ```

3. **Configurar variables de entorno**
   ```bash
   cp .env.example .env
   # Edita el archivo .env con tus configuraciones
   ```

4. **Levantar los servicios**
   ```bash
   # Para desarrollo
   docker-compose -f docker-compose.dev.yml up --build

   # Para producción
   docker-compose up --build
   ```

## Comandos Útiles

### Docker

```bash
# Ver contenedores activos
docker ps

# Ver logs específicos
docker-compose logs backend
docker-compose logs frontend
docker-compose logs postgres

# Ejecutar comandos en contenedores
docker exec -it devops-chaos-backend sh
docker exec -it devops-chaos-db psql -U postgres -d devops_chaos

# Limpiar todo
docker-compose down -v --rmi all
docker system prune -a
```

### Base de Datos

```bash
# Conectar a la base de datos
docker exec -it devops-chaos-db psql -U postgres -d devops_chaos

# Backup de la base de datos
docker exec devops-chaos-db pg_dump -U postgres devops_chaos > backup.sql

# Restaurar backup
docker exec -i devops-chaos-db psql -U postgres devops_chaos < backup.sql
```

## Desarrollo

Para desarrollo local con hot reload:

```bash
# Levantar solo la base de datos
docker-compose up postgres -d

# O usar el compose de desarrollo
docker-compose -f docker-compose.dev.yml up
```

El frontend estará disponible en http://localhost:5173 con hot reload activado.
El backend estará en http://localhost:8081 en modo debug.

## Problemas Comunes

### Puerto ocupado
```bash
# Verificar qué proceso usa el puerto
netstat -tulpn | grep :8080

# Cambiar el puerto en docker-compose.yml
ports:
  - "8081:8080"  # Puerto host:Puerto contenedor
```

### Problemas de permisos
```bash
# Linux/Mac: Dar permisos
sudo chown -R $USER:$USER .

# Windows: Ejecutar Docker Desktop como administrador
```

### Base de datos no se conecta
```bash
# Verificar que el servicio esté corriendo
docker-compose ps

# Ver logs de la base de datos
docker-compose logs postgres

# Recrear la base de datos
docker-compose down -v
docker-compose up postgres -d
```

### Reconstruir imágenes
```bash
# Forzar reconstrucción
docker-compose build --no-cache

# Reconstruir servicio específico
docker-compose build backend --no-cache
```

## Despliegue en Producción

### Usando Docker Compose

```bash
# 1. Configurar variables de producción
export JWT_SECRET="your_production_secret"
export DATABASE_URL="your_production_db_url"

# 2. Levantar servicios
docker-compose up -d

# 3. Verificar estado
docker-compose ps
```

### Usando Docker Swarm

```bash
# Inicializar swarm
docker swarm init

# Desplegar stack
docker stack deploy -c docker-compose.yml devops-chaos

# Ver servicios
docker service ls
```

## Monitoreo

### Logs
```bash
# Ver todos los logs
docker-compose logs -f

# Logs específicos con filtro de tiempo
docker-compose logs --since="1h" backend
```

### Health Checks
Los servicios incluyen health checks que se pueden monitorear:

```bash
# Ver estado de salud
docker inspect devops-chaos-backend | grep -A 5 Health
```

## Licencia

Este proyecto está bajo la licencia MIT.