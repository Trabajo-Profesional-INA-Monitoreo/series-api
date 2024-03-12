package middlewares

import (
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func IsAValidToken(apiConfig *config.ApiConfig) gin.HandlerFunc {
	authService := services.NewKeycloakAuthService(apiConfig.KeycloakConfig)
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("missing token in request")))
			return
		}
		isValid := authService.IsAValidToken(c, bearerToken[1])
		if !isValid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

func IsAnAdminToken(apiConfig config.ApiConfig) gin.HandlerFunc {
	authService := services.NewKeycloakAuthService(apiConfig.KeycloakConfig)
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("missing token in request")))
			return
		}
		isValid := authService.IsAnAdminToken(c, bearerToken[1])
		if !isValid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
