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

func (fr *ForumPgRepos) SelectBySlug(slug string, withPosts bool) (*models.Forum, error) {
	f := &models.Forum{}

	if err := fr.conn.Get(f,
		`SELECT * from forums where slug = $1`,
		slug);
	err != nil {
		return nil, PgxErrToCustom(err)
	}

	//if withPosts && f.Posts == 0 {
	//	if err := fr.conn.Get(&f.Posts, `SELECT posts FROM forums WHERE slug = $1`, f.Slug); err != nil {
	//		return nil, err
	//	}
	//
	//	if err := fr.UpdatePostsCount(f.Slug, f.Posts); err != nil {
	//		return nil, err
	//	}
	//}

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

	queryString := "SELECT U.nickname, U.fullname, U.about, U.email " +
		"FROM users U JOIN user_forum UF on " +
		"U.nickname = UF.user_nickname WHERE UF.forum_slug = $1"
	var values []interface{}
	values = append(values, slug)

	query, val := MakeQueryForUsers(values, queryString, limit, desc, since)

	if err := fr.conn.Select(&users, query, val...); err != nil {
		return nil, PgxErrToCustom(err)
	}

	return users, nil
}

func (fr *ForumPgRepos) UpdatePostsCount(forumSlug string, postsCount int) error {
	_, err := fr.conn.Exec(`UPDATE forums SET posts = posts + $1 where slug = $2`,
		postsCount, forumSlug)

	if err != nil {
		return err
	}

	return nil
}

func (fr *ForumPgRepos) SelectForumSlug(slug string) (string, error) {
	var selectedSlug string
	if err := fr.conn.Get(&selectedSlug, `SELECT slug from forums where slug = $1`, slug); err != nil {
		return "", PgxErrToCustom(err)
	}

	return selectedSlug, nil
}
