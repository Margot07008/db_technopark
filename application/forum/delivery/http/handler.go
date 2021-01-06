package http

import (
	"db_technopark/application/forum"
	"db_technopark/application/models"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type ForumHandler struct {
	forumUsecase forum.Usecase
}

func NewForumHandler(router *fasthttprouter.Router, forumUsecase forum.Usecase) {
	handler := &ForumHandler{
		forumUsecase: forumUsecase,
	}
	router.POST("/api/forum/create", handler.CreateForum)
	router.GET("/api/forum/:slug/details", handler.GetForumBySlug)
}

func (u ForumHandler) GetForumBySlug(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	if slug == "" {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
	}
	foundedForum, err := u.forumUsecase.GetForumBySlug(slug)
	if err != nil {
		err.SetToContext(ctx)
		return
	}
	jsonBlob, e := foundedForum.MarshalJSON()
	if e != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}
	ctx.SetBody(jsonBlob)
}

func (u ForumHandler) CreateForum(ctx *fasthttp.RequestCtx) {
	var buffer models.Forum
	body := ctx.PostBody()
	e := buffer.UnmarshalJSON(body)
	if e != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	createdForum, err := u.forumUsecase.CreateForum(buffer)
	if err != nil && err.StatusCode == 404 {
		err.SetToContext(ctx)
		return
	}
	if err != nil && err.StatusCode == 409 {
		//TODO add get forum by slug
	} else {
		ctx.SetStatusCode(201)
	}
	jsonBlob, e := createdForum.MarshalJSON()
	if e != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}
	ctx.SetBody(jsonBlob)
}
