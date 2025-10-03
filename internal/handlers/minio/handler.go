package minio

import (
	"net/http"

	minioSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/minio"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type FileHandler interface {
	Upload(c echo.Context) error
	GetImage(c echo.Context) error
}

type fileHandler struct {
	service minioSvc.Service
}

func NewFileHandler(service minioSvc.Service) FileHandler {
	return &fileHandler{
		service: service,
	}
}

// Upload implements FileHandler.
func (f *fileHandler) Upload(c echo.Context) error {
	log.Debug().Msgf("Starting file upload process , %v", c.Request().Header)
	file, err := c.FormFile("file")
	if err != nil {
		log.Error().Msgf("File upload bind error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "failed to get file from form data",
		})
	}

	if file == nil || file.Size == 0 {
		log.Warn().Msg("No file provided or file is empty")
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "file is required",
		})
	}

	// Debug logging for file details
	log.Debug().Msgf("Received file - Name: %s, Size: %d bytes, Headers: %+v",
		file.Filename, file.Size, file.Header)

	url, err := f.service.Upload(c.Request().Context(), *file)
	if err != nil {
		log.Error().Msgf("File upload error: %v", err)
		// Check for specific error types to provide better error messages
		if err.Error() == "invalid file format" {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "invalid file format. Only JPEG, PNG, and PDF files are allowed",
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to upload file",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"url": url,
	})
}

// GetImage implements FileHandler.
func (f *fileHandler) GetImage(c echo.Context) error {
	var req struct {
		URL string `json:"url"`
	}
	if err := c.Bind(&req); err != nil {
		log.Error().Msgf("Get image bind error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request body",
		})
	}

	imageData, contentType, err := f.service.GetImage(c.Request().Context(), req.URL)
	if err != nil {
		log.Error().Msgf("Get image error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to get image",
		})
	}

	return c.Blob(http.StatusOK, contentType, imageData)
}
