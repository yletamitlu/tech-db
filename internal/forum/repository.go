package forum

import "github.com/yletamitlu/tech-db/internal/models"

type ForumRepository interface {
	SelectBySlug(slug string) (*models.Forum, error)
	InsertInto(forum *models.Forum) error
	SelectUsers(slug string, limit int, desc bool, since string) ([]*models.User, error)
	UpdatePostsCount(forumSlug string, postsCount int) error

	SelectForumSlug(slug string) (string, error)
}
