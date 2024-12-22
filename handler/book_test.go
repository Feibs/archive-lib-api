package handler_test

import (
	"archive_lib/apperror"
	"archive_lib/dto"
	"archive_lib/handler"
	"archive_lib/middleware"
	"archive_lib/mocks"
	"archive_lib/util"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	authorId = 5
	quantity = 10
	cover    = "hardcover"

	authorResponse = &dto.AuthorResponse{
		Id:   authorId,
		Name: "John The Poet",
	}

	bookWithAuthorResponse = &dto.BookResponse{
		Id:          1,
		Author:      authorResponse,
		Title:       "Test Book",
		Description: "Cool book",
		Quantity:    quantity,
		Cover:       cover,
	}
	booksWithAuthorResponse = []*dto.BookResponse{bookWithAuthorResponse}

	bookResponse = &dto.BookResponse{
		Id:          1,
		Title:       "Test Book",
		Description: "Cool book",
		Quantity:    quantity,
		Cover:       cover,
	}

	bookRequest = &dto.BookRequest{
		Title:       "Test Book",
		AuthorId:    &authorId,
		Description: "Cool book",
		Quantity:    &quantity,
		Cover:       &cover,
	}
)

func TestGetBooksHandler(t *testing.T) {
	t.Run("should return StatusOK with list of books when no error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBookUsecase := new(mocks.BookUsecase)
		mockBookUsecase.On("ListBooks", ctx).Return(booksWithAuthorResponse, nil)
		bookHandler := handler.NewBookHandler(mockBookUsecase)
		router.GET("/books", bookHandler.GetBooksHandler)
		expectedResponse, _ := json.Marshal(gin.H{"data": booksWithAuthorResponse})

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/books", nil)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return StatusInternalServerError when list books encounters server error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBookUsecase := new(mocks.BookUsecase)
		mockBookUsecase.On("ListBooks", ctx).Return(nil, errors.New("server error"))
		bookHandler := handler.NewBookHandler(mockBookUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.GET("/books", bookHandler.GetBooksHandler)
		expectedResponse, _ := json.Marshal(gin.H{"message": "Server error"})

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/books", nil)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return StatusOK with list of books matching the input title when no error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBookUsecase := new(mocks.BookUsecase)
		mockBookUsecase.On("GetBooksByTitle", ctx, "tes").Return(booksWithAuthorResponse, nil)
		bookHandler := handler.NewBookHandler(mockBookUsecase)
		router.GET("/books", bookHandler.GetBooksHandler)
		expectedResponse, _ := json.Marshal(gin.H{"data": booksWithAuthorResponse})

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/books?title=tes", nil)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return StatusInternalServerError when search books encounters server error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBookUsecase := new(mocks.BookUsecase)
		mockBookUsecase.On("GetBooksByTitle", ctx, "any").Return(nil, errors.New("server error"))
		bookHandler := handler.NewBookHandler(mockBookUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.GET("/books", bookHandler.GetBooksHandler)
		expectedResponse, _ := json.Marshal(gin.H{"message": "Server error"})

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/books?title=any", nil)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})
}

func TestAddBookHandler(t *testing.T) {
	t.Run("should return StatusCreated with the added book when no error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBookUsecase := new(mocks.BookUsecase)
		mockBookUsecase.On("AddBook", ctx, bookRequest).Return(bookResponse, nil)
		bookHandler := handler.NewBookHandler(mockBookUsecase)
		router.POST("/books", bookHandler.AddBookHandler)
		expectedResponse, _ := json.Marshal(gin.H{"data": bookResponse})
		bookRequestJSON, _ := json.Marshal(*bookRequest)
		body := strings.NewReader(string(bookRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/books", body)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return StatusBadRequest when add book encounters decode error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBookUsecase := new(mocks.BookUsecase)
		mockBookUsecase.On("AddBook", ctx, bookRequest).Return(nil, errors.New("decode error"))
		bookHandler := handler.NewBookHandler(mockBookUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.POST("/books", bookHandler.AddBookHandler)
		errDecodeResponse := []util.FieldError{{Field: "", Message: "Mismatch data type or malformed request"}}
		expectedResponse, _ := json.Marshal(gin.H{"message": errDecodeResponse})
		bookRequestJSON, _ := json.Marshal("")
		body := strings.NewReader(string(bookRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/books", body)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return StatusInternalServerError when add book encounters server error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBookUsecase := new(mocks.BookUsecase)
		mockBookUsecase.On("AddBook", ctx, bookRequest).Return(nil, errors.New("server error"))
		bookHandler := handler.NewBookHandler(mockBookUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.POST("/books", bookHandler.AddBookHandler)
		expectedResponse, _ := json.Marshal(gin.H{"message": "Server error"})
		bookRequestJSON, _ := json.Marshal(*bookRequest)
		body := strings.NewReader(string(bookRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/books", body)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return StatusBadRequest when adding book with empty field", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBookUsecase := new(mocks.BookUsecase)
		invalidBookRequest := &dto.BookRequest{}
		mockBookUsecase.On("AddBook", ctx, invalidBookRequest).Return(nil, errors.New("error"))
		bookHandler := handler.NewBookHandler(mockBookUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.POST("/books", bookHandler.AddBookHandler)
		fieldErrors := []util.FieldError{
			{Field: "Title", Message: "Required"},
			{Field: "AuthorId", Message: "Required"},
			{Field: "Description", Message: "Required"},
			{Field: "Quantity", Message: "Required"},
		}
		expectedResponse, _ := json.Marshal(gin.H{"message": fieldErrors})
		bookRequestJSON, _ := json.Marshal(*invalidBookRequest)
		body := strings.NewReader(string(bookRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/books", body)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return StatusBadRequest when adding book with negative quantity", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBookUsecase := new(mocks.BookUsecase)
		invalidQuantity := -1
		invalidBookRequest := &dto.BookRequest{
			Title:       "Passed",
			AuthorId:    &authorId,
			Description: "Passed",
			Quantity:    &invalidQuantity,
		}
		mockBookUsecase.On("AddBook", ctx, invalidBookRequest).Return(nil, errors.New("error"))
		bookHandler := handler.NewBookHandler(mockBookUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.POST("/books", bookHandler.AddBookHandler)
		fieldErrors := []util.FieldError{
			{Field: "Quantity", Message: "Should be greater than 0"},
		}
		expectedResponse, _ := json.Marshal(gin.H{"message": fieldErrors})
		bookRequestJSON, _ := json.Marshal(*invalidBookRequest)
		body := strings.NewReader(string(bookRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/books", body)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return StatusBadRequest when adding book with existing title", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBookUsecase := new(mocks.BookUsecase)
		invalidBookRequest := &dto.BookRequest{
			Title:       "Duplicate",
			AuthorId:    &authorId,
			Description: "Passed",
			Quantity:    &quantity,
		}
		mockBookUsecase.On("AddBook", ctx, invalidBookRequest).Return(nil, apperror.ErrDuplicateTitle{})
		bookHandler := handler.NewBookHandler(mockBookUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.POST("/books", bookHandler.AddBookHandler)
		fieldErrors := []util.FieldError{
			{Field: "title", Message: "Already existed"},
		}
		expectedResponse, _ := json.Marshal(gin.H{"message": fieldErrors})
		bookRequestJSON, _ := json.Marshal(*invalidBookRequest)
		body := strings.NewReader(string(bookRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/books", body)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return StatusBadRequest when adding book with exceeding title", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBookUsecase := new(mocks.BookUsecase)
		invalidBookRequest := &dto.BookRequest{
			Title:       "Very loooooooooooooooooooooooong title",
			AuthorId:    &authorId,
			Description: "Passed",
			Quantity:    &quantity,
		}
		mockBookUsecase.On("AddBook", ctx, invalidBookRequest).Return(nil, errors.New("error"))
		bookHandler := handler.NewBookHandler(mockBookUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.POST("/books", bookHandler.AddBookHandler)
		fieldErrors := []util.FieldError{
			{Field: "Title", Message: "Should be less than 35 characters"},
		}
		expectedResponse, _ := json.Marshal(gin.H{"message": fieldErrors})
		bookRequestJSON, _ := json.Marshal(*invalidBookRequest)
		body := strings.NewReader(string(bookRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/books", body)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return StatusBadRequest when adding book with nonpositive author id", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBookUsecase := new(mocks.BookUsecase)
		invalidAuthorId := -1
		invalidBookRequest := &dto.BookRequest{
			Title:       "Passed",
			AuthorId:    &invalidAuthorId,
			Description: "Passed",
			Quantity:    &quantity,
		}
		mockBookUsecase.On("AddBook", ctx, invalidBookRequest).Return(nil, errors.New("error"))
		bookHandler := handler.NewBookHandler(mockBookUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.POST("/books", bookHandler.AddBookHandler)
		fieldErrors := []util.FieldError{
			{Field: "AuthorId", Message: "Should be greater than 0"},
		}
		expectedResponse, _ := json.Marshal(gin.H{"message": fieldErrors})
		bookRequestJSON, _ := json.Marshal(*invalidBookRequest)
		body := strings.NewReader(string(bookRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/books", body)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})
}
