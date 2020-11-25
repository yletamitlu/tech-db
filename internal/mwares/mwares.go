package mwares

import (
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)


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
