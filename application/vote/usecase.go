package vote

import "db_technopark/application/models"

type Usecase interface {
	UpsertVote(vote models.Vote) (models.Thread, *models.Error)
}
