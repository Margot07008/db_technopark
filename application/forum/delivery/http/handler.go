package http

import (
	"db_technopark/application/forum"
	"db_technopark/application/models"
	"db_technopark/pkg/queryWorker"
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
	router.POST("/api/forum/:path1", handler.CheckPath) ///create
	router.GET("/api/forum/:slug/details", handler.GetForumBySlug)
	router.GET("/api/forum/:slug/users", handler.GerForumUsers)
	router.GET("/api/forum/:slug/threads", handler.GetForumThreads)
}

func (h ForumHandler) CheckPath(ctx *fasthttp.RequestCtx) {
	path := ctx.UserValue("path1")
	if path == "create" {
		h.CreateForum(ctx)
	} else {
		ctx.SetStatusCode(404)
		ctx.SetBody(models.BadRequestErrorBytes)
	}
}

func (h ForumHandler) GetForumThreads(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	query := models.PostsRequestQuery{
		Limit: queryWorker.GetIntParam(ctx, "limit"),
		Since: queryWorker.GetStringParam(ctx, "since"),
		Desc:  queryWorker.GetBoolParam(ctx, "desc"),
	}
	threads, err := h.forumUsecase.GetForumThreads(slug, query)
	if err != nil {
		err.SetToContext(ctx)
		return
	}

	jsonBlob, e := threads.MarshalJSON()
	if e != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)
}

func (h ForumHandler) GerForumUsers(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	query := models.PostsRequestQuery{
		Limit: queryWorker.GetIntParam(ctx, "limit"),
		Since: queryWorker.GetStringParam(ctx, "since"),
		Desc:  queryWorker.GetBoolParam(ctx, "desc"),
	}
	users, err := h.forumUsecase.GetForumUsers(slug, query)
	if err != nil {
		err.SetToContext(ctx)
		return
	}

	jsonBlob, e := users.MarshalJSON()
	if e != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)
}

func (h ForumHandler) GetForumBySlug(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	if slug == "" {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
	}
	foundedForum, err := h.forumUsecase.GetForumBySlug(slug)
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

func (h ForumHandler) CreateForum(ctx *fasthttp.RequestCtx) {
	var buffer models.Forum
	e := buffer.UnmarshalJSON(ctx.PostBody())
	if e != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	createdForum, err := h.forumUsecase.CreateForum(buffer)
	if err != nil && err.StatusCode == 404 {
		err.SetToContext(ctx)
		return
	}
	if err != nil && err.StatusCode == 409 {
		createdForum, err = h.forumUsecase.GetForumBySlug(buffer.Slug)
		if err != nil {
			err.SetToContext(ctx)
			return
		}
		ctx.SetStatusCode(409)
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
