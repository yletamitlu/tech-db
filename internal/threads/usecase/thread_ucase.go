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
	found, _ := tUc.GetBySlug(thread.Slug)

	if found != nil {
		return found, ErrAlreadyExists
	}

	foundAuthor, _ := tUc.userUcase.GetByNickname(thread.AuthorNickname)

	if foundAuthor == nil {
		return nil, ErrNotFound
	}

	foundForum, _ := tUc.forumUcase.GetBySlug(thread.ForumSlug)

	if foundForum != nil && thread.ForumSlug != foundForum.Slug {
		thread.ForumSlug = foundForum.Slug
	}

	if foundForum == nil {
		return nil, ErrNotFound
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

	foundAuthor, _ := tUc.userUcase.GetByNickname(vote.AuthorNickname)

	if foundAuthor == nil {
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
	foundForum, err := tUc.forumUcase.GetBySlug(forumSlug)

	if err != nil {
		return nil, err
	}

	found, err := tUc.threadRepos.SelectByForumSlug(forumSlug, limit, desc, since)

	if err != nil {
		return nil, err
	}

	if found == nil && foundForum != nil {
		found = []*models.Thread{}
	}

	if found == nil {
		return nil, err
	}

	return found, nil
}

func (tUc *ThreadUcase) Update(updatedThread *models.Thread) (*models.Thread, error) {
	u, _ := tUc.threadRepos.SelectById(updatedThread.Id)

	if u != nil {
		return nil, ErrConflict
	}

	tUc.threadRepos.Update(updatedThread)

	//u, err := tUc.threadRepos.SelectByNickname(updatedThread.Nickname)
	//if u == nil {
	//	return nil, ErrNotFound
	//}
	//if err != nil {
	//	return nil, err
	//}

	return u, nil
}
