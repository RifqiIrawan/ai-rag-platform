package storage

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/RifqiIrawan/ai-rag-platform/services/document-service/internal/config"
)

// NewMinioClient constructs a client without making a network call,
// so the service can start even if MinIO is briefly unavailable.
func NewMinioClient(cfg *config.Config) (*minio.Client, error) {
	return minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioRootUser, cfg.MinioRootPass, ""),
		Secure: cfg.MinioUseSSL,
	})
}

// EnsureBucket creates the configured bucket if it doesn't already exist.
func EnsureBucket(ctx context.Context, client *minio.Client, bucket string) error {
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
}
