package delivery

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	. "github.com/yletamitlu/tech-db/internal/consts"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/threads"
	. "github.com/yletamitlu/tech-db/tools"
	"strconv"
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

	router.GET("/api/thread/:slug/details", td.getThreadDetailsHandler())
	router.POST("/api/thread/:slug/details", td.updateThreadsHandler())
	router.POST("/api/thread/:slug/vote", td.voteThreadsHandler())
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

		found, err := td.threadUcase.Create(thr)

		if found != nil {
			logrus.Info(err)
			SendResponse(ctx, 409, found)
			return
		}

		if err == ErrNotFound {
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

		SendResponse(ctx, 201, thr)
		return
	}
}

func (td *ThreadDelivery) getThreadsHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		slug, _ := ctx.UserValue("slug").(string)
		limitStr := string(ctx.QueryArgs().Peek("limit"))
		descStr := string(ctx.QueryArgs().Peek("desc"))
		since := string(ctx.QueryArgs().Peek("since"))

		limit, _ := strconv.Atoi(limitStr)
		desc, _ := strconv.ParseBool(descStr)

		found, err := td.threadUcase.GetByForumSlug(slug, limit, desc, since)

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

func (td *ThreadDelivery) getThreadDetailsHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		slugOrId, _ := ctx.UserValue("slug").(string)

		id, err := strconv.Atoi(slugOrId)

		var foundThr *models.Thread
		if err == nil {
			foundThr, err = td.threadUcase.GetById(id)
		} else {
			foundThr, err = td.threadUcase.GetBySlug(slugOrId)
		}

		if foundThr == nil {
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

		SendResponse(ctx, 200, foundThr)
		return
	}
}

func (td *ThreadDelivery) updateThreadsHandler() fasthttp.RequestHandler {
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

		changed, err := td.threadUcase.Update(thr)

		if err != nil {
			logrus.Info(err)
			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		SendResponse(ctx, 201, changed)
		return
	}
}

func (td *ThreadDelivery) voteThreadsHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		slug, _ := ctx.UserValue("slug").(string)

		vote := &models.Vote{}

		body := ctx.Request.Body()
		if err := json.Unmarshal(body, &vote); err != nil {
			logrus.Info(err)
			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		changed, err := td.threadUcase.CreateVote(vote, slug)

		if err != nil {
			logrus.Info(err)
			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		if vote.Voice == -1 && changed.Votes == 2 {
			changed.Votes = 1
		}

		SendResponse(ctx, 200, changed)
		return
	}
}
