package mwares

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"time"
)

type Middleware func(fasthttp.RequestHandler) fasthttp.RequestHandler

func Use(mainHandler fasthttp.RequestHandler, mwares ...Middleware) fasthttp.RequestHandler {
	for i := len(mwares) - 1; i >= 0; i-- {
		mainHandler = mwares[i](mainHandler)
	}

	return mainHandler
}

func PanicRecovering(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Info(err)
			}
		}()

		next(ctx)
	}
}

func SetHeaders(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-type", "application/json")

		next(ctx)
	}
}

func AccessLog(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		startTime := time.Now()

		next(ctx)

		fmt.Println()
		logrus.Info(ctx.RemoteAddr(), " ",
			string(ctx.Method()), " ",
			string(ctx.Request.URI().FullURI()),
			" Work time: ", time.Now().Sub(startTime))
	}
}
