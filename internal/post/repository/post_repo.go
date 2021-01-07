package repository

import (
	"github.com/jmoiron/sqlx"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/post"
)

type PostPgRepos struct {
	conn *sqlx.DB
}

func NewPostRepository(conn *sqlx.DB) post.PostRepository {
	return &PostPgRepos{
		conn: conn,
	}
}

func (pr *PostPgRepos) SelectById(id int) (*models.Post, error) {
	p := &models.Post{}

	if err := pr.conn.Get(p,
		`SELECT * from posts where id = $1`,
		id); err != nil {
		return nil, PgxErrToCustom(err)
	}

	return p, nil
}

func (pr *PostPgRepos) SelectByForumSlug(slug string) ([]*models.Post, error) {
	var posts []*models.Post

	if err := pr.conn.Select(&posts,
		`SELECT * from posts where forum_slug = $1`, slug); err != nil {
		return nil, PgxErrToCustom(err)
	}

	return posts, nil
}

func (pr *PostPgRepos) InsertInto(post *models.Post) (*models.Post, error) {
	if err := pr.conn.QueryRow(
		`INSERT INTO posts (author_nickname, forum_slug, message, created_at, thread_id) 
				VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		post.AuthorNickname,
		post.ForumSlug,
		post.Message,
		post.Created,
		post.Thread).Scan(&post.Id);
		err != nil {
		return nil, err
	}

	return post, nil
}

func (pr *PostPgRepos) InsertManyInto(posts []*models.Post) ([]*models.Post, error) {
	stmt, err := pr.conn.Prepare(`
		INSERT INTO posts (author_nickname, forum_slug, message, created_at, thread_id) 
				VALUES ($1, $2, $3, $4, $5) RETURNING id`)
	if err != nil {
		return nil, err
	}

	for _, pst := range posts {
		err := stmt.QueryRow(pst.AuthorNickname, pst.ForumSlug,
			pst.Message, pst.Created, pst.Thread).Scan(&pst.Id)
		if err != nil {
			return nil, err
		}
	}

	return posts, nil
}

func (pr *PostPgRepos) Update(updatedPost *models.Post) {
	_, _ = pr.conn.Exec(`UPDATE threads SET author_nickname = $1,
                   message = $2 where id = $3`,
		updatedPost.AuthorNickname,
		updatedPost.Message,
		updatedPost.Id)
}

//func (pr *PostPgRepos) SelectPostsFlat(id int, limit int, desc bool, since string) ([]*models.Post, error) {
//
//}
//
//func (pr *PostPgRepos) SelectPostsTree(id int, limit int, desc bool, since string) ([]*models.Post, error) {
//
//}
//
//func (pr *PostPgRepos) SelectPostsParentTree(id int, limit int, desc bool, since string) ([]*models.Post, error) {
//
//}
