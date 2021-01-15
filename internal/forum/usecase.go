package forum

import "github.com/yletamitlu/tech-db/internal/models"

type ForumUsecase interface {
	Create(forum *models.Forum) (*models.Forum, error)

	GetBySlug(slug string, withPosts bool) (*models.Forum, error)
	GetUsers(forumSlug string, limit int, desc bool, since string) ([]*models.User, error)

	Exists(forumSlug string) (string, bool)

	UpdatePostsCount(forumSlug string, count int) error
}
