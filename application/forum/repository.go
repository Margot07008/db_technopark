package forum

import "db_technopark/application/models"

type Repository interface {
	CreateForum(user string, forum models.Forum) (models.Forum, *models.Error)
}
