package usecase

import (
	. "github.com/yletamitlu/tech-db/internal/consts"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/user"
)

type UserUcase struct {
	userRepos user.UserRepository
}

func NewUserUcase(repos user.UserRepository) user.UserUsecase {
	return &UserUcase{
		userRepos: repos,
	}
}

func (uUc *UserUcase) Create(user *models.User) (error, []*models.User) {
	found, _ := uUc.userRepos.SelectByNicknameOrEmail(user.Nickname, user.Email)

	if found != nil {
		return ErrAlreadyExists, found
	}

	if err := uUc.userRepos.InsertInto(user); err != nil {
		return err, nil
	}

	return nil, nil
}
