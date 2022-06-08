package handler

import (
	"kiddou/base"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type middleWare struct {
	authRedis base.AuthRedis
}

func NewMiddleware(authredis base.AuthRedis) *middleWare {
	return &middleWare{authRedis: authredis}
}

func (m *middleWare) GetTokenFromHeaderBearer(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			base.APIResponse(c, "unauthorized", 433, "error authentication", nil)
			return

		}

		if strings.HasPrefix(authorization, "Bearer ") {
			token = authorization[7:]
		}

		if token == "" {
			base.APIResponse(c, "unauthorized", 433, "error authentication", nil)
			return
		}

		user, err := m.authRedis.Authentication(c, token)
		if err != nil {
			base.APIResponse(c, "unauthorized", 433, err.Error(), nil)
			return
		}

		log.Println(user)
		c.Set("user", user)
		next(c)
	}
}
