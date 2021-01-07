package usecase

import (
	"db_technopark/application/forum"
	"db_technopark/application/models"
	"db_technopark/application/thread"
	"db_technopark/application/user"
	"strconv"
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

	thread.Forum = foundedForum.Slug
	createdThread, err := u.threadRepo.CreateThread(foundedForum, foundedUser, thread)
	if err != nil && err.StatusCode == 409 {
		thread, _ = u.threadRepo.GetBySlug(thread.Slug)
		return thread, models.NewError(409, models.ConflictError)
	}
	//err := u.forumRepo.UpdateThreadField(slug)
	return createdThread, err
}

func (u threadUsecase) GetThreadBySlugOrId(data string, isSlug bool) (models.Thread, *models.Error) {
	if isSlug {
		return u.threadRepo.GetBySlug(data)
	}
	id, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return models.Thread{}, models.NewError(400, models.BadRequestError)
	}
	return u.threadRepo.GetByID(int32(id))
}

func (u threadUsecase) ThreadUpdate(threadUpdate models.ThreadUpdate) (models.Thread, *models.Error) {
	var foundThread models.Thread
	var err *models.Error
	if threadUpdate.Slug != "" {
		foundThread, err = u.threadRepo.GetBySlug(threadUpdate.Slug)

	} else {
		foundThread, err = u.threadRepo.GetByID(threadUpdate.Id)
	}

	if err != nil {
		return models.Thread{}, err
	}

	return u.threadRepo.Update(foundThread, threadUpdate)
}
