package server

import (
	"encoding/json"
	"net/http"
	"scoreplay/internal/domain"
	"scoreplay/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"golang.org/x/xerrors"
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
	logger.Info().Msg("Creating tags...")

	tagToCreate := &domain.Tag{}
	err := json.NewDecoder(r.Body).Decode(&tagToCreate)
	if err != nil {
		sendJSON(w, &logger, http.StatusBadRequest, xerrors.Errorf("failed to parse request body: %w", err).Error())
		return
	}

	if tagToCreate.Name == "" {
		sendJSON(w, &logger, http.StatusBadRequest, "tag name is required")
		return
	}

	logger.With().Str("tag", tagToCreate.Name)

	tag, err := h.service.CreateTag(ctx, tagToCreate.Name)
	if err != nil {
		sendJSON(w, &logger, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info().Msgf("Created tag: %v", tag.ID)
	sendJSON(w, &logger, http.StatusCreated, tag)
}

func (h *TagsHandler) handleListTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With().Str("handler", "handleListTags").Logger()
	logger.Info().Msg("Listing tags...")

	tags, err := h.service.ListTags(ctx)
	if err != nil {
		sendJSON(w, &logger, http.StatusBadRequest, err.Error())
		return
	}

	sendJSON(w, &logger, http.StatusOK, tags)
}
