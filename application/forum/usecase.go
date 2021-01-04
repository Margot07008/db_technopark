package forum

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
)

type UseCase interface {
	CreateForum(template models.Forum) (*models.Forum, error)
	CreateThread(slug string, template models.Thread) (*models.Thread, error)
	GetForumBySlug(slug string) (*models.Forum, error)
	GetAllForumTreads(slug string, params models.ForumParams) ([]models.Thread, error)
	GetAllForumUsers(slug string, params models.ForumParams) ([]models.User, error)
}
