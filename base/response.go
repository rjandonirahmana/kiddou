package base

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

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	NewRensponse := Response{
		Meta: meta,
		Data: data,
	}

	return NewRensponse
}

func ResponseAPIToken(message string, code int, status string, data interface{}, GetToken string) ResponseToken {
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

	return NewRensponse
}
