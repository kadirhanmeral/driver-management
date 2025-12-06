package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kadirhanmeral/driver-management/configs"
	"github.com/kadirhanmeral/driver-management/internal/handlers"
	"github.com/kadirhanmeral/driver-management/internal/repository"
	"github.com/kadirhanmeral/driver-management/server"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	services "github.com/kadirhanmeral/driver-management/internal/services"
	routes "github.com/kadirhanmeral/driver-management/server/router"
)

// @title           Driver Service API
// @version         1.0
// @description     API Documentation for Driver Service

// @host      localhost:8080
// @BasePath  /

func main() {

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	cfg := configs.NewConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.URI))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to Mongo")
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			logger.Error().Err(err).Msg("Mongo disconnect error")
		}
	}()

	driverRepo := repository.NewDriverRepository(client, cfg.Database.Database, "drivers")

	driverService := services.NewDriverService(driverRepo)

	driverHandler := handlers.NewDriverHandler(driverService)

	router := gin.Default()

	routes.RegisterDriverEndpoints(router, driverHandler)

	srv := server.NewServer(logger, router, cfg)
	srv.Serve()
}
