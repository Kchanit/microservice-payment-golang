package handler

import (
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	*fiber.App
}

func NewRouter(userHandler UserHandler) (*Router, error) {
	router := fiber.New()

	router.Get("/", userHandler.Hello)
	router.Get("/users/:id", userHandler.GetUser)

	return &Router{
		router,
	}, nil
}

func (r *Router) Start() error {
	return r.Listen(":8080")
}
