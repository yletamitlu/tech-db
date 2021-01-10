package models

type Forum struct {
	AuthorNickname string `json:"user" db:"author_nickname"`
	Title          string `json:"title"`
	Slug           string `json:"slug"`
	Threads        int    `json:"threads"`
	Posts          int    `json:"posts"`
}
