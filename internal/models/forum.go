package models

type Forum struct {
	AuthorNickname string `json:"user"`
	Title          string `json:"title"`
	Slug           string `json:"slug"`
	Threads        int    `json:"threads"`
	Posts          int    `json:"posts"`
}
