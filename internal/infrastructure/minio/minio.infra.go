package minio

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

func NewMinioConnection(cfg *configs.Config) (*minio.Client, error) {
	ctx := context.Background()
	minioClient, err := minio.New(cfg.MINIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MINIO.User, cfg.MINIO.Password, ""),
		Secure: cfg.MINIO.UseSSL,
		Region: cfg.MINIO.Region, // เพิ่ม Region สำหรับ S3
	})
	if err != nil {
		log.Error().Msgf("Error initializing MinIO client: %v", err)
		return nil, err
	}
	_, err = minioClient.ListBuckets(ctx)
	if err != nil {
		log.Error().Msgf("Error connecting to MinIO server: %v", err)
		return nil, err
	}
	return minioClient, nil
}
