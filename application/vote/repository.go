package vote

import "db_technopark/application/models"

type Repository interface {
	Create(vote models.Vote) *models.Error
	Update(vote models.Vote) *models.Error
	//GetByNicknameAndThreadID(nickname string, threadID int32) (models.Vote, *models.Error)
}
