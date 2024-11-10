package server

import (
	"fmt"
	"net/http"
	"scoreplay/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type TagsHandler struct {
	logger  *zerolog.Logger
	service service.TagService
}

func NewTagsHandler(logger *zerolog.Logger, service service.TagService) *TagsHandler {
	return &TagsHandler{
		logger:  logger,
		service: service,
	}
}

func (h *TagsHandler) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Get("/", h.handleListTags)
		r.Post("/", h.handleCreateTag)
	})

	return r
}

func (h *TagsHandler) handleCreateTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With().Str("handler", "handleCreateTag").Logger()
	tag, err := h.service.CreateTag(ctx, "test")
	if err != nil {
		logger.Error().Err(err).Msg("Could not create tag")
		//return error json response

		return
	}

	fmt.Println("Test create tag", tag)
}

func (h *TagsHandler) handleListTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With().Str("handler", "handleListTags").Logger()

	tags, err := h.service.ListTags(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Could not list tags")

		//return error json response
		return
	}

	fmt.Println("Test list tags", tags)
}
