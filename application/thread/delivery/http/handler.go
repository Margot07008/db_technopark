package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/thread"
	"net/http"
)

type UserHandler struct {
	UserUseCase thread.UseCase
}

type Request struct {
	ID   int32 `uri:"identification" binding:"int32"`
	Slug string `uri:"identification" binding:"string"`
}

func NewRest(router *gin.RouterGroup, useCase thread.UseCase) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router)
	return rest
}

func (u *UserHandler) routes(router *gin.RouterGroup) {
	router.POST("/:identification/create", u.CreateThread)
	router.GET("/:identification/details", u.GetThreadDetails)
	router.POST("/:identification/details", u.UpdateThread)
	router.GET("/:identification/posts", u.GetPostsThread)
	router.GET("/:identification/vote", u.VoteOnThread)
}

func (u *UserHandler) GetThreadDetails(ctx *gin.Context) {
	var req Request
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	thread, err := u.UserUseCase.GetThreadDetails(req.ID, req.Slug)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, thread)
}

func (u *UserHandler) CreateThread(ctx *gin.Context) {
	var req Request
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if (Request{} == req) {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.UriErrorThread})
		return
	}

	var thread models.Thread
	if err := ctx.ShouldBindJSON(&thread); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	if req.Slug != "" {
		thread.Slug = req.Slug
	}
	if req.ID != 0 {
		thread.ID = req.ID
	}
	//if err := common.ReqValidation(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
	//	return
	//}

	newThread, err := u.UserUseCase.CreateThread(thread)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}

	ctx.JSON(http.StatusOK, newThread)
}

func (u *UserHandler) UpdateThread(ctx *gin.Context) {
	var req Request
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if (Request{} == req) {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.UriErrorThread})
		return
	}

	var thread models.Thread
	if err := ctx.ShouldBindJSON(&thread); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	if req.Slug != "" {
		thread.Slug = req.Slug
	}
	if req.ID != 0 {
		thread.ID = req.ID
	}
	//if err := common.ReqValidation(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
	//	return
	//}

	newThread, err := u.UserUseCase.UpdateThread(thread)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}

	ctx.JSON(http.StatusOK, newThread)
}

func (u *UserHandler) GetPostsThread(ctx *gin.Context) {
	var req Request
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if (Request{} == req) {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.UriErrorThread})
		return
	}

	var query models.ThreadParams
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	if req.Slug != "" {
		query.SlugOrID = req.Slug
	} else if req.ID != 0 {
		query.SlugOrID = string(req.ID)
	}
	//if err := common.ReqValidation(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
	//	return
	//}

	threads, err := u.UserUseCase.GetPostsThread(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}

	ctx.JSON(http.StatusOK, threads)
}

func (u *UserHandler) VoteOnThread(ctx *gin.Context) {
	var req Request
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if (Request{} == req) {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.UriErrorThread})
		return
	}

	var vote models.Vote
	if err := ctx.ShouldBindJSON(&vote); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	if req.Slug != "" {
		vote.SlugOrID = req.Slug
	} else if req.ID != 0 {
		vote.SlugOrID = string(req.ID)
	}

	//if err := common.ReqValidation(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
	//	return
	//}

	threads, err := u.UserUseCase.VoteOnThread(vote)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}

	ctx.JSON(http.StatusOK, threads)
}