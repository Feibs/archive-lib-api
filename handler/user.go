package handler

import (
	"archive_lib/dto"
	"archive_lib/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
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
	var authRequest dto.AuthRequest
	err := ctx.ShouldBindJSON(&authRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	accessToken, err := h.usecase.Login(ctx, &authRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": accessToken})
}
