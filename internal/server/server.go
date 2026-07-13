package server

import (
	"context"
	"net/http"
	"time"

	"github.com/ryanrmg/backend-alpha/internal/api"
	config "github.com/ryanrmg/backend-alpha/internal/config"
	"github.com/ryanrmg/backend-alpha/internal/repository"
	"github.com/ryanrmg/backend-alpha/internal/service"
	projectx "github.com/ryanrmg/projectx-api"
)

type Server struct {
	http *http.Server
	db   *repository.Database
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
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

	client := projectx.NewProjectXClient(
		cfg.ProjectXHttps,
		cfg.ProjectXSocket,
		cfg.ProjectXUsername,
		cfg.ProjectXApiKey,
	)

	repo := repository.NewPostgresTradeRepository(db.Pool)
	tradeService := service.NewTradeService(repo)
	accountService := service.NewAccountService(client)

	tradeHandler := api.NewTradeHandler(tradeService)
	accountHandler := api.NewAccountHandler(accountService)
	router := api.NewRouter(tradeHandler, accountHandler)

	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      cors(router),
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
