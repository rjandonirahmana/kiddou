package handler

import (
	"io/ioutil"
	"kiddou/base"
	"kiddou/domain"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type HandlerVideos struct {
	useCaseVideo domain.UsecaseVideos
}

func NewHandlerVideo(usecaseVideo domain.UsecaseVideos) *HandlerVideos {
	return &HandlerVideos{useCaseVideo: usecaseVideo}
}

func (h *HandlerVideos) CreateVideosAdmin(c *gin.Context) {
	user := c.MustGet("user").(base.AuthToken)
	log.Println(user)
	if user.Role != "admin" {
		base.APIResponse(c, "you are not admin", 422, "unauthorized admin", nil)
		return
	}

	var input domain.InserVideosRequest
	err := c.MustBindWith(&input, binding.Form)
	if err != nil {
		base.APIResponse(c, "failed parse form file", 422, err.Error(), nil)
		return

	}

	err = validator.New().Struct(&input)
	if err != nil {
		base.APIResponse(c, "failed to pass form json", 422, err.Error(), nil)
		return

	}

	file, err := c.FormFile("videos")
	if err != nil {
		base.APIResponse(c, "failed parse form file videoo", 500, err.Error(), nil)

	}

	multiPartFile, err := file.Open()
	if err != nil {
		base.APIResponse(c, "failed parse form file videoo", 500, err.Error(), nil)
		return
	}
	defer multiPartFile.Close()

	fileByte, err := ioutil.ReadAll(multiPartFile)
	if err != nil {
		base.APIResponse(c, "failed parse form file videoo", 500, err.Error(), nil)
		return

	}

	err = h.useCaseVideo.InsertVideos(c, fileByte, &input)
	if err != nil {
		base.APIResponse(c, "failed parse form file videoo", 500, err.Error(), nil)
		return

	}

	base.APIResponse(c, "success", 200, "success", nil)
	return

}
