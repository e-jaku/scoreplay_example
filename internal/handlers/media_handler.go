package handlers

import (
	"fmt"
	"net/http"
	"scoreplay/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type MediaHandler struct {
	logger  *zerolog.Logger
	service service.MediaService
}

func NewMediaHandler(logger *zerolog.Logger, service service.MediaService) *MediaHandler {
	return &MediaHandler{}
}

func (h *MediaHandler) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Get("/", h.handleListMediaByTagId)
		r.Post("/", h.handleCreateMedia)
	})

	return r
}

func (h *MediaHandler) handleCreateMedia(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With().Str("handler", "handleCreateMedia").Logger()

	media, err := h.service.CreateMedia(ctx, "test", []string{})
	if err != nil {
		logger.Error().Err(err).Msg("could not create media")
	}
	fmt.Println("Test create media", media)
}

func (h *MediaHandler) handleListMediaByTagId(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := h.logger.With().Str("handler", "handleListMediaByTagId").Logger()

	medias, err := h.service.ListMediaByTagId(ctx, "test")
	if err != nil {
		logger.Error().Err(err).Msg("could not list media by tag id")
	}

	fmt.Println("Test create media", medias)
}
