package handler

import (
	"fmt"

	"github.com/Kchanit/microservice-payment-golang/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := h.userService.GetUser(id)
	if err != nil {
		fmt.Println("Error while getting user", err)
		return ctx.SendStatus(404)
	}
	return ctx.JSON(user)
}

func (h *UserHandler) Hello(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello World!")
}
