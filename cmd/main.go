package main

import (
	"log"
	"log/slog"
	"os"

	handler "github.com/Kchanit/microservice-payment-golang/internal/adapter/handler/http"
	repository "github.com/Kchanit/microservice-payment-golang/internal/adapter/repository/postgres"
	"github.com/Kchanit/microservice-payment-golang/internal/core/services"
	"github.com/Kchanit/microservice-payment-golang/internal/core/utils"
	"github.com/Kchanit/microservice-payment-golang/internal/core/utils/event"
)

func main() {
	utils.LoadSecret()
	facade := utils.FacadeSingleton()

	repository.ConnectDb(facade.Vault.GetSecretKey("DB_USER"), facade.Vault.GetSecretKey("DB_PASSWORD"), facade.Vault.GetSecretKey("DB_HOST"), facade.Vault.GetSecretKey("DB_NAME"), facade.Vault.GetSecretKey("DB_PORT"))

	userRepo := repository.NewUserRepository(repository.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	transactionRepo := repository.NewTransactionRepository(repository.DB)
	transactionService := services.NewTransactionService(transactionRepo, userRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	omiseService := services.NewOmiseService(userRepo, transactionRepo, transactionService)
	omiseHandler := handler.NewOmiseHandler(omiseService, userService, transactionService)

	// Start the Kafka consumer in a separate goroutine
	go func() {

		err := facade.ReceiverKafka([]string{"user_create"}, "payment", event.Listen)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Init router
	router, err := handler.NewRouter(
		*userHandler,
		*omiseHandler,
		*transactionHandler,
	)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	// Start server
	log.Fatal(router.Start())
}
