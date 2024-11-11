package server

import (
	"io"
	"net/http"
	"scoreplay/internal/service"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"golang.org/x/xerrors"
)

const MAX_SIZE = 1024 * 1024 * 10 // 10MB

type MediaHandler struct {
	logger  *zerolog.Logger
	service service.MediaService
}

func NewMediaHandler(logger *zerolog.Logger, service service.MediaService) *MediaHandler {
	return &MediaHandler{
		logger:  logger,
		service: service,
	}
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
	logger.Info().Msg("Creating media...")

	r.Body = http.MaxBytesReader(w, r.Body, MAX_SIZE)
	if err := r.ParseMultipartForm(MAX_SIZE); err != nil {
		sendJSON(w, &logger, http.StatusBadRequest, xerrors.Errorf("failed to parse request body: %w", err).Error())
		return
	}

	name := r.FormValue("name")
	if name == "" {
		sendJSON(w, &logger, http.StatusBadRequest, "missing media name")
		return
	}

	tagsField := r.FormValue("tags")
	if tagsField == "" {
		sendJSON(w, &logger, http.StatusBadRequest, "missing tags for media")
		return
	}

	var tags []string
	if tagsField != "" {
		tags = strings.Split(tagsField, ",")
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		sendJSON(w, &logger, http.StatusBadRequest, xerrors.Errorf("failed to read file: %w", err).Error())
		return
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		sendJSON(w, &logger, http.StatusBadRequest, xerrors.Errorf("failed to read from file: %w", err).Error())
		return
	}

	fileType := ""
	mimeType := http.DetectContentType(buffer)
	if mimeType == "image/jpeg" {
		fileType = ".jpg"
	} else if mimeType == "image/png" {
		fileType = ".png"
	} else {
		sendJSON(w, &logger, http.StatusUnsupportedMediaType, "invalid file type. Only JPEG and PNG are allowed.")
		return
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		sendJSON(w, &logger, http.StatusInternalServerError, xerrors.Errorf("failed to reset file pointer: %w", err).Error())
		return
	}

	media, err := h.service.CreateMedia(ctx, name, tags, file, fileType)
	if err != nil {
		sendJSON(w, &logger, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info().Msgf("Created media entry %v", media.ID)
	sendJSON(w, &logger, http.StatusCreated, media)
}

func (h *MediaHandler) handleListMediaByTagId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With().Str("handler", "handleListMediaByTagId").Logger()
	logger.Info().Msg("Listing media...")

	tagID := r.URL.Query().Get("tag")
	if tagID == "" {
		sendJSON(w, &logger, http.StatusBadRequest, "missing tag query parameter")
		return
	}

	medias, err := h.service.ListMediaByTagId(ctx, tagID)
	if err != nil {
		sendJSON(w, &logger, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info().Msgf("Found %d media entries matching tag ID: %v", len(medias), tagID)
	sendJSON(w, &logger, http.StatusOK, medias)
}
