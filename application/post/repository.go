package post

import "db_technopark/application/models"

type Repository interface {
	CreatePosts(posts models.Posts, thread models.Thread) (models.Posts, *models.Error)
	GetById(id int64) (models.Post, *models.Error)
	GetMany(thread models.Thread, query models.PostsRequestQuery) (models.Posts, *models.Error)
	Update(post models.Post, postUpdate models.PostUpdate) (models.Post, *models.Error)
}
