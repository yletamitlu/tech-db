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
	found, _ := fUc.forumRepos.SelectBySlug(forum.Slug)

	if found != nil {
		return found, ErrAlreadyExists
	}

	foundUser, _ := fUc.userUcase.GetByNickname(forum.AuthorNickname)

	if foundUser == nil {
		return nil, ErrNotFound
	}

	forum.AuthorNickname = foundUser.Nickname

	if err := fUc.forumRepos.InsertInto(forum); err != nil {
		return nil, err
	}

	return forum, nil
}

func (fUc *ForumUcase) GetBySlug(slug string) (*models.Forum, error) {
	found, err := fUc.forumRepos.SelectBySlug(slug)

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

func (fUc *ForumUcase) UpdatePostsCount(delta int, forumSlug string) error {
	foundForum, _ := fUc.forumRepos.SelectBySlug(forumSlug)

	if foundForum == nil {
		return ErrNotFound
	}

	err := fUc.forumRepos.UpdatePostsCount(forumSlug, foundForum.Posts + delta)

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
