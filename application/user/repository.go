package user

import "db_technopark/application/models"

type Repository interface {
	Create(userNew models.User) (models.User, *models.Error)
}
