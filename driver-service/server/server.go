package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kadirhanmeral/driver-management/configs"
	"github.com/rs/zerolog"
)

type Server struct {
	l      zerolog.Logger
	router *gin.Engine
	config *configs.Config
}

func NewServer(l zerolog.Logger, router *gin.Engine, config *configs.Config) *Server {
	return &Server{l: l, router: router, config: config}
}

func (s *Server) Serve() {
	srv := &http.Server{
		Addr:    s.config.Server.Address,
		Handler: s.router,
	}

	go func() {
		s.l.Info().Msgf("Server listening on %s", s.config.Server.Address)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.l.Fatal().Err(err).Msg("Server listen error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.l.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.l.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	<-ctx.Done()
	s.l.Info().Msg("Server exiting")
}
