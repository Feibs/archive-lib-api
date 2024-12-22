package handler

import (
	"archive_lib/apperror"
	"archive_lib/dto"
	"archive_lib/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BorrowHandler struct {
	usecase usecase.BorrowUsecase
}

func NewBorrowHandler(uc usecase.BorrowUsecase) BorrowHandler {
	return BorrowHandler{
		usecase: uc,
	}
}

func (h BorrowHandler) BorrowBookHandler(ctx *gin.Context) {
	rawUserId, found := ctx.Get("subject")
	if !found {
		ctx.Error(apperror.ErrRequestUnrecognized{})
		return
	}

	userId, err := strconv.Atoi(rawUserId.(string))
	if err != nil {
		ctx.Error(apperror.ErrRequestUnrecognized{})
		return
	}

	var borrowRequest dto.BorrowRequest
	borrowRequest.UserId = userId
	err = ctx.ShouldBindJSON(&borrowRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	borrowResponse, err := h.usecase.Record(ctx, &borrowRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": borrowResponse})
}

func (h BorrowHandler) ReturnBookHandler(ctx *gin.Context) {
	rawUserId, found := ctx.Get("subject")
	if !found {
		ctx.Error(apperror.ErrRequestUnrecognized{})
		return
	}

	userId, err := strconv.Atoi(rawUserId.(string))
	if err != nil {
		ctx.Error(apperror.ErrRequestUnrecognized{})
		return
	}

	var returnRequest dto.ReturnRequest
	returnRequest.UserId = userId
	err = ctx.ShouldBindJSON(&returnRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	returnResponse, err := h.usecase.Return(ctx, &returnRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": returnResponse})
}
