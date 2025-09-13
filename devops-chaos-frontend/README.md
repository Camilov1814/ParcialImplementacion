# DevOps Chaos Frontend

Frontend para la aplicación DevOps Chaos Backend - Una aplicación temática de juego de roles cyberpunk.

## 🚀 Inicio Rápido

### Prerequisitos
- Node.js 20.19+ o 22.12+
- npm o yarn
- Backend de DevOps Chaos ejecutándose en `http://localhost:8080`

### Instalación

```bash
# Instalar dependencias
npm install

# Ejecutar en modo desarrollo
npm run dev

# Construir para producción
npm run build

# Preview del build de producción
npm run preview
```

La aplicación estará disponible en `http://localhost:5173`

## 🎮 Credenciales de Prueba

### Andrei (Supreme Leader)
- **Username:** `andrei`
- **Password:** `AndreI2024!`
- **Rol:** Control total del sistema

### Daemon (Hacker)
- **Username:** `daemon_alpha`
- **Password:** `DaemonAlpha123!`
- **Rol:** Operativo de caos

### Network Admin (Resistance)
- **Username:** `netadmin_alice`
- **Password:** `NetworkAlice123!`
- **Rol:** Resistencia

## 🎨 Características del Frontend

### Tema Visual
- **Estilo:** Cyberpunk/Hacker oscuro
- **Colores:** Verde matrix (#00ff00), fondo negro, acentos rojos
- **Efectos:** Glitch effects, terminal windows, animaciones
- **Fuente:** Fira Code (monospace)

### Funcionalidades por Rol

#### 🎛️ Andrei Dashboard
- Estadísticas del sistema completo
- Top 5 daemons más efectivos
- Reportes de intel recientes
- Panel de control de operaciones
- Feed de actividad en tiempo real

#### 💻 Daemon Dashboard
- Estadísticas personales (capturas, puntos, ranking)
- Leaderboard completo de daemons
- Misiones activas disponibles
- Castigos activos (si los hay)
- Panel de operaciones daemon

#### 🛡️ Resistance Page
- Portal de bienvenida para la resistencia
- Estadísticas de supervivencia
- Guía de supervivencia con consejos
- Contactos de emergencia
- Formulario de reportes anónimos
- Memes de resistencia

### 🔐 Sistema de Autenticación
- JWT Bearer token authentication
- Context API para estado global
- Rutas protegidas por rol
- Renovación automática de sesión

## 🗂️ Estructura del Proyecto

```
src/
├── components/
│   ├── auth/           # Login y protección de rutas
│   ├── dashboard/      # Dashboards específicos por rol
│   └── common/         # Layout y navegación
├── contexts/           # Context API para autenticación
├── services/          # Llamadas a la API
├── types/             # Tipos TypeScript
└── styles/            # Tema cyberpunk global
```

## 🔗 Rutas Disponibles

- `/` - Redirección automática al dashboard
- `/login` - Página de autenticación
- `/dashboard` - Dashboard específico por rol
- `/users` - Gestión de usuarios (solo Andrei)
- `/reports` - Gestión de reportes
- `/punishments` - Sistema de castigos
- `/leaderboard` - Ranking de daemons
- `/capture` - Misión de captura (solo Daemons)
- `/contacts` - Contactos de emergencia (solo Network Admins)
- `/guide` - Guía de supervivencia (solo Network Admins)

## 🛠️ Tecnologías Utilizadas

- **React 18** con hooks
- **TypeScript** para tipado seguro
- **Vite** como bundler
- **React Router** para navegación
- **Axios** para llamadas HTTP
- **CSS3** con variables custom para theming

## 🔧 Configuración del Backend

El frontend espera que el backend esté ejecutándose en:
- **URL:** `http://localhost:8080`
- **API Base:** `http://localhost:8080/api`

### Estructura de Respuesta Esperada del Backend

```json
{
  "success": true,
  "message": "Operation completed",
  "data": {
    // Datos específicos del endpoint
  }
}
```

## 🎯 Funcionalidades Principales

### Autenticación
- Login con credenciales
- Manejo automático de tokens JWT
- Renovación de sesión
- Logout seguro

### Navegación por Roles
- Sidebar dinámico según rol
- Rutas protegidas
- Redirección automática

### Dashboards Interactivos
- Estadísticas en tiempo real
- Gráficos y métricas
- Actividad del sistema

### Tema Cyberpunk
- Efectos visuales de terminal
- Animaciones de glitch
- Paleta de colores hacker

## 🐛 Troubleshooting

### El login no funciona
1. Verificar que el backend esté ejecutándose en `http://localhost:8080`
2. Revisar las credenciales de prueba
3. Verificar la consola del navegador para errores de CORS

### Errores de compilación TypeScript
1. Ejecutar `npm install` para reinstalar dependencias
2. Verificar que todos los imports usen `type` cuando sea necesario

### Problemas de estilo
1. Verificar que `src/styles/globals.css` esté importado
2. Limpiar caché del navegador
3. Ejecutar `npm run build` para verificar errores

## 📱 Responsive Design

La aplicación está optimizada para:
- Desktop (1200px+)
- Tablet (768px - 1199px)
- Mobile (320px - 767px)

---

**Temática:** Una realidad distópica donde Andrei controla un ejército de hackers (daemons) que buscan capturar administradores de red, mientras estos últimos forman una resistencia clandestina.

¡Bienvenido al caos digital! 🚀💀
