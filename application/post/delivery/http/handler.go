package http

import (
	"db_technopark/application/models"
	"db_technopark/application/post"
	"db_technopark/application/vote"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type PostHandler struct {
	postUsecase post.Usecase
	voteUsecase vote.Usecase
}

func NewPostHandler(router *fasthttprouter.Router, posrUsecase post.Usecase, voteUsecase vote.Usecase) {
	handler := &PostHandler{
		postUsecase: posrUsecase,
		voteUsecase: voteUsecase,
	}

	router.POST("/api/thread/:slug_or_id/create", handler.CreatePosts)
	router.GET("/api/post/:id/details", handler.GetOnePost)
	router.POST("/api/thread/:slug_or_id/vote", handler.CreateVote)
}

func (p PostHandler) CreateVote(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	id, _ := strconv.ParseInt(slugOrID, 10, 64)
	createdVote := models.Vote{}

	if id == 0 {
		createdVote.Thread = -1
		createdVote.ThreadSlug = slugOrID
	} else {
		createdVote.Thread = int32(id)
	}

	err := createdVote.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}

	existingThread, e := p.voteUsecase.UpsertVote(createdVote)

	if e != nil {
		e.SetToContext(ctx)
		return
	}
	jsonBlob, err := existingThread.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)
}

func (p PostHandler) GetOnePost(ctx *fasthttp.RequestCtx) {
	var id int64 = -1
	id, _ = strconv.ParseInt(ctx.UserValue("id").(string), 10, 64)
	if id == -1 {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	queryParams := strings.Split(string(ctx.URI().QueryArgs().Peek("related")), ",")

	var query models.PostsRelatedQuery

	for _, param := range queryParams {
		if param == "user" {
			query.NeedAuthor = true
		} else if param == "forum" {
			query.NeedForum = true
		} else if param == "thread" {
			query.NeedThread = true
		}
	}

	existingPost, err := p.postUsecase.GetPostDetails(int32(id), query)
	if err != nil {
		err.SetToContext(ctx)
		return
	}

	jsonBlob, e := existingPost.MarshalJSON()
	if e != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)
}

func (p PostHandler) CreatePosts(ctx *fasthttp.RequestCtx) {
	slugOrId := ctx.UserValue("slug_or_id").(string)
	id, _ := strconv.ParseInt(slugOrId, 10, 64)
	if id == 0 {
		id = -1
	}

	posts := models.Posts{}
	err := posts.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}

	posts, e := p.postUsecase.CreatePosts(slugOrId, int32(id), posts)
	if e != nil {
		e.SetToContext(ctx)
		return
	}

	jsonBlob, err := posts.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetStatusCode(201)
	ctx.SetBody(jsonBlob)
}
