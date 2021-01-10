package threads

import (
	"github.com/yletamitlu/tech-db/internal/models"
)

type ThreadUsecase interface {
	Create(thread *models.Thread) (*models.Thread, error)
	GetBySlug(slug string) (*models.Thread, error)
	GetById(id int) (*models.Thread, error)
	GetByForumSlug(forumSlug string, limit int, desc bool, since string) ([]*models.Thread, error)
	CreateVote(vote *models.Vote, slugOrId string) (*models.Thread, error)
	Update(updatedThread *models.Thread) (*models.Thread, error)
}
