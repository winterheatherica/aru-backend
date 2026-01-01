package config

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MinioConfig struct {
	Endpoint      string
	AccessKey     string
	SecretKey     string
	UseSSL        bool
	Bucket        string
	PublicBaseURL string
}

func NewMinioConfig(v *viper.Viper) MinioConfig {
	return MinioConfig{
		Endpoint:      v.GetString("minio.endpoint"),
		AccessKey:     v.GetString("minio.accessKey"),
		SecretKey:     v.GetString("minio.secretKey"),
		UseSSL:        v.GetBool("minio.useSSL"),
		Bucket:        v.GetString("minio.bucket"),
		PublicBaseURL: v.GetString("minio.publicBaseURL"),
	}
}

func NewMinioClient(cfg MinioConfig, log *logrus.Logger) *minio.Client {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		log.Fatalf("failed to init MinIO client: %v", err)
	}

	return client
}
