package repository

import (
	"github.com/jmoiron/sqlx"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/user"
)

type UserPgRepos struct {
	conn      *sqlx.DB
	userCache *UserCache
}

func NewUserRepository(conn *sqlx.DB) user.UserRepository {
	return &UserPgRepos{
		conn: conn,
		userCache: NewUserCache(),
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

func (ur *UserPgRepos) SelectUserNickname(nickname string) (string, error) {
	if realNickname, err := ur.userCache.GetNickname(nickname); err == nil {
		return realNickname, nil
	}

	user := &models.User{}
	if err := ur.conn.Get(
		user,
		`SELECT id, nickname FROM users WHERE nickname = $1`,
		nickname,
	); err != nil {
		return "", PgxErrToCustom(err)
	}

	ur.userCache.Set(user.Id, user.Nickname)

	return user.Nickname, nil
}

func (ur *UserPgRepos) InsertInto(user *models.User) error {
	if err := ur.conn.QueryRow(
		`INSERT INTO users(nickname, fullname, email, about) VALUES ($1, $2, $3, $4) RETURNING id`,
		user.Nickname,
		user.FullName,
		user.Email,
		user.About).Scan(&user.Id); err != nil {
		return err
	}

	ur.userCache.Set(user.Id, user.Nickname)

	return nil
}

func (ur *UserPgRepos) SelectByNickname(nickname string) (*models.User, error) {
	u := &models.User{}

	if err := ur.conn.Get(u, `SELECT * from users where nickname = $1`, nickname); err != nil {
		return nil, PgxErrToCustom(err)
	}

	return u, nil
}

func (ur *UserPgRepos) SelectByEmail(email string) (*models.User, error) {
	u := &models.User{}

	if err := ur.conn.Get(u, `SELECT * from users where email = $1`, email); err != nil {
		return nil, PgxErrToCustom(err)
	}

	return u, nil
}

func (ur *UserPgRepos) Update(updatedUser *models.User) {

	if updatedUser.Email != "" {
		_, _ = ur.conn.Exec(`UPDATE users SET email = $1 where nickname = $2`,
			updatedUser.Email, updatedUser.Nickname)
	}

	if updatedUser.FullName != "" {
		_, _ = ur.conn.Exec(`UPDATE users SET fullname = $1 where nickname = $2`,
			updatedUser.FullName, updatedUser.Nickname)
	}

	if updatedUser.About != "" {
		_, _ = ur.conn.Exec(`UPDATE users SET about = $1 where nickname = $2`,
			updatedUser.About, updatedUser.Nickname)
	}
}
