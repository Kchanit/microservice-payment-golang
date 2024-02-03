package secret

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/vault-client-go"
)

type Vault struct {
	ClientVault *vault.Client
	Path        string
}

func NewVault(endpoint, token, path string) (*Vault, error) {

	// prepare a client with the given base address
	client, err := vault.New(
		vault.WithAddress(endpoint),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatal("vault can't connect", err)
	}
	// authenticate with a root token (insecure)
	if err := client.SetToken(token); err != nil {
		log.Fatal("vault Token Error", err)
	}

	return &Vault{
		ClientVault: client,
		Path:        path,
	}, nil

}

func (v *Vault) GetSecretKey(key string) string {
	ctx := context.Background()

	response, err := v.ClientVault.Read(ctx, v.Path)

	if err != nil {
		log.Fatal("can't get secret key:", err, key, response)
	}

	data, ok := response.Data["data"].(map[string]interface{})

	if !ok {
		log.Fatal("Error Vault Change to Map")
		return ""
	}

	fmt.Printf("%s type is %T\n", key, data[key])
	keyData, ok := data[key].(string)

	if !ok {

		log.Fatal("Error Vault Change to String")
		return ""
	}

	return keyData

}
