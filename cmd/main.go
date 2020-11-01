package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

func requestHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetContentType("application/json")
}

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://techdbuser@localhost:5432/techdb")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	if err = conn.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Fatal(fasthttp.ListenAndServe(":5000", requestHandler))
}
