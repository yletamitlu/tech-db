package delivery

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"
	"github.com/sirupsen/logrus"
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
}

func (ud *UserDelivery) createUserHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		nickname, _ := ctx.UserValue("nickname").(string)

		u := &models.User{
			Nickname: nickname,
		}

		body := ctx.Request.Body()
		if err := json.Unmarshal(body, &u); err != nil {
			logrus.Info(err)
			SendResponse(ctx, 500, &ErrorResponse{
				Message: ErrInternal.Error(),
			})
			return
		}

		err, found := ud.userUcase.Create(u)

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

		SendResponse(ctx, 201, u)
		return
	}
}
