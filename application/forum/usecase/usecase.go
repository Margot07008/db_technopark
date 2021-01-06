package usecase

import (
	"db_technopark/application/forum"
	"db_technopark/application/models"
	"db_technopark/application/user"
)

type forumUsecase struct {
	userRepo  user.Repository
	forumRepo forum.Repository
}

func NewForumUsecase(userRepo user.Repository, forumRepo forum.Repository) forum.Usecase {
	return &forumUsecase{forumRepo: forumRepo, userRepo: userRepo}
}

func (u forumUsecase) CreateForum(forumNew models.Forum) (models.Forum, *models.Error) {
	author, err := u.userRepo.GetByNickname(forumNew.User)
	if err != nil {
		return models.Forum{}, err
	}
	return u.forumRepo.CreateForum(author.Nickname, forumNew)
}

func (u forumUsecase) GetForumBySlug(slug string) (models.Forum, *models.Error) {
	return u.forumRepo.GetForumBySlug(slug)
}
