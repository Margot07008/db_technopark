package user

import "db_technopark/application/models"

type Repository interface {
	Create(userNew models.User) (models.User, *models.Error)
	GetByNickname(nickname string) (models.User, *models.Error)
	GetByEmail(email string) (models.User, *models.Error)
	Update(userUdp models.User) (models.User, *models.Error)
	GetByForum(forum models.Forum, query models.PostsRequestQuery) (models.Users, *models.Error)
}
