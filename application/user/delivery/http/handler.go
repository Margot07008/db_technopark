package http

import (
	"db_technopark/application/models"
	"db_technopark/application/user"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	usecase user.Usecase
}

func NewUserHandler(router *fasthttprouter.Router, usecase user.Usecase) {
	handler := &UserHandler{
		usecase: usecase,
	}
	router.POST("/api/user/:nickname/create", handler.CreateUser)
	router.GET("/api/user/:nickname/profile", handler.GetUser)
	router.POST("/api/user/:nickname/profile", handler.UpdateUser)
}

func (u UserHandler) UpdateUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	buffer := models.User{}
	body := ctx.PostBody()
	err := buffer.UnmarshalJSON(body)
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	buffer.Nickname = nickname
	returnedUser, e := u.usecase.UpdateUser(buffer)
	if e != nil {
		//e.SetToContext(ctx)
		//return
	}
	jsonBlob, err := returnedUser.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}
	ctx.SetBody(jsonBlob)
}


func (u UserHandler) GetUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	foundedUser, err := u.usecase.GetUserByNickname(nickname)
	if err != nil {
		err.SetToContext(ctx)
		return
	}
	jsonBlob, _ := foundedUser.MarshalJSON()
	ctx.SetBody(jsonBlob)
}

func (u UserHandler) CreateUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	//e := validation.ValidateNickname(nickname)
	//if e != nil
	//	e.SetToContext(ctx)
	//	return
	//}
	buffer := models.User{}
	body := ctx.PostBody()
	err := buffer.UnmarshalJSON(body)
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	buffer.Nickname = nickname
	returnedUser, e := u.usecase.CreateUser(buffer)
	if e != nil {
		foundedUser := models.Users{}
		ctx.SetStatusCode(409)
		existedUserNick, err := u.usecase.GetUserByNickname(nickname)
		if err == nil {
			foundedUser = append(foundedUser, existedUserNick)
		}
		existedUserEmail, err := u.usecase.GetUserByEmail(buffer.Email)
		if err == nil && existedUserNick.Email != existedUserEmail.Email{
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