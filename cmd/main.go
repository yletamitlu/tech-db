package main

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
	forumD "github.com/yletamitlu/tech-db/internal/forum/delivery"
	forumR "github.com/yletamitlu/tech-db/internal/forum/repository"
	forumU "github.com/yletamitlu/tech-db/internal/forum/usecase"
	. "github.com/yletamitlu/tech-db/internal/mwares"
	userD "github.com/yletamitlu/tech-db/internal/user/delivery"
	userR "github.com/yletamitlu/tech-db/internal/user/repository"
	userU "github.com/yletamitlu/tech-db/internal/user/usecase"

	threadD "github.com/yletamitlu/tech-db/internal/threads/delivery"
	threadR "github.com/yletamitlu/tech-db/internal/threads/repository"
	threadU "github.com/yletamitlu/tech-db/internal/threads/usecase"
	"log"
)

type CrutchRouter struct {
	r *fasthttprouter.Router
	fd *forumD.ForumDelivery
}

func NewCrutchRouter(r *fasthttprouter.Router, fd *forumD.ForumDelivery) *CrutchRouter {
	return &CrutchRouter{
		r: r,
		fd: fd,
	}
}

func (cr *CrutchRouter) GetHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())

		if path == "/api/forum/create" {
			cr.fd.CreateForumHandler(ctx)
		} else {
			cr.r.Handler(ctx)
		}
	}
}

func main() {
	conn, err := sqlx.Connect("pgx", "postgres://techdbuser@localhost:5432/techdb")
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	router := fasthttprouter.New()

	userRepos := userR.NewUserRepository(conn)
	userUcase := userU.NewUserUcase(userRepos)
	userDelivery := userD.NewUserDelivery(userUcase)

	forumRepos := forumR.NewForumRepository(conn)
	forumUcase := forumU.NewForumUcase(forumRepos, userUcase)
	forumDelivery := forumD.NewForumDelivery(forumUcase)

	threadRepos := threadR.NewThreadRepository(conn)
	threadUcase := threadU.NewThreadUcase(threadRepos, userUcase)
	threadDelivery := threadD.NewThreadDelivery(threadUcase)

	userDelivery.Configure(router)
	threadDelivery.Configure(router)
	forumDelivery.Configure(router)

	crutchRouter := NewCrutchRouter(router, forumDelivery)

	fmt.Printf("Server started...")
	log.Fatal(fasthttp.ListenAndServe(":5000", Use(crutchRouter.GetHandler(), PanicRecovering, SetHeaders, AccessLog)))
}
