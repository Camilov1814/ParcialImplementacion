package middleware

import (
	"devops-chaos-backend/internal/dto"
	"devops-chaos-backend/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware valida el token JWT y agrega información del usuario al contexto
func AuthMiddleware() gin.HandlerFunc {
	authService := services.NewAuthService()

	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "Authorization header required",
				Error:   "Missing Authorization header",
			})
			c.Abort()
			return
		}

		// Verificar formato del token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "Invalid authorization format",
				Error:   "Authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}

		// Validar token y obtener usuario
		user, err := authService.ValidateTokenAndGetUser(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "Invalid or expired token",
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		// Agregar información del usuario al contexto
		c.Set("userID", user.ID)
		c.Set("username", user.Username)
		c.Set("userRole", user.Role)
		c.Set("userStatus", user.Status)
		c.Set("user", user)

		c.Next()
	}
}

// RoleMiddleware verifica que el usuario tenga uno de los roles permitidos
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "User role not found",
				Error:   "User role not found in context",
			})
			c.Abort()
			return
		}

		// Verificar si el rol del usuario está en los roles permitidos
		roleStr := userRole.(string)
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Success: false,
			Message: "Insufficient permissions",
			Error:   "User role not authorized for this operation",
		})
		c.Abort()
	}
}

// AndreiOnlyMiddleware - Solo permite acceso a Andrei
func AndreiOnlyMiddleware() gin.HandlerFunc {
	return RoleMiddleware("andrei")
}

// DaemonOnlyMiddleware - Solo permite acceso a Daemons
func DaemonOnlyMiddleware() gin.HandlerFunc {
	return RoleMiddleware("daemon")
}

// NetworkAdminOnlyMiddleware - Solo permite acceso a Network Admins
func NetworkAdminOnlyMiddleware() gin.HandlerFunc {
	return RoleMiddleware("network_admin")
}

// AndreiAndDaemonMiddleware - Permite acceso a Andrei y Daemons
func AndreiAndDaemonMiddleware() gin.HandlerFunc {
	return RoleMiddleware("andrei", "daemon")
}
