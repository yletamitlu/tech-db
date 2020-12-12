package threads

import (
	"github.com/yletamitlu/tech-db/internal/models"
	"time"
)

type ThreadRepository interface {
	SelectBySlug(slug string) (*models.Thread, error)
	SelectById(id int) (*models.Thread, error)
	SelectByForumSlug(slug string, limit int, desc bool, since time.Time) ([]*models.Thread, error)
	InsertInto(thread *models.Thread) error
	Update(updatedThread *models.Thread)
	SelectByDate(forumSlug string, since time.Time) (*models.Thread, error)
}
