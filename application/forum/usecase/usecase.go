package usecase

import (
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/forum"
)

type UseCase struct {
	infoLogger       *logger.Logger
	errorLogger      *logger.Logger
	strg             forum.Repository
}

func (u *UseCase) CreateThread(slug string, template models.Thread) (*models.Thread, error) {
	panic("implement me")
}

func (u *UseCase) GetForumBySlug(slug string) (*models.Forum, error) {
	panic("implement me")
}

func (u *UseCase) GetAllForumTreads(slug string, params models.ForumParams) ([]models.Thread, error) {
	panic("implement me")
}

func (u *UseCase) GetAllForumUsers(slug string, params models.ForumParams) ([]models.User, error) {
	panic("implement me")
}

func NewUseCase(infoLogger *logger.Logger,
	errorLogger *logger.Logger,
	strg forum.Repository) forum.UseCase {
	usecase := UseCase{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		strg:        strg,
	}
	return &usecase
}

func (u *UseCase) CreateForum(template models.Forum) (*models.Forum, error) {
	return u.strg.CreateForum(template)
}

