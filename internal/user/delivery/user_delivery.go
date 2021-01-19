package delivery

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"

	"github.com/valyala/fasthttp"
	. "github.com/yletamitlu/tech-db/internal/consts"
	. "github.com/yletamitlu/tech-db/internal/helpers"
	"github.com/yletamitlu/tech-db/internal/models"
	"github.com/yletamitlu/tech-db/internal/user"
	. "github.com/yletamitlu/tech-db/tools"
)

type UserDelivery struct {
	userUcase user.UserUsecase
}

func NewUserDelivery(uUc user.UserUsecase) *UserDelivery {
	return &UserDelivery{
		userUcase: uUc,
	}
}

func (ud *UserDelivery) Configure(router *fasthttprouter.Router) {
	router.POST("/api/user/:nickname/create", ud.createUserHandler())
	router.GET("/api/user/:nickname/profile", ud.getUserProfile())
	router.POST("/api/user/:nickname/profile", ud.updateProfile())
}

func (ud *UserDelivery) getUserProfile() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		nickname, _ := ctx.UserValue("nickname").(string)

		u, err := ud.userUcase.GetByNickname(nickname)

		if err != nil && err == ErrNotFound {

			SendResponse(ctx, 404, &ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		if err != nil {

			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		SendResponse(ctx, 200, u)
	}
}

func (ud *UserDelivery) createUserHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		nickname, _ := ctx.UserValue("nickname").(string)

		u := &models.User{
			Nickname: nickname,
		}

		body := ctx.PostBody()
		if err := json.Unmarshal(body, &u); err != nil {
			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		err, found := ud.userUcase.Create(u)

		if found != nil {

			SendResponse(ctx, 409, found)
			return
		}

		if err != nil {

			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		SendResponse(ctx, 201, u)
		return
	}
}

func (ud *UserDelivery) updateProfile() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		nickname, _ := ctx.UserValue("nickname").(string)

		updatedUser := &models.User{
			Nickname: nickname,
		}

		body := ctx.PostBody()
		if err := json.Unmarshal(body, &updatedUser); err != nil {

			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		resU, err := ud.userUcase.Update(updatedUser)

		if err != nil && err == ErrNotFound {

			SendResponse(ctx, 404, &ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		if err != nil && err == ErrConflict {

			SendResponse(ctx, 409, &ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		if err != nil {

			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		if resU != nil {
			SendResponse(ctx, 200, resU)
			return 
		}

		SendResponse(ctx, 200, updatedUser)
	}
}
