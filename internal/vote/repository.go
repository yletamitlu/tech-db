package vote

import (
	"github.com/yletamitlu/tech-db/internal/models"
)

type VoteRepository interface {
	InsertInto(vote *models.Vote) error
	SelectByThreadAndUser(vote *models.Vote) (*models.Vote, error)
	Update(updatedVote *models.Vote) error
}
