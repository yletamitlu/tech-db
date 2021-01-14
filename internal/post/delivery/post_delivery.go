package delivery

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	. "github.com/yletamitlu/tech-db/internal/consts"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/post"
	. "github.com/yletamitlu/tech-db/tools"
	"strconv"
	"strings"
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
	router.GET("/api/thread/:slug/posts", pd.getPostsHandler())
	router.POST("/api/post/:id/details", pd.updatePost())
	router.GET("/api/post/:id/details", pd.getPostHandler())
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

		if err == ErrAlreadyExists || err == ErrConflict {
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
		idStr, _ := ctx.UserValue("id").(string)
		related := string(ctx.QueryArgs().Peek("related"))

		id, _ := strconv.Atoi(idStr)

		found, err := pd.postUcase.GetById(id)

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

		type wrapper struct {
			Post   *models.Post   `json:"post"`
			Author *models.User   `json:"author,omitempty"`
			Thread *models.Thread `json:"thread,omitempty"`
			Forum  *models.Forum  `json:"forum,omitempty"`
		}

		postWrapper := &wrapper{
			Post: found,
		}

		if related != "" {
			if strings.Contains(related, "user") {
				foundAuthor, err := pd.postUcase.GetPostAuthor(found.AuthorNickname)

				if err != nil {
					logrus.Info(err)
					SendResponse(ctx, 500, &ErrorResponse{
						Message: ErrInternal.Error(),
					})
					return
				}

				postWrapper.Author = foundAuthor
			}

			if strings.Contains(related, "thread") {
				foundThread, err := pd.postUcase.GetPostThread(found.Thread)

				if err != nil {
					logrus.Info(err)
					SendResponse(ctx, 500, &ErrorResponse{
						Message: ErrInternal.Error(),
					})
					return
				}

				postWrapper.Thread = foundThread
			}

			if strings.Contains(related, "forum") {
				foundForum, err := pd.postUcase.GetPostForum(found.ForumSlug)

				if err != nil {
					logrus.Info(err)
					SendResponse(ctx, 500, &ErrorResponse{
						Message: ErrInternal.Error(),
					})
					return
				}

				postWrapper.Forum = foundForum
			}
		}

		SendResponse(ctx, 200, postWrapper)
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

		if posts == nil {
			var posts []*models.Post
			posts = []*models.Post{}
			SendResponse(ctx, 200, posts)
			return
		}

		SendResponse(ctx, 200, posts)
		return
	}
}

func (pd *PostDelivery) updatePost() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		idStr, _ := ctx.UserValue("id").(string)

		id, _ := strconv.Atoi(idStr)
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
