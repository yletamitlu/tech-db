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
	router.POST("/api/thread/:slug/create", pd.createThreadHandler())
	router.GET("/api/post/:slug/details", pd.getThreadsHandler())
}

func (pd *PostDelivery) createThreadHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		slugStr, _ := ctx.UserValue("slug").(string)

		var posts []models.Post

		body := ctx.Request.Body()
		if string(body) == "[]\n" {
			var posts []*models.Post
			posts = []*models.Post{}
			SendResponse(ctx, 201, posts)
			return
		}

		if err := json.Unmarshal(body, &posts); err != nil {
			logrus.Info(err)
			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		var resultPosts []*models.Post
		for _, pst := range posts {
			resultPost, err := pd.postUcase.Create(&pst, slugStr)

			if err != nil {
				logrus.Info(err)
				SendResponse(ctx, 500, &ErrorResponse{
					Message: ErrInternal.Error(),
				})
				return
			}

			resultPosts = append(resultPosts, resultPost)
		}

		SendResponse(ctx, 201, resultPosts)
		return
	}
}

func (pd *PostDelivery) getThreadsHandler() fasthttp.RequestHandler {
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
