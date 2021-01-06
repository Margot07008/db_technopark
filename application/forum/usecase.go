package forum

import "db_technopark/application/models"

type Usecase interface {
	CreateForum(forum models.Forum) (models.Forum, *models.Error)
}
