package repository

import (
	"github.com/jmoiron/sqlx"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/user"
)

type UserPgRepos struct {
	conn *sqlx.DB
}

func NewUserRepository(conn *sqlx.DB) user.UserRepository {
	return &UserPgRepos{
		conn: conn,
	}
}

func (ur *UserPgRepos) SelectByNicknameOrEmail(nickname string, email string) ([]*models.User, error) {
	var users []*models.User

	if err := ur.conn.Select(&users,
		`SELECT * from users where nickname = $1 or email = $2`,
		nickname,
		email); err != nil {
		return nil, PgxErrToCustom(err)
	}

	return users, nil
}

func (ur *UserPgRepos) InsertInto(user *models.User) error {
	if _, err := ur.conn.Exec(
		`INSERT INTO users(nickname, fullname, email, about) VALUES ($1, $2, $3, $4)`,
		user.Nickname,
		user.FullName,
		user.Email,
		user.About); err != nil {
		return err
	}

	return nil
}
