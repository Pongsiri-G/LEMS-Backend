package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/minio/minio-go/v7"
)

type Repository interface {
	Upload(ctx context.Context, objectName string, file io.Reader, fileSize int64, contentType string) (string, error)
	GetImage(ctx context.Context, bucketName string, objectName string) ([]byte, string, error)
}

type minioAdaptor struct {
	cfg    *configs.Config
	client *minio.Client
}

func NewMinioRepository(cfg *configs.Config, client *minio.Client) Repository {
	return &minioAdaptor{
		cfg:    cfg,
		client: client,
	}
}

func (m *minioAdaptor) Upload(ctx context.Context, objectName string, file io.Reader, fileSize int64, contentType string) (string, error) {
	_, err := m.client.PutObject(ctx, m.cfg.MINIO.Bucket, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s", m.cfg.MINIO.Bucket, objectName)

	return url, nil
}

func (m *minioAdaptor) GetImage(ctx context.Context, bucketName string, objectName string) ([]byte, string, error) {
	obj, err := m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", err
	}
	defer obj.Close()

	// Get object info to determine content type
	objInfo, err := obj.Stat()
	if err != nil {
		return nil, "", err
	}

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, "", err
	}

	return data, objInfo.ContentType, nil
}
