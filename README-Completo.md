# DevOps Chaos - Sistema Completo de Gestión de Roles Cyberpunk

## 🎯 Descripción del Proyecto

**DevOps Chaos** es una aplicación web full-stack que simula un mundo cyberpunk donde diferentes roles interactúan en un sistema de captura y resistencia. El proyecto combina una temática de hacking/terminal con funcionalidades reales de gestión de usuarios, autenticación JWT y control de acceso basado en roles.

### 🎭 Narrativa del Sistema
Una realidad distópica donde **Andrei** (líder supremo) controla un ejército de hackers (**daemons**) que buscan capturar administradores de red (**network admins**), mientras estos últimos forman una resistencia clandestina.

---

## 🏗️ Arquitectura del Sistema

### Stack Tecnológico

#### Backend (Go)
- **Framework**: Gin (HTTP web framework)
- **Base de Datos**: PostgreSQL con GORM ORM
- **Autenticación**: JWT con bcrypt
- **Contenorización**: Docker + Docker Compose

#### Frontend (React + TypeScript)
- **Framework**: React 19.1.1 con TypeScript
- **Build Tool**: Vite
- **Routing**: React Router DOM
- **HTTP Client**: Axios
- **Gestión de Estado**: React Context API

---

## 🗃️ Modelo de Datos y Relaciones

### Entidades Principales

#### 1. **Users** (Usuario Base)
```sql
- id: Primary Key
- username: Identificador único
- password: Hash bcrypt
- role: "andrei" | "daemon" | "network_admin"
- status: "active" | "captured" | "punished"
- created_at, updated_at, deleted_at
```

#### 2. **Captures** (Sistema de Capturas)
```sql
- id: Primary Key
- daemon_id: FK → Users (daemon que realiza la captura)
- target_id: FK → Users (network_admin capturado)
- difficulty: "easy" | "medium" | "hard"
- points: 100 | 250 | 500 (según dificultad)
- method: Descripción del método de captura
```

#### 3. **Reports** (Sistema de Reportes)
```sql
- id: Primary Key
- user_id: FK → Users (puede ser null para reportes anónimos)
- type: "resistance" | "capture" | "anonymous"
- content: Contenido del reporte
- status: "pending" | "approved" | "rejected"
- is_anonymous: Boolean
```

#### 4. **Punishments** (Sistema de Castigos)
```sql
- id: Primary Key
- user_id: FK → Users
- assigned_by: FK → Users (quien asigna el castigo)
- type: "timeout" | "demotion" | "extra_tasks" | "reward"
- reason: Motivo del castigo
- expires_at: Fecha de expiración
```

#### 5. **Statistics** (Estadísticas de Rendimiento)
```sql
- id: Primary Key
- user_id: FK → Users (solo para daemons)
- total_captures: Número total de capturas
- total_points: Puntos acumulados
- rank_position: Posición en el ranking
- reports_made: Reportes realizados
```

### Relaciones Clave
- **User → Captures**: Un daemon puede tener múltiples capturas
- **User → Reports**: Un usuario puede crear múltiples reportes
- **User → Punishments**: Un usuario puede recibir múltiples castigos
- **User → Statistics**: Relación 1:1 para daemons

---

## 🔐 Sistema de Autenticación JWT

### Flujo de Autenticación

#### 1. **Estructura del Token JWT**
```go
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims // Incluye exp, iat, iss
}
```

#### 2. **Proceso de Login**
```
1. Usuario envía credenciales (username, password)
2. Backend valida contra hash bcrypt almacenado
3. Si es válido, genera JWT con claims del usuario
4. Token se firma con HMAC-SHA256 y secret key
5. Frontend recibe token y lo almacena en localStorage
6. Token se incluye en header Authorization: Bearer <token>
```

#### 3. **Validación del Token**
```go
// En cada request protegido:
1. Extraer token del header Authorization
2. Validar formato "Bearer <token>"
3. Parsear y verificar firma JWT
4. Validar expiración (24 horas)
5. Extraer claims del usuario
6. Inyectar datos del usuario en contexto Gin
```

### Seguridad del JWT

- **Expiración**: 24 horas de validez
- **Algoritmo**: HMAC-SHA256
- **Secret Key**: Variable de entorno JWT_SECRET
- **Storage**: localStorage en frontend (consideración de seguridad)

---

## 🛡️ Sistema de Middleware

