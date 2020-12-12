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

func (ur *ForumPgRepos) SelectBySlug(slug string) (*models.Forum, error) {
	f := &models.Forum{}

	if err := ur.conn.Get(f,
		`SELECT * from forums where slug = $1`,
		slug);
	err != nil {
		return nil, PgxErrToCustom(err)
	}

	return f, nil
}

func (ur *ForumPgRepos) InsertInto(forum *models.Forum) error {
	if _, err := ur.conn.Exec(
		`INSERT INTO forums(author_nickname, slug, title) VALUES ($1, $2, $3)`,
		forum.AuthorNickname,
		forum.Slug,
		forum.Title);
	err != nil {
		return err
	}

	return nil
}
