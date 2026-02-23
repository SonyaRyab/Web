package repository

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient структура для работы с Minio
type MinioClient struct {
	Client     *minio.Client
	BucketName string
	Endpoint   string
	UseSSL     bool
}

func NewMinioClient(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка создания Minio клиента: %v", err)
	}

	return &MinioClient{
		Client:     client,
		BucketName: bucket,
		Endpoint:   endpoint,
		UseSSL:     useSSL,
	}, nil
}

// GetFileURL возвращает публичный URL файла в Minio
func (m *MinioClient) GetFileURL(objectName string) string {
	if objectName == "" {
		return ""
	}

	protocol := "http"
	if m.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, m.Endpoint, m.BucketName, objectName)
}
