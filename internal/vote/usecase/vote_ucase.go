package usecase

import (
	"github.com/yletamitlu/tech-db/internal/consts"
	_ "github.com/yletamitlu/tech-db/internal/consts"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/user"
	"github.com/yletamitlu/tech-db/internal/vote"
)

type VoteUcase struct {
	voteRepos vote.VoteRepository
	userRepo user.UserRepository
}

func NewVoteUcase(repos vote.VoteRepository, usRepo user.UserRepository) vote.VoteUsecase {
	return &VoteUcase{
		voteRepos: repos,
		userRepo: usRepo,
	}
}

func (vUc *VoteUcase) Create(vote *models.Vote) (int, error) {
	foundNickname, _ := vUc.userRepo.SelectUserNickname(vote.AuthorNickname)

	if foundNickname == "" {
		return 0, consts.ErrNotFound
	}

	found, _ := vUc.voteRepos.SelectByThreadAndUser(vote)

	if found != nil {
		newVoice := vote.Voice - found.Voice

		if found.Voice == vote.Voice {
			return 0, nil
		}

		if err := vUc.voteRepos.Update(vote); err != nil {
			return 0, err
		}

		return newVoice, nil
	}

	if err := vUc.voteRepos.InsertInto(vote); err != nil {
		return 0, err
	}

	return vote.Voice, nil
}
