package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/yletamitlu/tech-db/internal/additions"
)

type AdditionsPgRepos struct {
	conn *sqlx.DB
}

func NewAdditionRepository(conn *sqlx.DB) additions.AdditionRepository {
	return &AdditionsPgRepos{
		conn: conn,
	}
}

func (ar *AdditionsPgRepos) Status() (uint64, uint64, uint64, uint64, error) {
	var forumsStatus, usersStatus, threadsStatus, postsStatus uint64

	statusQuery := "SELECT " +
		"(SELECT COUNT(*) FROM forums) as forums_status, " +
		"(SELECT COUNT(*) FROM threads) as threads_status, " +
		"(SELECT COUNT(*) FROM posts) as posts_status, " +
		"(SELECT COUNT(*) FROM users) as users_status"

	err := ar.conn.QueryRow(statusQuery).Scan(&forumsStatus, &threadsStatus, &postsStatus, &usersStatus)

	return forumsStatus, threadsStatus, postsStatus, usersStatus, err
}

func (ar *AdditionsPgRepos) Clear() error {
	clearQuery := "TRUNCATE users, forums, user_forum, votes, posts, threads RESTART IDENTITY"
	_, err := ar.conn.Exec(clearQuery)

	if err != nil {
		return err
	}

	return nil
}
