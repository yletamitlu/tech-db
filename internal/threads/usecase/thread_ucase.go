package usecase

import (
	. "github.com/yletamitlu/tech-db/internal/consts"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/threads"
	"github.com/yletamitlu/tech-db/internal/user"
	"time"
)

type ThreadUcase struct {
	threadRepos threads.ThreadRepository
	userUcase user.UserUsecase
}

func NewThreadUcase(repos threads.ThreadRepository, userUcase user.UserUsecase) threads.ThreadUsecase {
	return &ThreadUcase{
		threadRepos: repos,
		userUcase: userUcase,
	}
}

func (tUc *ThreadUcase) Create(thread *models.Thread) (*models.Thread, error) {
	found, _ := tUc.GetBySlug(thread.Slug)

	if found != nil {
		return found, ErrAlreadyExists
	}

	if err := tUc.threadRepos.InsertInto(thread); err != nil {
		return nil, err
	}

	return nil, nil
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

func (tUc *ThreadUcase) GetByForumSlug(forumSlug string, limit int, desc bool, since time.Time) ([]*models.Thread, error) {
	found, err := tUc.threadRepos.SelectByForumSlug(forumSlug, limit, desc, since)

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

func (tUc *ThreadUcase) GetByDate(forumSlug string, date time.Time) (*models.Thread, error) {
	found, err := tUc.threadRepos.SelectByDate(forumSlug, date)

	if found == nil {
		return nil, err
	}

	return found, nil
}