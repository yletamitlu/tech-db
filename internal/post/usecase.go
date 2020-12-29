package post

import "github.com/yletamitlu/tech-db/internal/models"

type PostUsecase interface {
	Create(post *models.Post, thread string) (*models.Post, error)
	GetById(id int) (*models.Post, error)
	GetByForumSlug(slug string) ([]*models.Post, error)
	GetByThreadId(id int) ([]*models.Post, error)
	Update(updatedPost *models.Post) (*models.Post, error)
}
