package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Router struct {
	*fiber.App
}

func NewRouter(userHandler UserHandler, omiseHandler OmiseHandler) (*Router, error) {
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

	omise := router.Group("/omise")
	{
		omise.Get("/charge-credit-card/:amount/:token", omiseHandler.ChargeCreditCard)
		omise.Get("/charge-banking/:amount/:source", omiseHandler.ChargeBanking)
		omise.Get("/retrieve-charge/:charge_id", omiseHandler.RetrieveCharge)
		omise.Post("/token", omiseHandler.CreateToken)
		omise.Get("/customers", omiseHandler.ListCustomers)
		omise.Put("/attach-card", omiseHandler.AttchCardToCustomer)
		omise.Post("/webhook", omiseHandler.HandleWebhook)
		omise.Get("/charges", omiseHandler.GetCharges)
	}

	return &Router{
		router,
	}, nil
}

func (r *Router) Start() error {
	return r.Listen(":8080")
}
