package delivery

import (
	"github.com/buaazp/fasthttprouter"

	"github.com/valyala/fasthttp"
	"github.com/yletamitlu/tech-db/internal/additions"
	. "github.com/yletamitlu/tech-db/internal/consts"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	. "github.com/yletamitlu/tech-db/tools"
)

type AdditionsDelivery struct {
	additionsRepo additions.AdditionRepository
}

func NewAdditionsDelivery(ar additions.AdditionRepository) *AdditionsDelivery {
	return &AdditionsDelivery{
		additionsRepo: ar,
	}
}

func (ad *AdditionsDelivery) Configure(router *fasthttprouter.Router) {
	router.GET("/api/service/status", ad.Clear())
	router.POST("/api/service/clear", ad.Status())
}

func (ad *AdditionsDelivery) Clear() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		forumStatus, threadStatus, postStatus, userStatus, err := ad.additionsRepo.Status()

		if err != nil {

			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		type Statuses struct {
			User   uint64 `json:"user"`
			Forum  uint64 `json:"forum"`
			Thread uint64 `json:"thread"`
			Post   uint64 `json:"post"`
		}

		statuses := &Statuses{
			User: userStatus,
			Forum: forumStatus,
			Thread: threadStatus,
			Post: postStatus,
		}

		SendResponse(ctx, 200, statuses)
		return
	}
}

func (ad *AdditionsDelivery) Status() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		err := ad.additionsRepo.Clear()

		if err != nil {

			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		SendResponse(ctx, 200, nil)
		return
	}
}
