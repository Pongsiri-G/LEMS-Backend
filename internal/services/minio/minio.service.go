package minio

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	minioRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/minio"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
)

type Service interface {
	Upload(ctx context.Context, file multipart.FileHeader) (string, error)
	GetImage(ctx context.Context, url string) ([]byte, string, error)
}
type minioService struct {
	minioRepo minioRepo.Repository
}

func NewMinioService(minioRepo minioRepo.Repository) Service {
	return &minioService{
		minioRepo: minioRepo,
	}
}

// GetImage implements MinioService.
func (m *minioService) GetImage(ctx context.Context, url string) ([]byte, string, error) {
	bucketName, objectName, err := utils.ExtractUrl(url)
	if err != nil {
		return nil, "", err
	}
	return m.minioRepo.GetImage(ctx, bucketName, objectName)
}

// Upload implements MinioService.
func (m *minioService) Upload(ctx context.Context, file multipart.FileHeader) (string, error) {
	buffer, err := file.Open()
	if err != nil {
		return "", err
	}

	defer buffer.Close()

	objectName := uuid.New().String() + "-" + file.Filename
	fileBuffer := buffer

	var contentType string
	if ctHeaders, ok := file.Header["Content-Type"]; ok && len(ctHeaders) > 0 {
		contentType = ctHeaders[0]
	} else {
		if len(file.Filename) > 0 {
			switch {
			case strings.HasSuffix(strings.ToLower(file.Filename), ".jpg"), strings.HasSuffix(strings.ToLower(file.Filename), ".jpeg"):
				contentType = "image/jpeg"
			case strings.HasSuffix(strings.ToLower(file.Filename), ".png"):
				contentType = "image/png"
			default:
				return "", exceptions.ErrInvalidFileFormat
			}
		} else {
			return "", exceptions.ErrInvalidFileFormat
		}
	}

	fileSize := file.Size

	// Allow only image and pdf
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "application/pdf" {
		return "", exceptions.ErrInvalidFileFormat
	}

	return m.minioRepo.Upload(ctx, objectName, fileBuffer, fileSize, contentType)
}
