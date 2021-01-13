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
	foundUser, _ := vUc.userRepo.SelectByNickname(vote.AuthorNickname)

	found, _ := vUc.voteRepos.SelectByThreadAndUser(vote)

	if foundUser == nil || found == nil {
		return 0, consts.ErrNotFound
	}

	if err := vUc.voteRepos.Update(vote); err != nil {
		return 0, err
	}

	if vote.Voice == 1 {
		return 1, nil
	}

	if err := vUc.voteRepos.InsertInto(vote); err != nil {
		return 0, err
	}

	if vote.Voice == 1 {
		return 1, nil
	}

	return -1, nil
}
