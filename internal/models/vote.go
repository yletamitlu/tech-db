package models

type Vote struct {
	AuthorNickname string `json:"nickname" db:"user_nickname"`
	Thread         int    `json:"thread,omitempty" db:"thread_id"`
	Voice          int `json:"voice"`
}
