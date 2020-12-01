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
	"log"
)

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
	forumUcase := forumU.NewUserUcase(forumRepos)
	forumDelivery := forumD.NewForumDelivery(forumUcase)

	userDelivery.Configure(router)
	forumDelivery.Configure(router)

	fmt.Printf("Server started...")
	log.Fatal(fasthttp.ListenAndServe(":5000", Use(router.Handler, PanicRecovering, SetHeaders, AccessLog)))
}
