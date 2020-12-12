package forum

import "github.com/yletamitlu/tech-db/internal/models"

type ForumUsecase interface {
	Create(forum *models.Forum) (*models.Forum, error)
	GetBySlug(slug string) (*models.Forum, error)
}