package base

import "github.com/gin-gonic/gin"

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type ResponseToken struct {
	Meta  Meta        `json:"meta"`
	Data  interface{} `json:"data"`
	Token string      `json:"token"`
}

func APIResponse(c *gin.Context, message string, code int, status string, data interface{}) {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	NewRensponse := Response{
		Meta: meta,
		Data: data,
	}
	c.JSON(code, NewRensponse)
}

func ResponseAPIToken(c *gin.Context, message string, code int, status string, data interface{}, GetToken string) {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	NewRensponse := ResponseToken{
		Meta:  meta,
		Data:  data,
		Token: GetToken,
	}

	c.JSON(code, NewRensponse)
}
