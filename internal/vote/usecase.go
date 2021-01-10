package vote

import (
	"github.com/yletamitlu/tech-db/internal/models"
)

type VoteUsecase interface {
	Create(vote *models.Vote) (int, error)
}
