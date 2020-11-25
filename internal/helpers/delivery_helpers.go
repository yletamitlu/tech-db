package helpers

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

func SendResponse(ctx *fasthttp.RequestCtx, code int, content interface{}) {
	ctx.SetStatusCode(code)
	body, _ := json.Marshal(&content)
	if body != nil {
		ctx.SetBody(body)
	}
}
