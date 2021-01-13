package models

type Post struct {
	Id             int    `json:"id"`
	AuthorNickname string `json:"author" db:"author_nickname"`
	Created        string `json:"created" db:"created_at"`
	ForumSlug      string `json:"forum" db:"forum_slug"`
	Message        string `json:"message"`
	IsEdited       bool   `json:"isEdited" db:"is_edited"`
	Thread         int    `json:"thread,omitempty" db:"thread_id"`
	Parent         int    `json:"parent"`
}
