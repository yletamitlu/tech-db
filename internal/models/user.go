package models

type User struct {
	Id       int    `json:"-"`
	About    string `json:"about"`
	Email    string `json:"email"`
	FullName string `json:"fullname"`
	Nickname string `json:"nickname"`
}
