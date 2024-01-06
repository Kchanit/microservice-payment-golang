package main

import (
	"log"
	"log/slog"
	"os"

	handler "github.com/Kchanit/microservice-payment-golang/internal/adapter/handler/http"
	repository "github.com/Kchanit/microservice-payment-golang/internal/adapter/repository/mysql"
	"github.com/Kchanit/microservice-payment-golang/internal/core/services"
	"github.com/joho/godotenv"
)

func main() {
	LoadEnv()

	repository.ConnectDb()

	userRepo := repository.NewUserRepository(repository.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Init router
	router, err := handler.NewRouter(
		*userHandler,
	)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	// Start server
	log.Fatal(router.Start())
}

func LoadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
