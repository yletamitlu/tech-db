package threads

import (
	"github.com/yletamitlu/tech-db/internal/models"
	"time"
)

type ThreadUsecase interface {
	Create(thread *models.Thread) (*models.Thread, error)
	GetBySlug(slug string) (*models.Thread, error)
	GetById(id int) (*models.Thread, error)
	GetByForumSlug(forumSlug string, limit int, desc bool, since time.Time) ([]*models.Thread, error)
	GetByDate(forumSlug string, date time.Time) (*models.Thread, error)
	Update(updatedThread *models.Thread) (*models.Thread, error)
}
