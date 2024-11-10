package server

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

type HTTPError struct {
	Message string `json:"message"`
}

// sendJSON is the exit point of each handled request.
// It ensures that the endpoints return a valid json response with the correct HTTP status code.
// In case of error, the error message is attached to the json response.
func sendJSON(w http.ResponseWriter, logger *zerolog.Logger, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status > http.StatusCreated {
		logger.Error().Msgf("Request failed: %v", data)
		if err := json.NewEncoder(w).Encode(&HTTPError{Message: data.(string)}); err != nil {
			logger.Error().Err(err).Msg("Could not send json response")
		}

		return
	}

	if err := json.NewEncoder(w).Encode(&data); err != nil {
		logger.Error().Err(err).Msg("Could not send json response")
	}
}