### 1. **Middleware de Autenticación** (`auth_middleware.go`)

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Extraer token del header
        authHeader := c.GetHeader("Authorization")

        // 2. Validar formato Bearer
        if !strings.HasPrefix(authHeader, "Bearer ") {
            c.JSON(401, gin.H{"error": "Missing or invalid token"})
            c.Abort()
            return
        }

        // 3. Parsear y validar JWT
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := utils.ValidateJWT(tokenString)

        // 4. Inyectar usuario en contexto
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)

        c.Next()
    }
}
```

### 2. **Middleware de Roles**

#### Control de Acceso Granular:
```go
// Solo para Andrei (líder supremo)
func AndreiOnlyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        role := c.GetString("role")
        if role != "andrei" {
            c.JSON(403, gin.H{"error": "Access denied"})
            c.Abort()
            return
        }
        c.Next()
    }
}

// Solo para Daemons (hackers)
func DaemonOnlyMiddleware() gin.HandlerFunc {
    // Similar implementación para role == "daemon"
}

// Solo para Network Admins (resistencia)
func NetworkAdminOnlyMiddleware() gin.HandlerFunc {
    // Similar implementación para role == "network_admin"
}
```

### 3. **Middleware CORS**
```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(200)
            return
        }

        c.Next()
    }
}
```

---

## 🚀 Flujo de Aplicación por Roles

### 👑 **Andrei Dashboard** (Líder Supremo)

#### Funcionalidades:
- **Vista General del Sistema**: Estadísticas globales, usuarios activos
- **Gestión Completa de Usuarios**: CRUD operations para todos los roles
- **Control de Castigos**: Asignación y gestión de punishments
- **Supervisión de Capturas**: Monitoreo de actividad de daemons
- **Gestión de Reportes**: Aprobación/rechazo de reportes

#### Endpoints Exclusivos:
```
GET  /api/dashboard/andrei     - Estadísticas del sistema
GET  /api/users                - Lista todos los usuarios
POST /api/auth/register        - Registro de nuevos usuarios
PUT  /api/users/:id           - Actualizar cualquier usuario
DELETE /api/users/:id         - Eliminar usuarios
```

### 🤖 **Daemon Dashboard** (Hackers)

#### Funcionalidades:
- **Panel Operativo**: Estadísticas personales, ranking, puntos
- **Sistema de Capturas**: Capturar network admins disponibles
- **Leaderboard**: Ranking entre daemons
- **Gestión de Reportes**: Crear reportes de actividad
- **Estado de Castigos**: Ver castigos activos

#### Lógica de Capturas:
```go
// Sistema de dificultad basado en ID del target
func calculateDifficulty(targetID uint) (string, int) {
    switch targetID % 3 {
    case 0: return "easy", 100      // +100 puntos
    case 1: return "medium", 250    // +250 puntos
    case 2: return "hard", 500      // +500 puntos
    }
}
```

#### Endpoints Principales:
```
GET  /api/dashboard/daemon         - Estadísticas del daemon
POST /api/users/:id/capture        - Capturar network admin
GET  /api/users/:id/captures       - Historial de capturas
GET  /api/leaderboard              - Ranking de daemons
```

### 🛡️ **Network Admin Page** (Resistencia)

#### Funcionalidades:
- **Portal de Resistencia**: Información de supervivencia
- **Reportes Anónimos**: Sistema de denuncias seguras
- **Guía de Supervivencia**: Tips para evitar capturas
- **Contactos de Emergencia**: Red de apoyo de la resistencia
- **Estado de Supervivencia**: Estadísticas personales

#### Características Especiales:
- **Modo Anónimo**: Reportes sin identificación del usuario
- **Interfaz de Supervivencia**: UI enfocada en evasión y resistencia
- **Sistema de Alerts**: Notificaciones de peligro

---

## 🔄 Flujo de Datos en Tiempo Real

### Ciclo de Request-Response

#### 1. **Frontend → Backend**
```typescript
// Ejemplo: Login de usuario
const loginUser = async (username: string, password: string) => {
    const response = await api.post('/auth/login', {
        username,
        password
    });

    // Almacenar token en localStorage
    localStorage.setItem('token', response.data.data.token);

    // Actualizar contexto de autenticación
    setUser(response.data.data.user);
};
```

#### 2. **Backend: Flujo de Middleware**
```
Request → CORS → Auth → Role → Controller → Service → DAO → Database
Response ← JSON ← DTO ← Service ← Controller ← Middleware ← Database
```

#### 3. **Manejo de Estados**
```typescript
// Context API para estado global de autenticación
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Estados locales para componentes específicos
const [loading, setLoading] = useState(true);
const [error, setError] = useState<string | null>(null);
const [data, setData] = useState<DashboardData | null>(null);
```

---

## 🎨 Diseño y UX

### Tema Cyberpunk/Terminal

#### Colores:
- **Primario**: Verde Matrix (#00ff00)
- **Secundario**: Cyan (#00ffff)
- **Peligro**: Rojo Neon (#ff0040)
- **Fondo**: Negro (#000000)
- **Texto**: Verde claro (#00cc00)

#### Tipografía:
- **Fuente Principal**: Fira Code (monospace)
- **Efectos**: Glitch, terminal cursor, scanlines

#### Componentes UI:
- **Terminal Windows**: Simulan interfaces de hacking
- **Progress Bars**: Estilo retro con efectos de carga
- **Buttons**: Diseño de teclas de terminal
- **Cards**: Ventanas de sistema con bordes pixelados

---

## 🚦 Gestión de Estados y Errores

### Estados de Usuario:
- **active**: Usuario normal operativo
- **captured**: Network admin capturado por daemon
- **punished**: Usuario con restricciones temporales

### Manejo de Errores:
```typescript
// Frontend: Error boundaries y states
try {
    const result = await apiCall();
    setData(result.data);
} catch (error) {
    setError(error.response?.data?.message || 'Error desconocido');
} finally {
    setLoading(false);
}
```

```go
// Backend: Respuestas estructuradas
type ApiResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}
```

---

## 🔧 Configuración y Deployment

### Variables de Entorno:
```env
# Backend
DB_HOST=localhost
DB_PORT=5432
DB_USER=chaos_user
DB_PASSWORD=chaos_password
DB_NAME=devops_chaos
JWT_SECRET=your-secret-key
GIN_MODE=release

