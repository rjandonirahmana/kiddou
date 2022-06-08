package handler

import (
	"kiddou/base"
	"kiddou/domain"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	usecaseUser domain.UsecaseUser
}

func NewUserHandler(useecaseUser domain.UsecaseUser) *UserHandler {
	return &UserHandler{usecaseUser: useecaseUser}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var input domain.RegisterInput

	err := c.BodyParser(&input)
	if err != nil {
		res := base.APIResponse("failed to pass form json", 422, err.Error(), nil)
		return c.Status(422).JSON(res)
	}

	err = validator.New().Struct(&input)
	if err != nil {
		res := base.APIResponse("failed to pass form json", 422, err.Error(), nil)
		return c.Status(422).JSON(res)
	}

	token, err := h.usecaseUser.Register(c.Context(), &input)
	if err != nil {
		res := base.APIResponse("system error failed", 422, err.Error(), nil)
		return c.Status(422).JSON(res)
	}

	res := base.ResponseAPIToken("success create user", 200, "success", nil, token)
	return c.Status(200).JSON(res)

}
