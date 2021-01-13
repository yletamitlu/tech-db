package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/yletamitlu/tech-db/internal/forum"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	"github.com/yletamitlu/tech-db/internal/models"
)

type ForumPgRepos struct {
	conn *sqlx.DB
}

func NewForumRepository(conn *sqlx.DB) forum.ForumRepository {
	return &ForumPgRepos{
		conn: conn,
	}
}

func (fr *ForumPgRepos) SelectBySlug(slug string) (*models.Forum, error) {
	f := &models.Forum{}

	if err := fr.conn.Get(f,
		`SELECT * from forums where slug = $1`,
		slug);
	err != nil {
		return nil, PgxErrToCustom(err)
	}

	return f, nil
}

func (fr *ForumPgRepos) InsertInto(forum *models.Forum) error {
	if _, err := fr.conn.Exec(
		`INSERT INTO forums(author_nickname, slug, title) VALUES ($1, $2, $3)`,
		forum.AuthorNickname,
		forum.Slug,
		forum.Title);
	err != nil {
		return err
	}

	return nil
}

func (fr *ForumPgRepos) SelectUsers(slug string, limit int, desc bool, since string) ([]*models.User, error) {
	var users []*models.User

	//

	return users, nil
}

func (fr *ForumPgRepos) UpdatePostsCount(forumSlug string, postsCount int) error {
	_, err := fr.conn.Exec(`UPDATE forums SET posts = $1 where slug = $2`,
		postsCount, forumSlug)

	if err != nil {
		return err
	}

	return nil
}
