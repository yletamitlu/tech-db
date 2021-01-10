package usecase

import (
	"fmt"
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

func (pUc *PostUcase) Create(posts []*models.Post, slugOrId string) ([]*models.Post, error) {
	if len(posts) == 0 {
		return nil, nil
	}

	var resultPosts []*models.Post

	for _, pst := range posts {
		id, err := strconv.Atoi(slugOrId)
		if err == nil {
			foundThr, _ := pUc.threadUcase.GetById(id)
			if foundThr == nil {
				return nil, ErrNotFound
			}
			pst.ForumSlug = foundThr.ForumSlug
			pst.Thread = id
		} else {
			foundThr, _ := pUc.threadUcase.GetBySlug(slugOrId)
			if foundThr == nil {
				return nil, ErrNotFound
			}
			pst.ForumSlug = foundThr.ForumSlug
			pst.Thread = foundThr.Id
		}
	}

	resultPosts, err := pUc.postRepos.InsertManyInto(posts)
	if  err != nil {
		return nil, err
	}

	return resultPosts, nil
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

func (pUc *PostUcase) GetPosts(slugOrId string, limit int, desc bool, since string, sort string) ([]*models.Post, error) {
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		foundThr, _ := pUc.threadUcase.GetBySlug(slugOrId)
		if foundThr == nil {
			return nil, ErrNotFound
		}

		id = foundThr.Id
	}

	var posts []*models.Post

	fmt.Println(id)

	switch sort {
	case "tree":
		//posts, err = pUc.postRepos.SelectPostsTree(id, limit, desc, since)
	case "parent_tree":
		//posts, err = pUc.postRepos.SelectPostsParentTree(id, limit, desc, since)
	default:
		//posts, err = pUc.postRepos.SelectPostsFlat(id, limit, desc, since)
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}
