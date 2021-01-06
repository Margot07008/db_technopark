package server

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	_ "github.com/valyala/fasthttp"

	deliveryUser "db_technopark/application/user/delivery/http"
	repositoryUser "db_technopark/application/user/repository"
	usecaseUser "db_technopark/application/user/usecase"

	deliveryForum "db_technopark/application/forum/delivery/http"
	repositoryForum "db_technopark/application/forum/repository"
	usecaseForum "db_technopark/application/forum/usecase"
)

type server struct {
	Host   string
	router *fasthttprouter.Router
}

func NewServer(host string, conn *pgx.ConnPool) *server {

	userRepo := repositoryUser.NewPgUserRepository(conn)
	forumRepo := repositoryForum.NewPgForumRepository(conn)

	userUsecase := usecaseUser.NewUserUsecase(userRepo)
	forumUsecase := usecaseForum.NewForumUsecase(userRepo, forumRepo)

	router := fasthttprouter.New()

	deliveryUser.NewUserHandler(router, userUsecase)
	deliveryForum.NewForumHandler(router, forumUsecase)

	return &server{
		Host:   host,
		router: router,
	}
}

func (s server) ListenAndServe() error {
	return fasthttp.ListenAndServe(s.Host, DefaultHeaders(s.router.Handler))
}
