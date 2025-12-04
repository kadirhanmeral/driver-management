package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kadirhanmeral/driver-management/configs"
	"github.com/kadirhanmeral/driver-management/internal/handlers"
	"github.com/kadirhanmeral/driver-management/internal/repository"
	services "github.com/kadirhanmeral/driver-management/internal/services"
	"github.com/kadirhanmeral/driver-management/server"
	routes "github.com/kadirhanmeral/driver-management/server/router"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Logger initialization
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Load configuration
	cfg := configs.NewConfig()

	// MongoDB client settings
	clientOpts := options.Client().ApplyURI(cfg.Database.URI)
	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create Mongo client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to Mongo")
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			logger.Error().Err(err).Msg("Mongo disconnect error")
		}
	}()

	// Repository creation
	driverRepo := repository.NewDriverRepository(client, cfg.Database.Database, "drivers")

	// Service creation
	driverService := services.NewDriverService(driverRepo)

	// Handler creation
	driverHandler := handlers.NewDriverHandler(driverService)

	// Gin router
	router := gin.Default()

	// Routes
	routes.RegisterDriverEndpoints(router, driverHandler)

	// Server initialization and start
	srv := server.NewServer(logger, router, cfg)
	srv.Serve()
}
