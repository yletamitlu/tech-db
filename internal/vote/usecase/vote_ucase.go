package usecase

import (
	_ "github.com/yletamitlu/tech-db/internal/consts"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/vote"
)

type VoteUcase struct {
	voteRepos vote.VoteRepository
}

func NewVoteUcase(repos vote.VoteRepository) vote.VoteUsecase {
	return &VoteUcase{
		voteRepos: repos,
	}
}

func (vUc *VoteUcase) Create(vote *models.Vote) (int, error) {
	found, _ := vUc.voteRepos.SelectByThreadAndUser(vote)

	if found != nil {
		if found.Voice == vote.Voice {
			//return 0, nil
		}

		if err := vUc.voteRepos.Update(vote); err != nil {
			return 0, err
		}

		if vote.Voice == 1 {
			return 1, nil
		}

		return -1, nil
	}

	if err := vUc.voteRepos.InsertInto(vote); err != nil {
		return 0, err
	}

	if vote.Voice == 1 {
		return 1, nil
	}

	return -1, nil
}
