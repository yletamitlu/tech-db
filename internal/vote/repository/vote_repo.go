package repository

import (
	"github.com/jmoiron/sqlx"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/vote"
)

type VotePgRepos struct {
	conn *sqlx.DB
}

func NewVoteRepository(conn *sqlx.DB) vote.VoteRepository {
	return &VotePgRepos{
		conn: conn,
	}
}

func (vr *VotePgRepos) InsertInto(vote *models.Vote) error {
	if _, err := vr.conn.Exec(
		`INSERT INTO votes(thread_id, user_nickname, voice) VALUES ($1, $2, $3)`,
		vote.Thread,
		vote.AuthorNickname,
		vote.Voice); err != nil {
		return err
	}

	return nil
}

func (vr *VotePgRepos) SelectByThreadAndUser(vote *models.Vote) (*models.Vote, error) {
	v := &models.Vote{}

	if err := vr.conn.Get(v, `SELECT * from votes where user_nickname = $1 and thread_id = $2`,
		vote.AuthorNickname, vote.Thread); err != nil {
		return nil, PgxErrToCustom(err)
	}

	return v, nil
}

func (pr *VotePgRepos) Update(updatedVote *models.Vote) error {
	_, err := pr.conn.Exec(`UPDATE votes SET voice = $1 where user_nickname = $2 and thread_id = $3`,
		updatedVote.Voice,
		updatedVote.AuthorNickname,
		updatedVote.Thread)

	if err != nil {
		return err
	}

	return nil
}