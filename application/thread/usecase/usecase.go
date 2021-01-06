package usecase

import (
	"db_technopark/application/forum"
	"db_technopark/application/models"
	"db_technopark/application/thread"
	"db_technopark/application/user"
)

type threadUsecase struct {
	threadRepo thread.Repository
	userRepo   user.Repository
	forumRepo  forum.Repository
}

func NewThreadUsecase(threadRepo thread.Repository, userRepo user.Repository, forumRepo forum.Repository) thread.Usecase {
	return &threadUsecase{
		threadRepo: threadRepo,
		userRepo:   userRepo,
		forumRepo:  forumRepo,
	}
}

func (u threadUsecase) CreateThread(slug string, thread models.Thread) (models.Thread, *models.Error) {
	foundedForum, err := u.forumRepo.GetForumBySlug(slug)
	if err != nil {
		return models.Thread{}, err
	}
	foundedUser, err := u.userRepo.GetByNickname(thread.Author)
	if err != nil {
		return models.Thread{}, err
	}
	createdThread, err := u.threadRepo.CreateThread(foundedForum, foundedUser, thread)
	if err != nil && err.StatusCode == 409 {
		//TODO add getting thread by Slug
		return thread, models.NewError(409, models.ConflictError)
	}
	return createdThread, err
}
