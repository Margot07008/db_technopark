package thread

import "db_technopark/application/models"

type Repository interface {
	CreateThread(forum models.Forum, user models.User, threadNew models.Thread) (models.Thread, *models.Error)
	GetThreadsBySlug(forum models.Forum, query models.PostsRequestQuery) (models.Threads, *models.Error)
}
