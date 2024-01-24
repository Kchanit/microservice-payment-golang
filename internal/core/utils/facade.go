package utils

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/omise/omise-go"
)

type UtilsFacade struct {
	Omise *omise.Client
	Vault *Vault
}

var lock = &sync.Mutex{}
var singleInstance *UtilsFacade

func NewUtilsFacade(omise *omise.Client, vault *Vault) *UtilsFacade {

	return &UtilsFacade{
		Omise: omise,
		Vault: vault,
	}
}

func FacadeSingleton() *UtilsFacade {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating single instance now.")
			// Create Vault
			vault, err := NewVault(os.Getenv("VAULT_ADDR"), os.Getenv("VAULT_TOKEN"), os.Getenv("VAULT_PATH"))

			if err != nil {
				log.Fatal("vault error", err)
			}

			//Create Omise

			omiseClient, err := NewOmiseClient(vault.GetSecretKey("OMISE_PUBLIC_KEY"), vault.GetSecretKey("OMISE_SECRET_KEY"))

			if err != nil {
				log.Fatal("omise error", err)
			}

			singleInstance = NewUtilsFacade(omiseClient, vault)

		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singleInstance
}
