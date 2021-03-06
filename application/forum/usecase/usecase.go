package usecase

import (
	"db_technopark/application/forum"
	"db_technopark/application/models"
	"db_technopark/application/thread"
	"db_technopark/application/user"
)

type forumUsecase struct {
	userRepo   user.Repository
	forumRepo  forum.Repository
	threadRepo thread.Repository
}

func NewForumUsecase(userRepo user.Repository, forumRepo forum.Repository, threadRepo thread.Repository) forum.Usecase {
	return &forumUsecase{forumRepo: forumRepo, userRepo: userRepo, threadRepo: threadRepo}
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

func (u forumUsecase) GetForumUsers(slug string, query models.PostsRequestQuery) (models.Users, *models.Error) {
	foundedForum, err := u.forumRepo.GetForumBySlug(slug)
	if err != nil {
		return models.Users{}, models.NewError(404, models.NotFoundError)
	}
	return u.userRepo.GetByForum(foundedForum, query)
}

func (u forumUsecase) GetForumThreads(slug string, query models.PostsRequestQuery) (models.Threads, *models.Error) {
	foundedForum, err := u.forumRepo.GetForumBySlug(slug)
	if err != nil {
		return models.Threads{}, err
	}
	return u.threadRepo.GetThreadsBySlug(foundedForum, query)
}
