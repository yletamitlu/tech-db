package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/post"
)

type PostPgRepos struct {
	conn *sqlx.DB
	postIdsGenerator *Generator
}

func NewPostRepository(conn *sqlx.DB) post.PostRepository {
	return &PostPgRepos{
		conn: conn,
		postIdsGenerator: NewGenerator(),
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

func (pr *PostPgRepos) InsertManyInto(posts []*models.Post, createdAt string) ([]*models.Post, error) {
	var queryStringAdditional string
	var args []interface{}

	queryStringMain := "INSERT INTO posts (author_nickname, forum_slug, message, thread_id, id, parent, created_at) VALUES "

	chunks := pr.makeChunks(posts)

	numb := 1

	for _, chunk := range chunks {
		ids := pr.postIdsGenerator.Next(len(chunk))
		for i, pst := range chunk {
			queryStringAdditional = ""

			pst.Id = ids[i]

			queryStringAdditional = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				numb, numb + 1, numb + 2, numb + 3, numb + 4, numb + 5, numb + 6)

			if i + 1 < len(chunk) {
				queryStringAdditional += ","
			}

			queryStringMain += queryStringAdditional

			pst.Created = createdAt

			args = append(args, pst.AuthorNickname, pst.ForumSlug, pst.Message, pst.Thread, pst.Id, pst.Parent, pst.Created)

			numb = numb + 7
		}

		_, err := pr.conn.Exec(queryStringMain, args...)

		if err != nil {
			return nil, err
		}
	}

	return posts, nil
}

func (pr *PostPgRepos) makeChunks(posts []*models.Post) [][]*models.Post {
	postsChunk := 100
	var chunks [][]*models.Post

	for i := 0; i < len(posts); i += postsChunk {
		bound := i + postsChunk

		if bound > len(posts) {
			bound = len(posts)
		}

		chunks = append(chunks, posts[i:bound])
	}

	return chunks
}

func (pr *PostPgRepos) Update(updatedPost *models.Post) {
	_, _ = pr.conn.Exec(`UPDATE threads SET author_nickname = $1,
                   message = $2 where id = $3`,
		updatedPost.AuthorNickname,
		updatedPost.Message,
		updatedPost.Id)
}

func (pr *PostPgRepos) SelectPostsFlat(id int, limit int, desc bool, since string) ([]*models.Post, error) {
	var posts []*models.Post

	queryString := "SELECT * from posts where thread_id = $1"
	var values []interface{}
	values = append(values, id)

	query, val := MakeQuery(values, queryString, limit, desc, since)

	if err := pr.conn.Select(&posts, query, val...); err != nil {
		return nil, PgxErrToCustom(err)
	}

	return posts, nil
}

func (pr *PostPgRepos) SelectPostsTree(id int, limit int, desc bool, since string) ([]*models.Post, error) {
	var posts []*models.Post

	queryString := "SELECT * from posts where thread_id = $1"
	var values []interface{}
	values = append(values, id)

	query, val := MakeQuery(values, queryString, limit, desc, since)

	if err := pr.conn.Select(&posts, query, val...); err != nil {
		return nil, PgxErrToCustom(err)
	}

	return posts, nil
}
//
//func (pr *PostPgRepos) SelectPostsParentTree(id int, limit int, desc bool, since string) ([]*models.Post, error) {
//
//}
