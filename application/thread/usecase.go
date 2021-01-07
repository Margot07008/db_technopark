package thread

import "db_technopark/application/models"

type Usecase interface {
	CreateThread(slug string, thread models.Thread) (models.Thread, *models.Error)
	GetThreadBySlugOrId(data string, isSlug bool) (models.Thread, *models.Error)
	ThreadUpdate(threadUpdate models.ThreadUpdate) (models.Thread, *models.Error)
}
