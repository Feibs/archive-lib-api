package setup

import (
	"archive_lib/handler"
	"archive_lib/middleware"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	userHandler   *handler.UserHandler
	bookHandler   *handler.BookHandler
	borrowHandler *handler.BorrowHandler
}

func NewHandlers(userHandler *handler.UserHandler, bookHandler *handler.BookHandler, borrowHandler *handler.BorrowHandler) *Handlers {
	return &Handlers{
		userHandler,
		bookHandler,
		borrowHandler,
	}
}

func NewRouter(h *Handlers) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware)
	router.Use(middleware.ErrorMiddleware)

	router.POST("/login", h.userHandler.Login)
	router.GET("/books", h.bookHandler.GetBooksHandler)
	router.POST("/books", h.bookHandler.AddBookHandler)
	router.POST("/borrowing-records", middleware.AuthMiddleware, h.borrowHandler.BorrowBookHandler)
	router.PATCH("/borrowing-records", middleware.AuthMiddleware, h.borrowHandler.ReturnBookHandler)

	return router
}
