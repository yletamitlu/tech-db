package repository

import (
	"github.com/jmoiron/sqlx"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/threads"
)

type ThreadPgRepos struct {
	conn *sqlx.DB
}

func NewThreadRepository(conn *sqlx.DB) threads.ThreadRepository {
	return &ThreadPgRepos{
		conn: conn,
	}
}

func (tr *ThreadPgRepos) SelectBySlug(slug string) (*models.Thread, error) {
	thread := &models.Thread{}

	if err := tr.conn.Get(thread,
		`SELECT * from threads where slug = $1`,
		slug); err != nil {
		return nil, PgxErrToCustom(err)
	}

	if thread.Slug == "" {
		return nil, nil
	}

	return thread, nil
}

func (tr *ThreadPgRepos) SelectById(id int) (*models.Thread, error) {
	thread := &models.Thread{}

	if err := tr.conn.Get(thread,
		`SELECT * from threads where id = $1`,
		id); err != nil {
		return nil, PgxErrToCustom(err)
	}

	return thread, nil
}

func (tr *ThreadPgRepos) SelectByForumSlug(slug string, limit int, desc bool, since string) ([]*models.Thread, error) {
	var threads []*models.Thread

	if since == "" {
		if desc {
			if err := tr.conn.Select(&threads,
				`SELECT * from threads where forum_slug = $1 order by created_at desc limit $2`,
				slug, limit); err != nil {
				return nil, PgxErrToCustom(err)
			}
		} else {
			if err := tr.conn.Select(&threads,
				`SELECT * from threads where forum_slug = $1 order by created_at limit $2`,
				slug, limit); err != nil {
				return nil, PgxErrToCustom(err)
			}
		}
	} else {
		if desc {
			if err := tr.conn.Select(&threads,
				`SELECT * from threads where forum_slug = $1 and created_at <= $2 order by created_at desc limit $3`,
				slug, since, limit); err != nil {
				return nil, PgxErrToCustom(err)
			}
		} else {
			if err := tr.conn.Select(&threads,
				`SELECT * from threads where forum_slug = $1 and created_at >= $2 order by created_at limit $3`,
				slug, since, limit); err != nil {
				return nil, PgxErrToCustom(err)
			}
		}
	}

	return threads, nil
}

func (tr *ThreadPgRepos) InsertInto(thread *models.Thread) error {
	if err := tr.conn.QueryRow(
		`INSERT INTO threads (author_nickname, forum_slug, message, title, created_at, slug) 
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		thread.AuthorNickname,
		thread.ForumSlug,
		thread.Message,
		thread.Title,
		thread.Created,
		thread.Slug).Scan(&thread.Id);
	err != nil {
		return err
	}

	return nil
}

func (tr *ThreadPgRepos) Update(updatedThread *models.Thread) {
	_, _ = tr.conn.Exec(`UPDATE threads SET slug = $1, author_nickname = $2, title = $3,
                   message = $4 where id = $5`,
		updatedThread.Slug,
		updatedThread.AuthorNickname,
		updatedThread.Title,
		updatedThread.Message,
		updatedThread.Id)
}

func (tr *ThreadPgRepos) UpdateVotes(updatedThread *models.Thread) {
	_, _ = tr.conn.Exec(`UPDATE threads SET votes = $1 where id = $2`,
		updatedThread.Votes, updatedThread.Id)
}
