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
}

func (u UserHandler) CreateUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	//e := validation.ValidateNickname(nickname)
	//if e != nil {
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
	returnUser, _ := u.usecase.CreateUser(buffer)
	//if e != nil {
	//	foundUsers := models.Users{}
	//	ctx.SetStatusCode(409)
	//	returnUserByNickname, e := u.usecase.GetUserByNickname(nickname)
	//	if e == nil {
	//		foundUsers = append(foundUsers, returnUserByNickname)
	//	}
	//	returnUserByEmail, e := u.usecase.GetUserByEmail(buffer.Email)
	//	if e == nil && returnUserByNickname.Email != returnUserByEmail.Email {
	//		foundUsers = append(foundUsers, returnUserByEmail)
	//	}
	//	jsonBlob, err := foundUsers.MarshalJSON()
	//	if err != nil {
	//		ctx.SetStatusCode(500)
	//		ctx.SetBody(models.InternalErrorBytes)
	//		return
	//	}
	//	ctx.SetBody(jsonBlob)
	//} else {
		ctx.SetStatusCode(201)
		jsonBlob, err := returnUser.MarshalJSON()
		if err != nil {
			ctx.SetStatusCode(500)
			ctx.SetBody(models.InternalErrorBytes)
			return
		}
		ctx.SetBody(jsonBlob)
	//}

}