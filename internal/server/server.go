package server

import (
	"context"
	"net/http"
	"time"

	"github.com/ryanrmg/backend-alpha/internal/api"
	config "github.com/ryanrmg/backend-alpha/internal/config"
	"github.com/ryanrmg/backend-alpha/internal/repository"
	"github.com/ryanrmg/backend-alpha/internal/service"
)

type Server struct {
	http *http.Server
	db   *repository.Database
}

func New(
	ctx context.Context,
	cfg config.Config,
) (*Server, error) {

	db, err := repository.NewDatabase(
		ctx,
		cfg.DBConnString,
	)

	if err != nil {
		return nil, err
	}

	repo := repository.NewPostgresTradeRepository(db.Pool)
	service := service.NewTradeService(repo)
	handler := api.NewTradeHandler(service)
	router := api.NewRouter(handler)

	httpServer := &http.Server{
		Addr:         cfg.Port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		http: httpServer,
		db:   db,
	}, nil
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(
	ctx context.Context,
) error {

	if err := s.http.Shutdown(ctx); err != nil {
		return err
	}

	s.db.Close()

	return nil
}
