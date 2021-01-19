package usecase

import (
	. "github.com/yletamitlu/tech-db/internal/consts"
	"github.com/yletamitlu/tech-db/internal/forum"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/threads"
	"github.com/yletamitlu/tech-db/internal/user"
	"github.com/yletamitlu/tech-db/internal/vote"
	"strconv"
)

type ThreadUcase struct {
	threadRepos threads.ThreadRepository
	userUcase user.UserUsecase
	forumUcase forum.ForumUsecase
	voteUcase vote.VoteUsecase
}

func NewThreadUcase(
	repos threads.ThreadRepository,
	userUcase user.UserUsecase,
	forumUcase forum.ForumUsecase,
	voteUcase vote.VoteUsecase,
	) threads.ThreadUsecase {
	return &ThreadUcase{
		threadRepos: repos,
		userUcase: userUcase,
		forumUcase: forumUcase,
		voteUcase: voteUcase,
	}
}

func (tUc *ThreadUcase) Create(thread *models.Thread) (*models.Thread, error) {
	foundNickname, _ := tUc.userUcase.GetUserNickname(thread.AuthorNickname)

	if foundNickname == "" {
		return nil, ErrNotFound
	}

	forumSlug, exists := tUc.forumUcase.Exists(thread.ForumSlug)

	if !exists {
		return nil, ErrNotFound
	}

	if thread.ForumSlug != forumSlug {
		thread.ForumSlug = forumSlug
	}

	if thread.Slug != "" {
		found, err := tUc.GetBySlug(thread.Slug)

		if err != nil && err != ErrNotFound {
			return nil, err
		}

		if found != nil {
			return found, ErrAlreadyExists
		}
	}

	if err := tUc.threadRepos.InsertInto(thread); err != nil {
		return nil, err
	}

	return nil, nil
}

func (tUc *ThreadUcase) CreateVote(vote *models.Vote, slugOrId string) (*models.Thread, error) {
	id, err := strconv.Atoi(slugOrId)

	var foundThr *models.Thread
	if err == nil {
		foundThr, _ = tUc.GetById(id)
		if foundThr == nil {
			return nil, ErrNotFound
		}
		vote.Thread = id
	} else {
		foundThr, _ = tUc.GetBySlug(slugOrId)
		if foundThr == nil {
			return nil, ErrNotFound
		}
		vote.Thread = foundThr.Id
	}

	foundNickname, _ := tUc.userUcase.GetUserNickname(vote.AuthorNickname)

	if foundNickname == "" {
		return nil, ErrNotFound
	}

	result, err := tUc.voteUcase.Create(vote)

	if err != nil {
		return nil, err
	}

	foundThr.Votes += result

	tUc.threadRepos.UpdateVotes(foundThr)

	return foundThr, nil
}

func (tUc *ThreadUcase) GetThread(slugOrId string) (*models.Thread, error) {
	foundThr := &models.Thread{}

	id, err := strconv.Atoi(slugOrId)
	if err == nil {
		foundThr, err = tUc.GetById(id)
	} else {
		foundThr, err = tUc.GetBySlug(slugOrId)
	}

	if foundThr == nil {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return foundThr, nil
}

func (tUc *ThreadUcase) GetBySlug(slug string) (*models.Thread, error) {
	found, err := tUc.threadRepos.SelectBySlug(slug)

	if found == nil {
		return nil, err
	}

	return found, nil
}

func (tUc *ThreadUcase) GetById(id int) (*models.Thread, error) {
	found, err := tUc.threadRepos.SelectById(id)

	if found == nil {
		return nil, err
	}

	return found, nil
}

func (tUc *ThreadUcase) GetByForumSlug(forumSlug string, limit int, desc bool, since string) ([]*models.Thread, error) {
	_, forumExists := tUc.forumUcase.Exists(forumSlug)
	if !forumExists {
		return nil, ErrNotFound
	}

	found, err := tUc.threadRepos.SelectByForumSlug(forumSlug, limit, desc, since)

	if err != nil {
		return nil, err
	}

	if found == nil && forumExists {
		found = []*models.Thread{}
	}

	if found == nil {
		return nil, err
	}

	return found, nil
}

func (tUc *ThreadUcase) Update(updatedThread *models.Thread, slugOrId string) (*models.Thread, error) {
	id, err := strconv.Atoi(slugOrId)

	var foundThr *models.Thread
	if err == nil {
		foundThr, _ = tUc.GetById(id)
	} else {
		foundThr, _ = tUc.GetBySlug(slugOrId)
	}

	if foundThr == nil {
		return nil, ErrNotFound
	}

	if updatedThread.Message == "" && updatedThread.Title == "" {
		return foundThr, nil
	}

	updatedThread.Id = foundThr.Id

	err = tUc.threadRepos.Update(updatedThread)

	if err != nil {
		return nil, err
	}

	if updatedThread.Title != "" {
		foundThr.Title = updatedThread.Title
	}

	if updatedThread.Message != "" {
		foundThr.Message = updatedThread.Message
	}

	return foundThr, nil
}

func (tUc *ThreadUcase) GetExactFields(fields string, slugOrId string) (*models.Thread, error) {
	if id, err := strconv.Atoi(slugOrId); err == nil {
		return tUc.threadRepos.SelectThreadFields(fields, "id=$1", id)
	}

	return tUc.threadRepos.SelectThreadFields(fields, "slug=$1", slugOrId)
}
