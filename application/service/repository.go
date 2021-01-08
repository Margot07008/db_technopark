package service

import "db_technopark/application/models"

type Repository interface {
	GetStatus() (models.Status, *models.Error)
	Clear() *models.Error
}
