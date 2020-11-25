package user

import "github.com/yletamitlu/tech-db/internal/models"

type UserUsecase interface {
	Create(user *models.User) (error, []*models.User) // ошибка, найденные пользователи
}
