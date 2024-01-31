package utils

import (
	"log"
	"os"
	"time"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
)

func test() {
	products := []domain.Product{
		{
			Name:        "T-shirt",
			Description: "White t-shirt, Size L",
			Price:       9300,
			Quantity:    1,
		},
		{
			Name:        "Test2",
			Description: "Temp Description",
			Price:       4800,
			Quantity:    5,
		},
	}
	customer := domain.User{
		Name: "John Doe",
		Addresses: []domain.Address{
			{
				Address:    "89 somewhere",
				PostalCode: "12345",
				City:       "Phuket",
				Country:    "Thailand",
			},
		},
	}
	ref := "trx10281723"
	outputName, err := GenerateInvoice(customer, products, ref)
	if err != nil {
		log.Fatal(err)
	}

	bucketName := "pixelmanstorage"
	objectName := "invoices/" + ref + time.Now().Format("2006-01-02") + ".pdf"

	// Get Minio client instance
	minioClientInstance, err := GetMinioClient()
	if err != nil {
		log.Fatal(err)
	}

	// Upload PDF file to Minio
	err = minioClientInstance.UploadImage(bucketName, objectName, outputName)
	if err != nil {
		log.Fatal(err)
	}

	os.Remove(outputName)
	log.Println("PDF file uploaded successfully.")
}
