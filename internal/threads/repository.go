package threads

import (
	"github.com/yletamitlu/tech-db/internal/models"
)

type ThreadRepository interface {
	SelectBySlug(slug string) (*models.Thread, error)
	SelectById(id int) (*models.Thread, error)
	SelectByForumSlug(slug string, limit int, desc bool, since string) ([]*models.Thread, error)
	InsertInto(thread *models.Thread) error
	Update(updatedThread *models.Thread)
}
