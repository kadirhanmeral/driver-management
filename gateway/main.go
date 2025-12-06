package main

import (
	"go-api-gateway/config"
	"go-api-gateway/handlers"
	"go-api-gateway/middleware"
	"go-api-gateway/utils"
	"log"
	"strings"
	"time"

	_ "go-api-gateway/docs"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

// @title           API Gateway
// @version         1.0
// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @title API Gateway
// @version 1.0

func main() {

	router := gin.Default()

	router.Use(middleware.RequestLogger())

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	utils.InitElasticsearch(cfg)

	cache := cache.New(5*time.Minute, cfg.RateLimitWindow)

	router.POST("/auth/token", handlers.GetAuthToken(cfg))

	registerServices(router, cfg, cache)

	log.Printf("API Gateway running on :%s", cfg.ServicePort)
	router.Run(":" + cfg.ServicePort)
}

func registerServices(router *gin.Engine, cfg *config.Config, cache *cache.Cache) {

	for serviceName, svc := range cfg.Services {
		for _, route := range svc.Routes {
			gatewayPath := route.Path

			ginPath := gatewayPath
			if strings.Contains(ginPath, "{") {
				ginPath = strings.ReplaceAll(ginPath, "{", ":")
				ginPath = strings.ReplaceAll(ginPath, "}", "")
			}

			log.Printf("Registering route: %s (Gin: %s) -> %s%s", gatewayPath, ginPath, svc.BaseURL, route.Path)

			router.Any(ginPath, middleware.JWTAuthMiddleware(cfg.JWTSecretKey), handlers.MakeProxyHandler(serviceName, svc.BaseURL, cfg, cache))
		}
	}

	router.GET("/swagger/*any", handlers.SwaggerHandler(cfg))

}
