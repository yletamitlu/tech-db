package user

import "github.com/yletamitlu/tech-db/internal/models"

type UserRepository interface {
	SelectByNicknameOrEmail(nickname string, email string) ([]*models.User, error)
	InsertInto(user *models.User) error
	SelectByNickname(nickname string) (*models.User, error)
	Update(updatedUser *models.User) error
}