# Frontend
VITE_API_URL=http://localhost:8080/api
```

### Docker Compose:
```yaml
services:
  backend:
    build: ./devops-chaos-backend
    ports: ["8080:8080"]
    depends_on: [postgres]

  frontend:
    build: ./devops-chaos-frontend
    ports: ["3000:80"]
    depends_on: [backend]

  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: devops_chaos
```

---

## 🎯 Justificación Técnica

### ¿Por qué JWT?
- **Stateless**: No requiere almacenamiento de sesiones en servidor
- **Escalable**: Ideal para arquitecturas distribuidas
- **Seguro**: Firmado criptográficamente, no se puede falsificar
- **Flexible**: Permite incluir metadata del usuario en el token

### ¿Por qué Middleware?
- **Separation of Concerns**: Lógica de autenticación separada de business logic
- **Reutilización**: Mismo middleware para múltiples endpoints
- **Seguridad**: Validación centralizada y consistente
- **Mantenibilidad**: Cambios de autenticación en un solo lugar

### ¿Por qué Context API?
- **Global State**: Estado de autenticación accesible desde cualquier componente
- **Eficiencia**: Evita prop drilling innecesario
- **Simplicidad**: Menos boilerplate que Redux para casos simples
- **Performance**: Re-renders optimizados para cambios específicos

---

## 📊 Métricas y Monitoreo

### Sistema de Puntos:
- **Captura Easy**: 100 puntos
- **Captura Medium**: 250 puntos
- **Captura Hard**: 500 puntos

### Rankings:
- Calculado automáticamente basado en puntos totales
- Actualización en tiempo real tras cada captura
- Leaderboard visible para todos los daemons

### Estadísticas Tracked:
- Total de capturas por daemon
- Puntos acumulados
- Posición en ranking
- Reportes realizados
- Status de castigos activos

---

## 🔮 Futuras Mejoras

### Backend:
- WebSockets para actualizaciones en tiempo real
- Sistema de notificaciones push
- Rate limiting para prevenir spam
- Logs de auditoría completos

### Frontend:
- Modo offline con sincronización
- Notificaciones browser nativas
- Lazy loading de componentes
- Progressive Web App (PWA)

### DevOps:
- Pipeline CI/CD automatizado
- Monitoring con Prometheus/Grafana
- Backup automático de base de datos
- SSL/TLS certificates automáticos

---

**DevOps Chaos** representa una implementación completa de una aplicación moderna con autenticación robusta, control de acceso granular y una experiencia de usuario inmersiva que combina funcionalidad técnica con narrativa cyberpunk.