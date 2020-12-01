package usecase

import (
	. "github.com/yletamitlu/tech-db/internal/consts"
	"github.com/yletamitlu/tech-db/internal/forum"
	"github.com/yletamitlu/tech-db/internal/models"
)

type ForumUcase struct {
	forumRepos forum.ForumRepository
}

func NewUserUcase(repos forum.ForumRepository) forum.ForumUsecase {
	return &ForumUcase{
		forumRepos: repos,
	}
}

func (uUc *ForumUcase) Create(forum *models.Forum) (error, *models.Forum) {
	found, _ := uUc.forumRepos.SelectBySlug(forum.Slug)

	if found != nil {
		return ErrAlreadyExists, found
	}

	if err := uUc.forumRepos.InsertInto(forum); err != nil {
		return err, nil
	}

	return nil, nil
}
