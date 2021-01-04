package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/forum"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"net/http"
)

type ResumeHandler struct {
	UseCaseResume forum.UseCase
}

const resumePath = "resume/"

func NewRest(router *gin.RouterGroup,
	useCaseResume forum.UseCase,
	AuthRequired gin.HandlerFunc) *ResumeHandler {
	rest := &ResumeHandler{
		UseCaseResume: useCaseResume,
	}
	rest.routes(router, AuthRequired)
	return rest
}

func (r *ResumeHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.POST("/create", r.CreateForum)
	router.POST("/:slug/create", r.CreateThread)
	router.GET("/:slug/details", r.GetForumBySlug)
	router.GET("/:slug/threads", r.GetAllForumTreads)
	router.GET("/:slug/users", r.GetAllForumUsers)
}

func (r *ResumeHandler) CreateForum(ctx *gin.Context) {
	var template models.Forum
	if err := ctx.ShouldBindBodyWith(&template, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	//if err := common.ReqValidation(template); err != nil {
	//	ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
	//	return
	//}

	result, err := r.UseCaseResume.CreateForum(template)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}
func (r *ResumeHandler) CreateThread(ctx *gin.Context) {
	var req struct {
		Slug string `uri:"slug" binding:"required,string"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var template models.Thread
	if err := ctx.ShouldBindBodyWith(&template, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	//if err := common.ReqValidation(template); err != nil {
	//	ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
	//	return
	//}

	result, err := r.UseCaseResume.CreateThread(req.Slug, template)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (r *ResumeHandler) GetForumBySlug(ctx *gin.Context) {
	var req struct {
		Slug string `uri:"slug" binding:"required,string"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := r.UseCaseResume.GetForumBySlug(req.Slug)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (r *ResumeHandler) GetAllForumTreads(ctx *gin.Context) {
	var req struct {
		Slug string `uri:"slug" binding:"required,string"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var params models.ForumParams

	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	result, err := r.UseCaseResume.GetAllForumTreads(req.Slug, params)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (r *ResumeHandler) GetAllForumUsers(ctx *gin.Context) {
	var req struct {
		Slug string `uri:"slug" binding:"required,string"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var params models.ForumParams

	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	result, err := r.UseCaseResume.GetAllForumUsers(req.Slug, params)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}