package delivery

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	. "github.com/yletamitlu/tech-db/internal/consts"
	"github.com/yletamitlu/tech-db/internal/forum"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	"github.com/yletamitlu/tech-db/internal/models"
	. "github.com/yletamitlu/tech-db/tools"
)

type ForumDelivery struct {
	forumUcase forum.ForumUsecase
}

func NewForumDelivery(fUc forum.ForumUsecase) *ForumDelivery {
	return &ForumDelivery{
		forumUcase: fUc,
	}
}

func (fd *ForumDelivery) Configure(router *fasthttprouter.Router) {
	router.POST("/api/forum/create", fd.createForumHandler())
}

func (fd *ForumDelivery) createForumHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		f := &models.Forum{}

		body := ctx.Request.Body()
		if err := json.Unmarshal(body, &f); err != nil {
			logrus.Info(err)
			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		err, found := fd.forumUcase.Create(f)

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

		SendResponse(ctx, 201, f)
		return
	}
}