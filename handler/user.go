package handler

import (
	"context"
	"archive_lib/config"
	"archive_lib/dto"
	"archive_lib/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) UserHandler {
	return UserHandler{
		usecase: uc,
	}
}

func (h UserHandler) Login(ctx *gin.Context) {
	span := config.StartSpanFromRequest(config.Tracer, ctx.Request, "archive-lib-service handler Login func")
	defer span.Finish()
	ctxTracer := opentracing.ContextWithSpan(context.Background(), span)

	var authRequest dto.AuthRequest
	err := ctx.ShouldBindJSON(&authRequest)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("error", err)
		ctx.Error(err)
		return
	}

	accessToken, err := h.usecase.Login(ctxTracer, &authRequest)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("error", err)
		ctx.Error(err)
		return
	}

	span.LogKV("info", "handler Login success")

	ctx.JSON(http.StatusOK, gin.H{"data": accessToken})
}
