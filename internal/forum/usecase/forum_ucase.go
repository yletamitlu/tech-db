package usecase

import (
	. "github.com/yletamitlu/tech-db/internal/consts"
	"github.com/yletamitlu/tech-db/internal/forum"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/user"
)

type ForumUcase struct {
	forumRepos forum.ForumRepository
	userUcase user.UserUsecase
}

func NewForumUcase(repos forum.ForumRepository, uUcase user.UserUsecase) forum.ForumUsecase {
	return &ForumUcase{
		forumRepos: repos,
		userUcase: uUcase,
	}
}

func (fUc *ForumUcase) Create(forum *models.Forum) (*models.Forum, error) {
	foundNickname, _ := fUc.userUcase.GetUserNickname(forum.AuthorNickname)

	if foundNickname == "" {
		return nil, ErrNotFound
	}

	found, _ := fUc.forumRepos.SelectBySlug(forum.Slug, false)

	if found != nil {
		return found, ErrAlreadyExists
	}

	forum.AuthorNickname = foundNickname

	if err := fUc.forumRepos.InsertInto(forum); err != nil {
		return nil, err
	}

	return forum, nil
}

func (fUc *ForumUcase) GetBySlug(slug string, withPosts bool) (*models.Forum, error) {
	found, err := fUc.forumRepos.SelectBySlug(slug, withPosts)

	if err != nil {
		return nil, err
	}

	return found, nil
}

func (fUc *ForumUcase) GetUsers(forumSlug string, limit int, desc bool, since string) ([]*models.User, error) {
	if _, exists := fUc.Exists(forumSlug); !exists {
		return nil, ErrNotFound
	}

	found, _ := fUc.forumRepos.SelectUsers(forumSlug, limit, desc, since)

	if found == nil {
		found = []*models.User{}
	}

	return found, nil
}

func (fUc *ForumUcase) UpdatePostsCount(forumSlug string, count int) error {
	err := fUc.forumRepos.UpdatePostsCount(forumSlug, count)

	if err != nil {
		return err
	}

	return nil
}

func (fUc *ForumUcase) Exists(forumSlug string) (string, bool) {
	slug, err := fUc.forumRepos.SelectForumSlug(forumSlug)
	if  err != nil {
		return "", false
	}

	return slug, true
}
