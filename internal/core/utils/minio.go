package utils

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

var (
	minioClient = &minio.Client{}
	ctx         = context.Background()
)

// Client minio client interface
type Client interface {
	UploadImage(bucketName string, objectName, filePath string) error
}

type MinioClient struct {
	client *minio.Client
}

// Configuration config minio for new connection
type MinioConfiguration struct {
	Host            string
	AccessKeyID     string
	SecretAccessKey string
}

// NewConnection new ftp connection
func NewConnection(config MinioConfiguration) (client *MinioClient, err error) {
	// Initialize minio client object.
	minioClient, err = minio.New(config.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		logrus.Errorf("[NewConnection] error creating Minio client: %v", err)
		return nil, err
	}

	return &MinioClient{client: minioClient}, nil
}

// GetMinioClient get minio client
func GetMinioClient() (*MinioClient, error) {
	minioConfig := MinioConfiguration{
		Host:            "localhost:9000",
		AccessKeyID:     "ROOTNAME",
		SecretAccessKey: "CHANGEME123",
		// Host:            os.Getenv("MINIO_HOST"),
		// AccessKeyID:     os.Getenv("MINIO_ACCESS_KEY"),
		// SecretAccessKey: os.Getenv("MINIO_SECRET_KEY"),
	}

	minioClient, err := NewConnection(minioConfig)
	if err != nil {
		log.Fatal("minio client error", err)
		return nil, err
	}

	return minioClient, nil
}

// UploadImage uploads an image to Minio using a pre-configured Minio client
func (m *MinioClient) UploadImage(bucketName string, objectName string, filePath string) error {
	err := m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		exists, errBucketExists := m.client.BucketExists(ctx, bucketName)
		if errBucketExists != nil {
			logrus.Errorf("[UploadImage] check bucket exists error: %s", err)
			return err
		}

		if !exists {
			logrus.Errorf("[UploadImage] make bucket error: %s", err)
			return err
		}
	}

	contentType := "image/pdf"
	info, err := m.client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		logrus.Errorf("[UploadImage] put object error: %s", err)
		return err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return nil
}
