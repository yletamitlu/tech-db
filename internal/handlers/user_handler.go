package handlers

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/valyala/fasthttp"
	"log"
)

type UserHandler struct {
	conn *pgx.Conn
}

type User struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	FullName string `json:"fullname"`
	Nickname string `json:"nickname"`
}

func MakeUserHandler(c *pgx.Conn) *UserHandler {
	return &UserHandler{
		conn: c,
	}
}

func (uh *UserHandler) parseQueryRows(rows pgx.Rows) []User {
	var users []User
	rowNum := 1
	for rows.Next() {
		var v []interface{}
		v, err := rows.Values()
		if err != nil {
			log.Fatal(err)
		}

		user := User{}
		user.Nickname = v[0].(string)
		user.FullName = v[1].(string)
		user.Email = v[2].(string)
		user.About = v[3].(string)

		users = append(users, user)
		rowNum++
	}
	return users
}

func (uh *UserHandler) CreateUser(ctx *fasthttp.RequestCtx) {
	nickname, _ := ctx.UserValue("nickname").(string)

	user := &User{}

	user.Nickname = nickname
	body := ctx.Request.Body()
	_ = json.Unmarshal(body, &user)

	rows, err := uh.conn.Query(context.Background(), `SELECT * from users where nickname=$1 or
														email=$2`, nickname, user.Email)

	users := uh.parseQueryRows(rows)

	ctx.Response.Header.Set("Content-type", "application/json")
	if len(users) > 0 {
		ctx.SetStatusCode(409)
		body, err := json.Marshal(users)
		if err == nil {
			ctx.SetBody(body)
		}
		return
	}

	_, err = uh.conn.Exec(context.Background(), `INSERT INTO users(nickname, fullname, email, about)
		VALUES ($1, $2, $3, $4)`, user.Nickname, user.FullName, user.Email, user.About)

	if err != nil {
		ctx.SetStatusCode(409)
		body, err = json.Marshal(users)
		if err == nil {
			ctx.SetBody(body)
		}
		return
	}

	ctx.SetStatusCode(201)
	body, err = json.Marshal(user)
	if err == nil {
		ctx.SetBody(body)
	}
	return
}
