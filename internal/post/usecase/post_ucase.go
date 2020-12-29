package usecase

import (
	. "github.com/yletamitlu/tech-db/internal/consts"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/post"
	"github.com/yletamitlu/tech-db/internal/threads"
	"github.com/yletamitlu/tech-db/internal/user"
	"strconv"
)

type PostUcase struct {
	postRepos post.PostRepository
	userUcase user.UserUsecase
	threadUcase threads.ThreadUsecase
}

func NewPostUcase(repos post.PostRepository, userUcase user.UserUsecase, threadUcase threads.ThreadUsecase) *PostUcase {
	return &PostUcase{
		postRepos: repos,
		userUcase:   userUcase,
		threadUcase:  threadUcase,
	}
}

func (pUc *PostUcase) Create(post *models.Post, thread string) (*models.Post, error) {
	id, err := strconv.Atoi(thread)
	if err == nil {
		foundThr, _ := pUc.threadUcase.GetById(id)
		if foundThr == nil {
			return nil, ErrNotFound
		}
		post.ForumSlug = foundThr.ForumSlug
		post.Thread = id
	} else {
		foundThr, _ := pUc.threadUcase.GetBySlug(thread)
		if foundThr == nil {
			return nil, ErrNotFound
		}
		post.ForumSlug = foundThr.ForumSlug
		post.Thread = foundThr.Id
	}

	found, _ := pUc.postRepos.SelectById(post.Id)

	if found != nil {
		return found, ErrAlreadyExists
	}

	resultPost, err := pUc.postRepos.InsertInto(post)
	if  err != nil {
		return nil, err
	}

	return resultPost, nil
}

func (pUc *PostUcase) GetById(id int) (*models.Post, error) {
	found, _ := pUc.postRepos.SelectById(id)

	if found == nil {
		return nil, ErrNotFound
	}

	return found, nil
}

func (pUc *PostUcase) GetByForumSlug(slug string) ([]*models.Post, error) {
	found, _ := pUc.postRepos.SelectByForumSlug(slug)

	if found == nil {
		return nil, ErrNotFound
	}

	return found, nil
}

func (pUc *PostUcase) GetByThreadId(id int) ([]*models.Post, error) {
	found, _ := pUc.GetByThreadId(id)

	if found == nil {
		return nil, ErrNotFound
	}

	return found, nil
}

func (pUc *PostUcase) Update(updatedPost *models.Post) (*models.Post, error) {
	pUc.postRepos.Update(updatedPost)

	u, _ := pUc.postRepos.SelectById(updatedPost.Id)

	return u, nil
}
