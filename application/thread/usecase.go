package thread

import "db_technopark/application/models"

type Usecase interface {
	CreateThread(slug string, thread models.Thread) (models.Thread, *models.Error)
}
