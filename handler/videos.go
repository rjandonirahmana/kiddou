package handler

import (
	"io/ioutil"
	"kiddou/base"
	"kiddou/domain"
	"strconv"

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

func (h *HandlerVideos) SubscribersVideo(c *gin.Context) {
	user := c.MustGet("user").(base.AuthToken)

	videoID := c.Request.FormValue("video_id")
	vidInt, err := strconv.Atoi(videoID)
	if err != nil {
		base.APIResponse(c, "video_id only accept integer", 403, err.Error(), nil)
		return

	}

	err = h.useCaseVideo.SubscribtionVideo(c, user.UserID, vidInt)
	if err != nil {
		base.APIResponse(c, "error system", 500, err.Error(), nil)
		return

	}

	base.APIResponse(c, "success to subscribe", 200, "success", nil)
	return

}

func (h *HandlerVideos) StatusSUbscribe(c *gin.Context) {
	user := c.MustGet("user").(base.AuthToken)

	videoID := c.Param("id")
	vidInt, err := strconv.Atoi(videoID)
	if err != nil {
		base.APIResponse(c, "video_id only accept integer", 403, err.Error(), nil)
		return

	}

	res, err := h.useCaseVideo.SubsribesStatus(c, user.UserID, vidInt)
	if err != nil {
		base.APIResponse(c, "error system", 500, err.Error(), nil)
		return

	}

	base.APIResponse(c, "success to subscribe", 200, "success", res)
	return
}

func (h *HandlerVideos) RenewSubscribe(c *gin.Context) {
	user := c.MustGet("user").(base.AuthToken)

	videoID := c.Request.FormValue("video_id")
	vidInt, err := strconv.Atoi(videoID)
	if err != nil {
		base.APIResponse(c, "video_id only accept integer", 403, err.Error(), nil)
		return

	}

	err = h.useCaseVideo.SubscribtionVideo(c, user.UserID, vidInt)
	if err != nil {
		base.APIResponse(c, "error system", 500, err.Error(), nil)
		return

	}

	base.APIResponse(c, "success to subscribe", 200, "success", nil)
	return
}
