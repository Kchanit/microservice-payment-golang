package utils

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Kchanit/microservice-payment-golang/internal/core/utils/broker"
	"github.com/Kchanit/microservice-payment-golang/internal/core/utils/payment"
	"github.com/Kchanit/microservice-payment-golang/internal/core/utils/secret"
	"github.com/omise/omise-go"
)

type UtilsFacade struct {
	Omise *omise.Client
	Vault *secret.Vault
}

var lock = &sync.Mutex{}
var singleInstance *UtilsFacade

func NewUtilsFacade(omise *omise.Client, vault *secret.Vault) *UtilsFacade {

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
			vault, err := secret.NewVault(os.Getenv("VAULT_ADDR"), os.Getenv("VAULT_TOKEN"), os.Getenv("VAULT_PATH"))

			if err != nil {
				log.Fatal("vault error", err)
			}

			//Create Omise

			omiseClient, err := payment.NewOmiseClient(vault.GetSecretKey("OMISE_PUBLIC_KEY"), vault.GetSecretKey("OMISE_SECRET_KEY"))

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

func (u *UtilsFacade) SendKafka(topic string, content map[string]interface{}) error {
	err := broker.KafkaProducer(topic, content)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (u *UtilsFacade) ReceiverKafka(topic []string, group string) error {
	err := broker.KafkaConsumer(topic, group)

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
