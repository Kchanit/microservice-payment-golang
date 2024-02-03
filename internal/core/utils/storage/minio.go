package storage

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioSetup struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

type MinioStorage struct {
	Client *minio.Client
}

func NewMinio(config MinioSetup) (*MinioStorage, error) {
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return &MinioStorage{
		Client: minioClient,
	}, nil

}

func (m *MinioStorage) UploadFile(bucketName string, objectName string, objectFile io.Reader, sizeFile int64, typeFile string, location string) error {
	ctx := context.Background()

	err := m.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})

	if err != nil {

		exists, errBucketExists := m.Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)

		} else {
			log.Fatalln(err)
			return err
		}

	} else {
		log.Printf("Successfully created %s\n", bucketName)

	}

	// Upload the file
	//  minio.PutObjectOptions{ContentType: "application/pdf"}
	info, err := m.Client.PutObject(ctx, bucketName, objectName, objectFile, sizeFile, minio.PutObjectOptions{ContentType: typeFile})

	if err != nil {
		log.Fatalln(err)
		return err
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	return nil
}
