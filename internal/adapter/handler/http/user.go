package handler

import (
	"fmt"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"github.com/Kchanit/microservice-payment-golang/internal/core/ports"
	"github.com/go-playground/validator/v10"
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

type userInput struct {
	Name  string `validate:"required" json:"name"`
	Email string `validate:"required,email" json:"email"`
}

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		fmt.Println("Error while getting all users", err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		fmt.Println("Error while getting user", err)
		return c.Status(fiber.StatusBadRequest).SendString("User not found")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	user := new(domain.User)
	if err := c.BodyParser(user); err != nil {
		fmt.Println("Error while parsing user", err)
		return c.SendStatus(400)
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		fmt.Println("Error while validating user", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: "Validation failed",
			Errors:  extractValidationErrors(err),
		})
	}

	user, err := h.userService.CreateUser(user)
	if err != nil {
		fmt.Println("Error while creating user", err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if the user exists
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		fmt.Println("Error while getting user", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User not found"})
	}

	// Parse the request body
	var userInput userInput
	if err := c.BodyParser(&userInput); err != nil {
		fmt.Println("Error while parsing user", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(userInput); err != nil {
		fmt.Println("Error while validating user", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: "Validation failed",
			Errors:  extractValidationErrors(err),
		})
	}

	// Update user fields
	user.Name = userInput.Name
	user.Email = userInput.Email

	// Update user
	updatedUser, err := h.userService.UpdateUser(user, id)
	if err != nil {
		fmt.Println("Error while updating user", err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(updatedUser)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if the user exists
	_, err := h.userService.GetUserByID(id)
	if err != nil {
		fmt.Println("Error while getting user", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User not found"})
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		fmt.Println("Error while deleting user", err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}

// function to extract validation errors
func extractValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	for _, fieldErr := range err.(validator.ValidationErrors) {
		errors[fieldErr.Field()] = getErrorMessage(fieldErr)
	}
	return errors
}

// function to get a user-friendly error message from a ValidationErrors instance
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is not valid", err.Field())
	}
}
