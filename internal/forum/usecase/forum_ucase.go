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

func (uUc *ForumUcase) Create(forum *models.Forum) (*models.Forum, error) {
	found, _ := uUc.forumRepos.SelectBySlug(forum.Slug)

	if found != nil {
		return found, ErrAlreadyExists
	}

	foundUser, _ := uUc.userUcase.GetByNickname(forum.AuthorNickname)

	if foundUser == nil {
		return nil, ErrNotFound
	}

	forum.AuthorNickname = foundUser.Nickname

	if err := uUc.forumRepos.InsertInto(forum); err != nil {
		return nil, err
	}

	return forum, nil
}

func (uUc *ForumUcase) GetBySlug(slug string) (*models.Forum, error) {
	found, err := uUc.forumRepos.SelectBySlug(slug)

	if err != nil {
		return nil, err
	}

	return found, nil
}

func (uUc *ForumUcase) GetUsers(forumSlug string, limit int, desc bool, since string) ([]*models.User, error) {
	foundForum, _ := uUc.forumRepos.SelectBySlug(forumSlug)
	if foundForum == nil {
		return nil, ErrNotFound
	}

	found, _ := uUc.forumRepos.SelectUsers(forumSlug, limit, desc, since)

	if found == nil {
		found = []*models.User{}
	}

	return found, nil
}
