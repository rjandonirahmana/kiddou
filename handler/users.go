package handler

import (
	"kiddou/base"
	"kiddou/domain"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	usecaseUser domain.UsecaseUser
}

func NewUserHandler(useecaseUser domain.UsecaseUser) *UserHandler {
	return &UserHandler{usecaseUser: useecaseUser}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input domain.RegisterInput

	err := c.MustBindWith(&input, binding.Form)
	if err != nil {
		base.APIResponse(c, "failed to pass form json", 422, err.Error(), nil)
		return
	}

	err = validator.New().Struct(&input)
	if err != nil {
		base.APIResponse(c, "failed to pass form json", 422, err.Error(), nil)
		return
	}

	token, err := h.usecaseUser.Register(c, &input)
	if err != nil {
		base.APIResponse(c, "system error failed", 422, err.Error(), nil)
		return

	}

	base.ResponseAPIToken(c, "success create user", 200, "success", nil, token)
	return

}

func (h *UserHandler) Login(c *gin.Context) {
	var input domain.LoginInput

	err := c.MustBindWith(&input, binding.Form)
	if err != nil {
		base.APIResponse(c, "failed to pass form json", 422, err.Error(), nil)
		return
	}

	err = validator.New().Struct(&input)
	if err != nil {
		base.APIResponse(c, "failed to pass form json", 422, err.Error(), nil)
		return
	}

	token, err := h.usecaseUser.Login(c, input.Email, input.Password)
	if err != nil {
		base.APIResponse(c, "failed to login", 422, err.Error(), nil)
		return
	}

	base.ResponseAPIToken(c, "success login", 200, "success", nil, token)
	return
}
