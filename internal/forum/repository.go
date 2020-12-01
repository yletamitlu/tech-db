package forum

import "github.com/yletamitlu/tech-db/internal/models"

type ForumRepository interface {
	SelectBySlug(slug string) (*models.Forum, error)
	InsertInto(forum *models.Forum) error
}
