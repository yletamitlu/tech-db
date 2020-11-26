package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/yletamitlu/tech-db/internal/consts"
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

func (ur *UserPgRepos) SelectByNickname(nickname string) (*models.User, error) {
	u := &models.User{}

	if err := ur.conn.Get(u, `SELECT * from users where nickname = $1`, nickname); err != nil {
		return nil, PgxErrToCustom(err)
	}

	return u, nil
}

func (ur *UserPgRepos) Update(updatedUser *models.User) error {
	var affected int64

	res, err := ur.conn.Exec(`UPDATE users SET email = $1, fullname = $2, about = $3 where nickname = $4`,
		updatedUser.Email, updatedUser.FullName, updatedUser.About, updatedUser.Nickname)

	if err != nil {
		return err
	}

	if res != nil {
		affected, _ = res.RowsAffected()
	}

	if affected <= 0 {
		return consts.ErrNotFound
	}

	return nil
}
