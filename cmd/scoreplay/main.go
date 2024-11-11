package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"scoreplay/internal/config"
	"scoreplay/internal/db"
	"scoreplay/internal/repository"
	"scoreplay/internal/server"
	"scoreplay/internal/service"
	"scoreplay/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	ctx = logger.WithContext(ctx)

	if err := run(ctx, &logger); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run app")
	}
}

func run(ctx context.Context, logger *zerolog.Logger) error {
	serverCfg, dbCfg, storageCfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load config")
	}

	dbConn, err := db.NewPostgresConnection(dbCfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize database connection")
	}
	defer dbConn.Close()

	done := make(chan struct{})
	quit := make(chan os.Signal, 1)

	mediaRepo := repository.NewPostgresMediaRepository(dbConn)
	tagRepo := repository.NewPostgresTagRepository(dbConn)

	storage, err := storage.NewMinioStorage(storageCfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create MinIO storage")
	}

	mediaService := service.NewMediaService(mediaRepo, tagRepo, storage)
	tagService := service.NewTagService(tagRepo)

	svr := server.NewServer(serverCfg, func(m chi.Router) {
		m.Mount("/media", server.NewMediaHandler(logger, mediaService).Router())
		m.Mount("/tags", server.NewTagsHandler(logger, tagService).Router())
	})

	signal.Notify(quit, os.Interrupt)

	// ensure that the server shuts down gracefully
	go func() {
		<-quit

		logger.Info().Msg("Server is shutting down...")

		c, cancel := context.WithTimeout(ctx, serverCfg.ShutdownTimeout)
		defer cancel()

		svr.SetKeepAlivesEnabled(false)

		if err := svr.Shutdown(c); err != nil {
			logger.Fatal().Err(err).Msg("Failed to gracefully shutdown server")
		}

		close(done)
	}()

	logger.Info().Msgf("Starting server at %s", svr.Addr)

	if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server at %s: %w", svr.Addr, err)
	}

	<-done

	logger.Info().Msg("Server stopped")

	return nil
}
