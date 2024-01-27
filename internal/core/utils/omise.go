package utils

import (
	"fmt"

	"github.com/omise/omise-go"
)

func NewOmiseClient(pubKey, secretKey string) (*omise.Client, error) {
	OmisePublicKey := pubKey
	OmiseSecretKey := secretKey
	client, e := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	if e != nil {
		fmt.Println(e)
	}
	return client, e
}
