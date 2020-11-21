package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
	"log"
	. "github.com/yletamitlu/tech-db/internal/handlers"
)

// нужна модель пользователя
// разобраться в библиотеке json marshal unmarshal
// разобраться как доставать параметры из запроса (из тела, из path параметров)
// написиат обработчик для создания пользователя - протестировать запуском сервера и прогоном тестов (1 тест должен пройти)
// не забывать про content-type application/json в response
// го лэнг struct json tag

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://techdbuser@localhost:5432/techdb")
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	defer conn.Close(context.Background())

	router := fasthttprouter.New()

	userHandler := MakeUserHandler(conn)
	router.POST("/api/user/:nickname/create", userHandler.CreateUser)

	log.Fatal(fasthttp.ListenAndServe(":5000", router.Handler))
}
