package post

import (
	"github.com/yletamitlu/tech-db/internal/models"
)

type PostRepository interface {
	SelectById(id int) (*models.Post, error)
	SelectByForumSlug(slug string) ([]*models.Post, error)
	//SelectByThreadId(id int) ([]*models.Post, error)
	InsertInto(post *models.Post) (*models.Post, error)
	Update(updatedPost *models.Post)
}

