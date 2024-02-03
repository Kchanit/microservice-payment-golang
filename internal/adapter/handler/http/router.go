package handler

import (
	"time"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"github.com/Kchanit/microservice-payment-golang/internal/core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	omise := router.Group("/omise")
	{
		omise.Post("/charge-credit-card/:userID", omiseHandler.ChargeCreditCard)
		omise.Post("/charge-banking/:userID", omiseHandler.ChargeBanking)
		omise.Get("/retrieve-charge/:charge_id", omiseHandler.RetrieveCharge)
		omise.Post("/token", omiseHandler.CreateToken)
		omise.Get("/customers", omiseHandler.ListCustomers)
		omise.Get("/customers/:customerToken", omiseHandler.GetCustomer)
		omise.Put("/attach-card", omiseHandler.AttachCardToCustomer)
		omise.Post("/webhook", omiseHandler.HandleWebhook)
		omise.Get("/charges", omiseHandler.GetCharges)
		omise.Get("/transactions/:transaction_id", omiseHandler.GetTransaction)
	}

	kafkaroute := router.Group("/kafka")
	{
		kafkaroute.Get("/produce", func(c *fiber.Ctx) error {
			facade := utils.FacadeSingleton()
			err := facade.SendKafka("test", map[string]interface{}{
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

	return &Router{
		router,
	}, nil
}

func (r *Router) Start() error {

	return r.Listen(":8000")
}
