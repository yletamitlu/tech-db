package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/yletamitlu/tech-db/internal/consts"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/post"
	"strconv"
	"strings"
)

const (
	pathItemLen = 8
	nullPathItem = "00000000"
	maxNesting = 5
	pathItemsSeparator = "."
)

type PostPgRepos struct {
	conn             *sqlx.DB
	postIdsGenerator *Generator
}

func NewPostRepository(conn *sqlx.DB) post.PostRepository {
	return &PostPgRepos{
		conn:             conn,
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

func (pr *PostPgRepos) InsertManyInto(posts []*models.Post, thread *models.Thread, createdAt string) ([]*models.Post, error) {
	var queryStringAdditional string
	var args []interface{}

	queryStringMain := "INSERT INTO posts (author_nickname, forum_slug, message, thread_id, id, parent, created_at, path) VALUES "

	chunks := pr.makeChunks(posts)

	numb := 1

	for _, chunk := range chunks {
		ids := pr.postIdsGenerator.Next(len(chunk))
		for i, pst := range chunk {
			queryStringAdditional = ""

			pst.Id = ids[i]

			queryStringAdditional = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				numb, numb+1, numb+2, numb+3, numb+4, numb+5, numb+6, numb+7)

			if i + 1 < len(chunk) {
				queryStringAdditional += ","
			}

			queryStringMain += queryStringAdditional

			pst.Created = createdAt
			pst.ForumSlug = thread.ForumSlug
			pst.Thread = thread.Id
			path, err := pr.createPath(pst.Id, pst.Parent)

			if err != nil {
				return nil, err
			}

			pst.Path = path

			args = append(args, pst.AuthorNickname,
				pst.ForumSlug, pst.Message, pst.Thread,
				pst.Id, pst.Parent, pst.Created, pst.Path)

			numb = numb + 8
		}

		_, err := pr.conn.Exec(queryStringMain, args...)

		if err != nil {
			return nil, err
		}
	}

	return posts, nil
}

func (pr *PostPgRepos) createPath(postId int, parentId int) (string, error) {
	currentIdStr := strconv.Itoa(postId)
	pathItem := strings.Repeat("0", pathItemLen - len(currentIdStr)) + currentIdStr

	if parentId == 0 {
		pathItems := []string{pathItem}

		for i := 0; i < maxNesting - 1; i++ {
			pathItems = append(pathItems, nullPathItem)
		}

		return strings.Join(pathItems, pathItemsSeparator), nil
	}

	foundParent, err := pr.SelectById(parentId)

	if err != nil {
		return "", consts.ErrNotFound
	}

	return strings.Replace(foundParent.Path, nullPathItem, pathItem, 1), nil
}

func (pr *PostPgRepos) extractParentPath(path string) string {
	parentPathItem := strings.Split(path, pathItemsSeparator)[0]

	pathItems := []string{parentPathItem}

	for i := 0; i < maxNesting - 1; i++ {
		pathItems = append(pathItems, nullPathItem)
	}

	return strings.Join(pathItems, pathItemsSeparator)
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

func (pr *PostPgRepos) Update(updatedPost *models.Post) error {
	_, err := pr.conn.Exec(`UPDATE posts SET message = $1, is_edited = true where id = $2`,
		updatedPost.Message,
		updatedPost.Id)

	if err != nil {
		return err
	}

	return nil
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
	i := 1
	values = append(values, id)
	i++

	if since != "" {
		postId, _ := strconv.Atoi(since)
		foundPost, err := pr.SelectById(postId)

		if err != nil {
			return nil, err
		}

		if desc {
			queryString += fmt.Sprintf(" and path < $%d", i)
		} else {
			queryString += fmt.Sprintf(" and path > $%d", i)
		}

		values = append(values, foundPost.Path)
		i++
	}

	queryString += " order by path"

	if desc {
		queryString += " desc"
	}

	queryString += fmt.Sprintf(" limit $%d", i)
	values = append(values, limit)

	q := queryString
	if err := pr.conn.Select(&posts, q, values...); err != nil {
		return nil, PgxErrToCustom(err)
	}

	return posts, nil
}

func (pr *PostPgRepos) SelectPostsParentTree(id int, limit int, desc bool, since string) ([]*models.Post, error) {
	var parentPosts []*models.Post

	queryString := "SELECT * from posts where thread_id = $1 and parent = 0"

	var values []interface{}
	i := 1
	values = append(values, id)
	i++

	if since != "" {
		postId, _ := strconv.Atoi(since)
		foundPost, err := pr.SelectById(postId)

		if err != nil {
			return nil, err
		}

		if desc {
			queryString += fmt.Sprintf(" and path < $%d", i)
		} else {
			queryString += fmt.Sprintf(" and path > $%d", i)
		}

		values = append(values, pr.extractParentPath(foundPost.Path))
		i++
	}

	queryString += " order by id"

	if desc {
		queryString += " desc"
	}

	queryString += fmt.Sprintf(" limit $%d", i)
	values = append(values, limit)

	q := queryString
	if err := pr.conn.Select(&parentPosts, q, values...); err != nil {
		return nil, PgxErrToCustom(err)
	}

	var posts []*models.Post

	for _, parent := range parentPosts {
		var children []*models.Post

		parentPathItem := strings.Split(parent.Path, pathItemsSeparator)[0]

		if err := pr.conn.Select(&children, "SELECT * FROM posts where substring(path, 1, 8) = $1 order by path",
			parentPathItem);
		err != nil {
			return nil, PgxErrToCustom(err)
		}

		posts = append(posts, children...)
	}

	return posts, nil
}
