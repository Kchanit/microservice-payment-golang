package handler

import (
	"time"

	"log"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"github.com/Kchanit/microservice-payment-golang/internal/core/utils"
	"github.com/Kchanit/microservice-payment-golang/internal/core/utils/broker"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type Router struct {
	*fiber.App
}

func NewRouter(userHandler UserHandler, omiseHandler OmiseHandler, transactionHandler TransactionHandler) (*Router, error) {
	router := fiber.New()

	router.Use(cors.New(
		cors.Config{
			AllowOrigins:     "*",
			AllowCredentials: true,
		}))

	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Yo, World 👋!")
	})

	user := router.Group("/users")
	{
		user.Post("/", userHandler.CreateUser)
		user.Get("/", userHandler.GetAllUsers)
		user.Get("/:id", userHandler.GetUserByID)
		user.Put("/:id", userHandler.UpdateUser)
		user.Delete("/:id", userHandler.DeleteUser)
	}

	router.Get("/transactions", transactionHandler.GetAllTransactions)

	omisePayment := router.Group("/omise")
	{
		omisePayment.Post("/charge-credit-card/:userID", omiseHandler.ChargeCreditCard)
		omisePayment.Post("/charge-banking/:userID", omiseHandler.ChargeBanking)
		omisePayment.Get("/retrieve-charge/:charge_id", omiseHandler.RetrieveCharge)
		omisePayment.Post("/token", omiseHandler.CreateToken)
		omisePayment.Get("/customers", omiseHandler.ListCustomers)
		omisePayment.Get("/customers/:customerToken", omiseHandler.GetCustomer)
		omisePayment.Put("/attach-card", omiseHandler.AttachCardToCustomer)
		omisePayment.Post("/webhook", omiseHandler.HandleWebhook)
		omisePayment.Get("/charges", omiseHandler.GetCharges)
		omisePayment.Get("/transactions/:transaction_id", omiseHandler.GetTransaction)
	}

	kafkaroute := router.Group("/kafka")
	{
		kafkaroute.Get("/produce", func(c *fiber.Ctx) error {

			err := broker.KafkaProducer("test", "key", map[string]interface{}{
				"message": "Hello World!",
			})
			if err != nil {
				return c.SendString("Error")
			}

			return c.SendString("Yo, World 👋!")
		})
	}
	billing := router.Group("/billing")
	{
		billing.Get("/invoice", func(c *fiber.Ctx) error {
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
			facade := utils.FacadeSingleton()
			bill, size, err := facade.GenerateInvoice(customer, products, ref)
			if err != nil {
				return c.SendString("Error Generate Invoice")
			}
			bucketName := "pixelmanstorage"
			objectName := "invoices/" + ref + time.Now().Format("2006-01-02") + ".pdf"

			// Get Minio client instance
			err = facade.Minio.UploadFile(bucketName, objectName, bill, size, "application/pdf", "us-east-1")
			if err != nil {
				return c.SendString("Error Get Minio Client")
			}

			return c.SendString("Yo, World 👋!")
		})

	}

	admin := router.Group("/admin")
	{
		admin.Get("/balance", func(c *fiber.Ctx) error {
			facade := utils.FacadeSingleton()
			result := &omise.Balance{}
			err := facade.Omise.Do(result, &operations.RetrieveBalance{})
			if err != nil {
				log.Fatalln(err)
			}

			return c.JSON(result)
		})

	}

	return &Router{
		router,
	}, nil
}

func (r *Router) Start() error {

	return r.Listen(":8000")
}
