package server

import (
	"net/http"

	"scoreplay/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer(cfg *config.Config, address string, f func(chi.Router)) *http.Server {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.Server.RequestTimeout))

	f(r)

	return &http.Server{
		Addr:         address,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		Handler:      r,
	}
}
