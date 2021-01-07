package post

import (
	"github.com/yletamitlu/tech-db/internal/models"
)

type PostRepository interface {
	SelectById(id int) (*models.Post, error)
	SelectByForumSlug(slug string) ([]*models.Post, error)
	InsertInto(post *models.Post) (*models.Post, error)
	InsertManyInto(posts []*models.Post) ([]*models.Post, error)
	Update(updatedPost *models.Post)
	//SelectPostsFlat(id int, limit int, desc bool, since string) ([]*models.Post, error)
	//SelectPostsTree(id int, limit int, desc bool, since string) ([]*models.Post, error)
	//SelectPostsParentTree(id int, limit int, desc bool, since string) ([]*models.Post, error)
}

