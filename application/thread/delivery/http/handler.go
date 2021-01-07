package http

import (
	"db_technopark/application/forum"
	"db_technopark/application/models"
	"db_technopark/application/thread"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"strconv"
)

type ThreadHandler struct {
	threadUsecase thread.Usecase
	forumUsecase  forum.Usecase
}

func NewThreadHandler(router *fasthttprouter.Router, threadUsecase thread.Usecase, forumUsecase forum.Usecase) {
	handler := &ThreadHandler{
		threadUsecase: threadUsecase,
		forumUsecase:  forumUsecase,
	}

	router.POST("/api/forum/:path1/:path2", handler.CheckPath) //:slug/create
	router.GET("/api/thread/:slug_or_id/details", handler.GetThread)
	router.POST("/api/thread/:slug_or_id/details", handler.UpdateThread)
}

func (h ThreadHandler) CheckPath(ctx *fasthttp.RequestCtx) {
	path1 := ctx.UserValue("path1")
	path2 := ctx.UserValue("path2")

	if path1 != "" && path2 == "create" {
		h.CreateThread(ctx)
	} else {
		ctx.SetStatusCode(404)
		ctx.SetBody(models.BadRequestErrorBytes)
	}
}

func (h ThreadHandler) UpdateThread(ctx *fasthttp.RequestCtx) {
	updateThread := models.ThreadUpdate{}

	slugOrId := ctx.UserValue("slug_or_id").(string)
	id, _ := strconv.ParseInt(slugOrId, 10, 64)
	if id == 0 {
		id = -1
		updateThread.Slug = slugOrId
	}
	updateThread.Id = int32(id)

	err := updateThread.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	updatedThread, e := h.threadUsecase.ThreadUpdate(updateThread)
	if e != nil {
		e.SetToContext(ctx)
		return
	}
	jsonBlob, err := updatedThread.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}
	ctx.SetBody(jsonBlob)
}

func (h ThreadHandler) GetThread(ctx *fasthttp.RequestCtx) {
	slugOrId := ctx.UserValue("slug_or_id").(string)
	id, _ := strconv.ParseInt(slugOrId, 10, 64)
	isSlug := false
	if id == 0 {
		isSlug = true
	}
	existingThread, err := h.threadUsecase.GetThreadBySlugOrId(slugOrId, isSlug)
	if err != nil {
		err.SetToContext(ctx)
		return
	}
	jsonBlob, e := existingThread.MarshalJSON()
	if e != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}
	ctx.SetBody(jsonBlob)
}

func (h ThreadHandler) CreateThread(ctx *fasthttp.RequestCtx) {
	createdThread := models.Thread{}
	slug := ctx.UserValue("path1").(string)
	err := createdThread.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	createdThread, e := h.threadUsecase.CreateThread(slug, createdThread)
	if e != nil && e.StatusCode == 409 {
		ctx.SetStatusCode(409)
	} else if e != nil {
		e.SetToContext(ctx)
		return
	} else {
		ctx.SetStatusCode(201)
	}

	jsonBlob, err := createdThread.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)
}
