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

func (uUu *UserUcase) GetByNickname(nickname string) (*models.User, error) {
	u, err := uUu.userRepos.SelectByNickname(nickname)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (uUu *UserUcase) Update(updatedUser *models.User) error {
	if err := uUu.userRepos.Update(updatedUser); err != nil {
		return err
	}

	return nil
}
