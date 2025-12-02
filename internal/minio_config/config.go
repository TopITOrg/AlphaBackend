package minio_config

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type MinioConfig struct {
	Endpoint    string `env:"MINIO_ENDPOINT"`
	AccessKeyID string `env:"MINIO_ACCESS_KEY"`
	SecretKey   string `env:"MINIO_SECRET_KEY"`
	UseSSL      bool   `env:"MINIO_USE_SSL"`
	BucketName  string `env:"MINIO_BUCKET_NAME"`
}

func InitBuckets(ctx context.Context, client *minio.Client, bucketName string) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return err
		}
	}
	return nil
}
