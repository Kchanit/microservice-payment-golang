package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"

	"github.com/Kchanit/microservice-payment-golang/internal/core/utils/billing"

	"github.com/Kchanit/microservice-payment-golang/internal/core/utils/payment"
	"github.com/Kchanit/microservice-payment-golang/internal/core/utils/secret"
	"github.com/Kchanit/microservice-payment-golang/internal/core/utils/storage"

	"github.com/joho/godotenv"
	"github.com/omise/omise-go"
)

type UtilsFacade struct {
	Omise *omise.Client
	Vault *secret.Vault
	Minio *storage.MinioStorage
}

var lock = &sync.Mutex{}
var singleInstance *UtilsFacade

func NewUtilsFacade(omise *omise.Client, vault *secret.Vault, minio *storage.MinioStorage) *UtilsFacade {

	return &UtilsFacade{
		Omise: omise,
		Vault: vault,
		Minio: minio,
	}
}

func FacadeSingleton() *UtilsFacade {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating single instance now.")
			// Create Vault
			vault, err := secret.NewVault(os.Getenv("VAULT_ADDR"), os.Getenv("VAULT_TOKEN"), os.Getenv("VAULT_PATH"))

			if err != nil {
				log.Fatal("vault error", err)
			}

			//Create Omise

			omiseClient, err := payment.NewOmiseClient(vault.GetSecretKey("OMISE_PUBLIC_KEY"), vault.GetSecretKey("OMISE_SECRET_KEY"))

			if err != nil {
				log.Fatal("omise error", err)
			}

			// Create Minio

			minioConfig := storage.MinioSetup{
				Endpoint:        vault.GetSecretKey("MINIO_ENDPOINT"),
				AccessKeyID:     vault.GetSecretKey("MINIO_ACCESS_KEY"),
				SecretAccessKey: vault.GetSecretKey("MINIO_SECRET_KEY"),
				UseSSL:          true,
			}

			minio, err := storage.NewMinio(minioConfig)

			if err != nil {
				log.Fatal("minio error", err)
			}

			singleInstance = NewUtilsFacade(omiseClient, vault, minio)

		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singleInstance
}

func (u *UtilsFacade) GenerateInvoice(customer domain.User, products []domain.Product, ref string) (io.Reader, int64, error) {
	bill, size, err := billing.GenerateInvoice(customer, products, ref)
	if err != nil {
		return nil, size, err
	}

	return bill, size, nil

}

func LoadSecret() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
