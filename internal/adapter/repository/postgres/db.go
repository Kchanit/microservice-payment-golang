package repository

import (
	"fmt"
	"log"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb(user, pass, host, dbname, port string) {
	USER := user
	PASS := pass
	HOST := host
	DBNAME := dbname
	DBPORT := port
	fmt.Println(DBPORT)
	fmt.Println("Connecting to database...")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", HOST, USER, PASS, DBNAME, DBPORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&domain.User{}, &domain.Transaction{}, &domain.CardToken{})
	if err != nil {
		fmt.Println("Error while migrating database", err)
	}

	DB = db
	fmt.Println("Successfully connected!")
}
