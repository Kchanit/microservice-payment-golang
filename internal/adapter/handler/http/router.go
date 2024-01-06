package handler

import (
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	*fiber.App
}

func NewRouter(userHandler UserHandler) (*Router, error) {
	router := fiber.New()

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

	return &Router{
		router,
	}, nil
}

func (r *Router) Start() error {
	return r.Listen(":8080")
}
