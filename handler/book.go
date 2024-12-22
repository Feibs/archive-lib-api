package handler

import (
	"archive_lib/dto"
	"archive_lib/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	usecase usecase.BookUsecase
}

func NewBookHandler(uc usecase.BookUsecase) BookHandler {
	return BookHandler{
		usecase: uc,
	}
}

func (h BookHandler) GetBooksHandler(ctx *gin.Context) {
	if title, found := ctx.GetQuery("title"); found {
		booksResponse, err := h.usecase.GetBooksByTitle(ctx, title)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": booksResponse})
		return
	}

	booksResponse, err := h.usecase.ListBooks(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": booksResponse})
}

func (h BookHandler) AddBookHandler(ctx *gin.Context) {
	var bookRequest dto.BookRequest
	err := ctx.ShouldBindJSON(&bookRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	bookResponse, err := h.usecase.AddBook(ctx, &bookRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": bookResponse})
}
