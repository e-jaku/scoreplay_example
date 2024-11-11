package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"scoreplay/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
)

type mockedMediaService struct {
	err error
}

func (s *mockedMediaService) CreateMedia(ctx context.Context, name string, tags []string, file io.Reader, fileType string) (*domain.Media, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &domain.Media{
		ID:      uuid.New(),
		Name:    name,
		FileURL: fileType, // for simplicity
	}, nil
}

func (s *mockedMediaService) ListMediaByTagId(ctx context.Context, tagId string) ([]*domain.Media, error) {
	if s.err != nil {
		return nil, s.err
	}
	return []*domain.Media{}, nil
}

func TestHandleCreateMedia(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Ctx(ctx)
	simulatedServiceError := xerrors.New("Test error")

	tests := []struct {
		name               string
		imagePath          string
		mediaName          string
		fileType           string
		tags               string
		expectedStatusCode int
		mockService        *mockedMediaService
	}{
		{
			name:               "Valid test-case",
			imagePath:          "../../samples/soccer.png",
			mediaName:          "Super nice picture",
			tags:               "tag1",
			fileType:           "png",
			expectedStatusCode: http.StatusCreated,
			mockService:        &mockedMediaService{},
		},
		{
			name:               "Missing tags",
			imagePath:          "../../samples/soccer.png",
			mediaName:          "Super nice picture",
			fileType:           "png",
			expectedStatusCode: http.StatusBadRequest,
			mockService:        &mockedMediaService{},
		},
		{
			name:               "Missing media name",
			imagePath:          "../../samples/soccer.png",
			tags:               "tag1",
			fileType:           "png",
			expectedStatusCode: http.StatusBadRequest,
			mockService:        &mockedMediaService{},
		},
		{
			name:               "Wrong media type(not jpeg or png)",
			mediaName:          "Super nice picture",
			tags:               "tag1",
			fileType:           "png",
			expectedStatusCode: http.StatusUnsupportedMediaType,
			mockService:        &mockedMediaService{err: simulatedServiceError},
		},
		{
			name:               "Error in service",
			mediaName:          "Super nice picture",
			imagePath:          "../../samples/soccer.png",
			tags:               "tag1",
			fileType:           "png",
			expectedStatusCode: http.StatusInternalServerError,
			mockService:        &mockedMediaService{err: simulatedServiceError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &MediaHandler{
				service: tt.mockService,
				logger:  logger,
			}

			var imageData []byte
			var err error
			if tt.imagePath != "" {
				imageData, err = os.ReadFile(tt.imagePath)
				require.NoError(t, err, "Failed to read image file")
			}

			var reqBody bytes.Buffer
			writer := multipart.NewWriter(&reqBody)

			writer.WriteField("name", tt.mediaName)
			writer.WriteField("tags", tt.tags)

			filePart, _ := writer.CreateFormFile("file", "testfile.png")
			if tt.imagePath != "" {
				_, err = filePart.Write(imageData)
			} else {
				_, err = filePart.Write([]byte("fake file content"))
			}
			require.NoError(t, err, "Failed to write image data to form")

			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/media", &reqBody)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			rec := httptest.NewRecorder()
			handler.handleCreateMedia(rec, req)

			require.Equal(t, tt.expectedStatusCode, rec.Code)
		})
	}
}

func TestHandleListMediaByTagId(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Ctx(ctx)
	simulatedServiceError := xerrors.New("Test error")

	tests := []struct {
		name               string
		tagId              string
		expectedStatusCode int
		mockService        *mockedMediaService
	}{
		{
			name:               "Valid test-case",
			tagId:              uuid.NewString(),
			expectedStatusCode: http.StatusOK,
			mockService:        &mockedMediaService{},
		},
		{
			name:               "Missing tag id",
			expectedStatusCode: http.StatusBadRequest,
			mockService:        &mockedMediaService{},
		},
		{
			name:               "Service throws an error",
			tagId:              uuid.NewString(),
			expectedStatusCode: http.StatusInternalServerError,
			mockService:        &mockedMediaService{err: simulatedServiceError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &MediaHandler{
				service: tt.mockService,
				logger:  logger,
			}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/media?tag=%s", tt.tagId), nil)
			rec := httptest.NewRecorder()

			handler.handleListMediaByTagId(rec, req)

			require.Equal(t, tt.expectedStatusCode, rec.Code)
		})
	}

}
