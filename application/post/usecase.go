package post

import "db_technopark/application/models"

type Usecase interface {
	CreatePosts(slug string, id int32, posts models.Posts) (models.Posts, *models.Error)
	GetPostDetails(id int32, query models.PostsRelatedQuery) (models.PostFull, *models.Error)
}
