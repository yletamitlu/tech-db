package usecase

import (
	. "github.com/yletamitlu/tech-db/internal/consts"
	"github.com/yletamitlu/tech-db/internal/forum"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/post"
	"github.com/yletamitlu/tech-db/internal/threads"
	"github.com/yletamitlu/tech-db/internal/user"
	"strconv"
	"time"
)

type PostUcase struct {
	postRepos        post.PostRepository
	userUcase        user.UserUsecase
	threadUcase      threads.ThreadUsecase
	forumUcase       forum.ForumUsecase
}

func NewPostUcase(repos post.PostRepository, userUcase user.UserUsecase,
	threadUcase threads.ThreadUsecase, forumUcase forum.ForumUsecase) *PostUcase {
	return &PostUcase{
		postRepos:        repos,
		userUcase:        userUcase,
		threadUcase:      threadUcase,
		forumUcase:       forumUcase,
	}
}

func (pUc *PostUcase) GetPostAuthor(nickname string) (*models.User, error) {
	foundUser, err := pUc.userUcase.GetByNickname(nickname)

	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (pUc *PostUcase) GetPostThread(threadId int) (*models.Thread, error) {
	foundThread, err := pUc.threadUcase.GetById(threadId)

	if err != nil {
		return nil, err
	}

	return foundThread, nil
}

func (pUc *PostUcase) GetPostForum(forumSlug string) (*models.Forum, error) {
	foundForum, err := pUc.forumUcase.GetBySlug(forumSlug)

	if err != nil {
		return nil, err
	}

	return foundForum, nil
}

func (pUc *PostUcase) Create(posts []*models.Post, slugOrId string) ([]*models.Post, error) {
	foundThr := &models.Thread{}

	id, err := strconv.Atoi(slugOrId)
	if err == nil {
		foundThr, _ = pUc.threadUcase.GetById(id)
		if foundThr == nil {
			return nil, ErrNotFound
		}
	} else {
		foundThr, _ = pUc.threadUcase.GetBySlug(slugOrId)
		if foundThr == nil {
			return nil, ErrNotFound
		}
	}

	if len(posts) == 0 {
		return nil, nil
	}

	createdAt := time.Now().Format(time.RFC3339)

	for _, pst := range posts {
		foundAuthor, err := pUc.userUcase.GetByNickname(pst.AuthorNickname)

		if foundAuthor == nil {
			return nil, ErrNotFound
		}

		if err != nil {
			return nil, err
		}

		pst.Created = createdAt
		pst.ForumSlug = foundThr.ForumSlug
		pst.Thread = foundThr.Id

		if pst.Parent != 0 {
			foundParent, err := pUc.postRepos.SelectById(pst.Parent)

			if (foundParent != nil && foundParent.Thread != pst.Thread) || err != nil {
				return nil, ErrConflict
			}
		}
	}

	resultPosts, err := pUc.postRepos.InsertManyInto(posts)

	if err != nil {
		return nil, err
	}

	err = pUc.forumUcase.UpdatePostsCount(len(posts), foundThr.ForumSlug)

	if err != nil {
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
	found, _ := pUc.postRepos.SelectById(updatedPost.Id)

	if found == nil {
		return nil, ErrNotFound
	}

	if updatedPost.Message == "" || updatedPost.Message == found.Message {
		return found, nil
	}

	err := pUc.postRepos.Update(updatedPost)

	if err != nil {
		return nil, err
	}

	found.Message = updatedPost.Message
	found.IsEdited = true

	return found, nil
}

func (pUc *PostUcase) GetPosts(slugOrId string, limit int, desc bool, since string, sort string) ([]*models.Post, error) {
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		foundThr, _ := pUc.threadUcase.GetBySlug(slugOrId)
		if foundThr == nil {
			return nil, ErrNotFound
		}

		id = foundThr.Id
	} else {
		foundThr, _ := pUc.threadUcase.GetById(id)

		if foundThr == nil {
			return nil, ErrNotFound
		}

		id = foundThr.Id
	}

	var posts []*models.Post

	switch sort {
	case "tree":
		posts, err = pUc.postRepos.SelectPostsTree(id, limit, desc, since)
	case "parent_tree":
		posts, err = pUc.postRepos.SelectPostsParentTree(id, limit, desc, since)
	default:
		posts, err = pUc.postRepos.SelectPostsFlat(id, limit, desc, since)
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}
