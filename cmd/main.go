package main

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
	. "github.com/yletamitlu/tech-db/internal/mwares"
	"github.com/yletamitlu/tech-db/internal/user/delivery"
	"github.com/yletamitlu/tech-db/internal/user/repository"
	"github.com/yletamitlu/tech-db/internal/user/usecase"
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

	userRepos := repository.NewUserRepository(conn)
	userUcase := usecase.NewUserUcase(userRepos)
	userDelivery := delivery.NewUserDelivery(userUcase)

	userDelivery.Configure(router)

	fmt.Printf("Server started...")
	log.Fatal(fasthttp.ListenAndServe(":5000", Use(router.Handler, PanicRecovering, SetHeaders, AccessLog)))
}
