package delivery

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"

	"github.com/valyala/fasthttp"
	. "github.com/yletamitlu/tech-db/internal/consts"
	"github.com/yletamitlu/tech-db/internal/forum"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	"github.com/yletamitlu/tech-db/internal/models"
	. "github.com/yletamitlu/tech-db/tools"
	"strconv"
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
	// /api/forum/create в CrutchRouter'e
	router.GET("/api/forum/:slug/details", fd.getForumDetailsHandler())
	router.GET("/api/forum/:slug/users", fd.getForumUsersHandler())
}

func (fd *ForumDelivery) CreateForumHandler(ctx *fasthttp.RequestCtx) {
	f := &models.Forum{}

	body := ctx.PostBody()
	if err := json.Unmarshal(body, &f); err != nil {

		SendResponse(ctx, 500, &ErrorResponse{
			Message: ErrInternal.Error(),
		})
		return
	}

	f, err := fd.forumUcase.Create(f)

	if err == ErrAlreadyExists {

		SendResponse(ctx, 409, f)
		return
	}

	if err == ErrNotFound {

		SendResponse(ctx, 404, &ErrorResponse{
			Message: ErrNotFound.Error(),
		})
		return
	}

	if err != nil {

		SendResponse(ctx, 500, &ErrorResponse{
			Message: ErrInternal.Error(),
		})
		return
	}

	SendResponse(ctx, 201, f)
	return
}

func (fd *ForumDelivery) getForumDetailsHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		slug, _ := ctx.UserValue("slug").(string)

		found, err := fd.forumUcase.GetBySlug(slug, true)

		if err == ErrNotFound {

			SendResponse(ctx, 404, &ErrorResponse{
				Message: ErrNotFound.Error(),
			})
			return
		}

		if err != nil {

			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		SendResponse(ctx, 200, found)
		return
	}
}

func (fd *ForumDelivery) getForumUsersHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		slug, _ := ctx.UserValue("slug").(string)

		limitStr := string(ctx.FormValue("limit"))
		descStr := string(ctx.FormValue("desc"))
		since := string(ctx.FormValue("since"))

		limit, _ := strconv.Atoi(limitStr)
		desc, _ := strconv.ParseBool(descStr)

		found, err := fd.forumUcase.GetUsers(slug, limit, desc, since)

		if found == nil {

			SendResponse(ctx, 404, &ErrorResponse{
				Message: ErrNotFound.Error(),
			})
			return
		}

		if err != nil {

			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		SendResponse(ctx, 200, found)
		return
	}
}
