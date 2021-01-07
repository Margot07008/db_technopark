package repository

import (
	"db_technopark/application/models"
	"db_technopark/application/vote"
	"github.com/jackc/pgx"
)

type pgVoteRepository struct {
	conn *pgx.ConnPool
}

func NewPgVoteRepository(db *pgx.ConnPool) vote.Repository {
	return &pgVoteRepository{conn: db}
}

func (p pgVoteRepository) Create(vote models.Vote) *models.Error {
	resInsert, err := p.conn.Exec(`insert into main.votes (nickname, voice, thread) values ($1, $2, $3)`,
		vote.Nickname, vote.Voice, vote.Thread)
	if err != nil {
		return models.NewError(404, models.InternalError)
	}
	if resInsert.RowsAffected() == 0 {
		return models.NewError(404, models.NotFoundError)
	}
	return nil
}

func (p pgVoteRepository) Update(vote models.Vote) *models.Error {
	res, err := p.conn.Exec(`update main.votes set voice = $1 where lower(nickname) = lower($2) and thread = $3`,
		vote.Voice, vote.Nickname, vote.Thread)
	if err != nil {
		return models.NewError(500, models.UpdateError)
	}
	if res.RowsAffected() == 0 {
		return models.NewError(404, models.NotFoundError)
	}
	return nil
}
