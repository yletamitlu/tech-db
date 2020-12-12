package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	. "github.com/yletamitlu/tech-db/internal/consts"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/threads"
	. "github.com/yletamitlu/tech-db/tools"
	"strconv"
	"time"
)

type ThreadDelivery struct {
	threadUcase threads.ThreadUsecase
}

func NewThreadDelivery(tUc threads.ThreadUsecase) *ThreadDelivery {
	return &ThreadDelivery{
		threadUcase: tUc,
	}
}

func (td *ThreadDelivery) Configure(router *fasthttprouter.Router) {
	router.POST("/api/forum/:slug/create", td.createThreadHandler())
	router.GET("/api/forum/:slug/threads", td.getThreadsHandler())
}

func (td *ThreadDelivery) createThreadHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		slug, _ := ctx.UserValue("slug").(string)

		thr := &models.Thread{
			ForumSlug: slug,
		}

		body := ctx.Request.Body()
		if err := json.Unmarshal(body, &thr); err != nil {
			logrus.Info(err)
			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		err, found := td.threadUcase.Create(thr)

		if found != nil {
			logrus.Info(err)
			SendResponse(ctx, 409, found)
			return
		}

		if err != nil {
			logrus.Info(err)
			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		SendResponse(ctx, 201, thr)
		return
	}
}

func (td *ThreadDelivery) getThreadsHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Println(ctx.URI())
		slug, _ := ctx.UserValue("slug").(string)
		limitStr := string(ctx.QueryArgs().Peek("limit"))
		descStr := string(ctx.QueryArgs().Peek("desc"))
		sinceStr := string(ctx.QueryArgs().Peek("since"))

		limit, _ := strconv.Atoi(limitStr)
		desc, _ := strconv.ParseBool(descStr)
		since, err := time.Parse(time.RFC3339Nano, sinceStr)
		fmt.Print(since)

		found, err := td.threadUcase.GetByForumSlug(slug, limit, desc, since)

		dat, _ := td.threadUcase.GetByDate(slug, since)
		fmt.Print(dat)

		if found == nil {
			logrus.Info(err)
			SendResponse(ctx, 404, &ErrorResponse{
				Message: ErrNotFound.Error(),
			})
			return
		}

		if err != nil {
			logrus.Info(err)
			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		SendResponse(ctx, 200, found)
		return
	}
}
