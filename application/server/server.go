package server

import (
	"db_technopark/application/user/usecase"
	"github.com/buaazp/fasthttprouter"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	_ "github.com/valyala/fasthttp"

	"db_technopark/application/user/delivery/http"
	"db_technopark/application/user/repository"
)

type server struct {
	Host   string
	router *fasthttprouter.Router
}

func NewServer(host string, conn *pgx.ConnPool) *server {

	userRepo := repository.NewPgUserRepository(conn)
	userUsecase := usecase.NewUserUsecase(userRepo)

	router := fasthttprouter.New()

	http.NewUserHandler(router, userUsecase)

	return &server{
		Host:   host,
		router: router,
	}
}

func (s server) ListenAndServe() error {
	return fasthttp.ListenAndServe(s.Host, DefaultHeaders(s.router.Handler))
}