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
	"github.com/yletamitlu/tech-db/internal/post"
	. "github.com/yletamitlu/tech-db/tools"
	"strconv"
)

type PostDelivery struct {
	postUcase post.PostUsecase
}

func NewPostDelivery(pUc post.PostUsecase) *PostDelivery {
	return &PostDelivery{
		postUcase: pUc,
	}
}

func (pd *PostDelivery) Configure(router *fasthttprouter.Router) {
	router.POST("/api/thread/:slug/create", pd.createPostHandler())
	router.GET("/api/post/:slug/details", pd.getPostHandler())
	router.GET("/api/thread/:slug/posts", pd.getPostsHandler())
	router.PUT("/api/post/:id/details", pd.updatePost())
}

func (pd *PostDelivery) createPostHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		slugStr, _ := ctx.UserValue("slug").(string)

		var posts []*models.Post

		err := json.Unmarshal(ctx.Request.Body(), &posts)

		if err != nil {
			logrus.Info(err)
			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		resultPosts, err := pd.postUcase.Create(posts, slugStr)

		if err == ErrNotFound {
			logrus.Info(err)
			SendResponse(ctx, 404, &ErrorResponse{
				Message: ErrNotFound.Error(),
			})
			return
		}

		if err == ErrAlreadyExists {
			logrus.Info(err)
			SendResponse(ctx, 409, &ErrorResponse{
				Message: ErrAlreadyExists.Error(),
			})
			return
		}

		if resultPosts == nil {
			var posts []*models.Post
			posts = []*models.Post{}
			SendResponse(ctx, 201, posts)
			return
		}

		if err != nil {
			logrus.Info(err)
			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		SendResponse(ctx, 201, resultPosts)
		return
	}
}

func (pd *PostDelivery) getPostHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Println(ctx.URI())
		slug, _ := ctx.UserValue("slug").(string)

		found, err := pd.postUcase.GetByForumSlug(slug)

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

func (pd *PostDelivery) getPostsHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		slugOrId, _ := ctx.UserValue("slug").(string)
		limitStr := string(ctx.QueryArgs().Peek("limit"))
		descStr := string(ctx.QueryArgs().Peek("desc"))
		since := string(ctx.QueryArgs().Peek("since"))
		sort := string(ctx.QueryArgs().Peek("sort"))

		limit, _ := strconv.Atoi(limitStr)
		desc, _ := strconv.ParseBool(descStr)

		posts, err := pd.postUcase.GetPosts(slugOrId, limit, desc, since, sort)

		if posts == nil {
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

		SendResponse(ctx, 200, posts)
		return
	}
}

func (pd *PostDelivery) updatePost() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		id, _ := ctx.UserValue("id").(int)

		pst := &models.Post{
			Id: id,
		}

		err := json.Unmarshal(ctx.Request.Body(), &pst)

		if err != nil {
			logrus.Info(err)
			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		pst, err = pd.postUcase.Update(pst)

		if err == ErrNotFound {
			logrus.Info(err)
			SendResponse(ctx, 404, &ErrorResponse{
				Message: ErrNotFound.Error(),
			})
			return
		}

		SendResponse(ctx, 200, pst)
		return
	}
}