package forum

import "github.com/yletamitlu/tech-db/internal/models"

type ForumUsecase interface {
	Create(forum *models.Forum) (error, *models.Forum)
}