package models

import "time"

type Post struct {
	Id             int       `json:"id"`
	AuthorNickname string    `json:"author" db:"author_nickname"`
	Created        time.Time `json:"created" db:"created_at"`
	ForumSlug      string    `json:"forum" db:"forum_slug"`
	Message        string    `json:"message"`
	IsEdited       bool      `json:"is_edited" db:"is_edited"`
	Thread         int       `json:"thread,omitempty" db:"thread_id"`
	Parent         int       `json:"parent"`
}
