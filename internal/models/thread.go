package models

import "time"

type Thread struct {
	Id             int       `json:"id"`
	AuthorNickname string    `json:"author" db:"author_nickname"`
	Created        time.Time `json:"created" db:"created_at"`
	ForumSlug      string    `json:"forum" db:"forum_slug"`
	Message        string    `json:"message"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug,omitempty"`
	Votes          int       `json:"votes"`
}
