package forum

import "db_technopark/application/models"

type Usecase interface {
	CreateForum(forum models.Forum) (models.Forum, *models.Error)
	GetForumBySlug(slug string) (models.Forum, *models.Error)
	GetForumUsers(slug string, query models.PostsRequestQuery) (models.Users, *models.Error)
}
