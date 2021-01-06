package http

import (
	"db_technopark/application/models"
	"db_technopark/application/user"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"regexp"
)

type UserHandler struct {
	usecase user.Usecase
}

var nicknamePattern = regexp.MustCompile("^[A-Za-z0-9_.]+$")

func ValidateNickname(nickname string) *models.Error {
	if !nicknamePattern.MatchString(nickname) {
		return models.NewError(400, models.BadRequestError, "invalid nickname")
	}
	return nil
}

func NewUserHandler(router *fasthttprouter.Router, usecase user.Usecase) {
	handler := &UserHandler{
		usecase: usecase,
	}
	router.POST("/api/user/:nickname/create", handler.CreateUser)
	router.GET("/api/user/:nickname/profile", handler.GetUser)
	router.POST("/api/user/:nickname/profile", handler.UpdateUser)
}

func (h UserHandler) UpdateUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	e := ValidateNickname(nickname)
	if e != nil {
		e.SetToContext(ctx)
		return
	}
	buffer := models.User{}
	err := buffer.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	buffer.Nickname = nickname
	returnedUser, e := h.usecase.UpdateUser(buffer)
	if e != nil {
		e.SetToContext(ctx)
		return
	}
	jsonBlob, err := returnedUser.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}
	ctx.SetBody(jsonBlob)
}

func (h UserHandler) GetUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	foundedUser, err := h.usecase.GetUserByNickname(nickname)
	if err != nil {
		err.SetToContext(ctx)
		return
	}
	jsonBlob, _ := foundedUser.MarshalJSON()
	ctx.SetBody(jsonBlob)
}

func (h UserHandler) CreateUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	e := ValidateNickname(nickname)
	if e != nil {
		e.SetToContext(ctx)
		return
	}
	buffer := models.User{}
	err := buffer.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	buffer.Nickname = nickname
	returnedUser, e := h.usecase.CreateUser(buffer)
	if e != nil {
		foundedUser := models.Users{}
		ctx.SetStatusCode(409)
		existedUserNick, err := h.usecase.GetUserByNickname(nickname)
		if err == nil {
			foundedUser = append(foundedUser, existedUserNick)
		}
		existedUserEmail, err := h.usecase.GetUserByEmail(buffer.Email)
		if err == nil && existedUserNick.Email != existedUserEmail.Email {
			foundedUser = append(foundedUser, existedUserEmail)
		}
		jsonBlob, e := foundedUser.MarshalJSON()
		if e != nil {
			ctx.SetStatusCode(500)
			ctx.SetBody(models.InternalErrorBytes)
			return
		}
		ctx.SetBody(jsonBlob)
	} else {
		ctx.SetStatusCode(201)
		jsonBlob, err := returnedUser.MarshalJSON()
		if err != nil {
			ctx.SetStatusCode(500)
			ctx.SetBody(models.InternalErrorBytes)
			return
		}
		ctx.SetBody(jsonBlob)
	}
}
