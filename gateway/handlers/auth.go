package handlers

import (
	"go-api-gateway/config"
	"go-api-gateway/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthDTO struct {
	ApiKey string `json:"apiKey" binding:"required"`
}

// @Summary      Get JWT Token
// @Description  Get a JWT token by providing a valid API key
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body	AuthDTO  true  "Get Auth Token"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /auth/token [post]
func GetAuthToken(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		token, err := middleware.GenerateToken(req.ApiKey, cfg.APIKey, cfg.JWTSecretKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
