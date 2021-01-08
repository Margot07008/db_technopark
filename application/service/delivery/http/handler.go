package http

import (
	"db_technopark/application/models"
	"db_technopark/application/service"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type ServiceHandler struct {
	serviceUsecase service.Usecase
}

func NewServiceHandler(router *fasthttprouter.Router, serviceUsecase service.Usecase) {
	handler := &ServiceHandler{
		serviceUsecase: serviceUsecase,
	}
	router.POST("/api/service/clear", handler.Clear)
	router.GET("/api/service/status", handler.GetStatus)
}

func (u ServiceHandler) GetStatus(ctx *fasthttp.RequestCtx) {
	dbStatus, e := u.serviceUsecase.GetDBStatus()
	if e != nil {
		e.SetToContext(ctx)
		return
	}
	jsonBlob, err := dbStatus.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
	}
	ctx.SetBody(jsonBlob)
}

func (u ServiceHandler) Clear(ctx *fasthttp.RequestCtx) {
	err := u.serviceUsecase.ClearDB()
	if err != nil {
		err.SetToContext(ctx)
	}
}
