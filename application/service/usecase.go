package service

import "db_technopark/application/models"

type Usecase interface {
	GetDBStatus() (models.Status, *models.Error)
	ClearDB() *models.Error
}
