package usecase_test

import (
	"archive_lib/apperror"
	"archive_lib/dto"
	"archive_lib/entity"
	"archive_lib/mocks"
	"archive_lib/usecase"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	authorId = 5
	quantity = 10
	cover    = "hardcover"

	author = entity.Author{
		Id:   authorId,
		Name: "John The Poet",
	}

	authorResponse = &dto.AuthorResponse{
		Id:   authorId,
		Name: "John The Poet",
	}

	bookPost = entity.BookPost{
		Title:       "Test Book",
		Description: "Cool book",
		Quantity:    quantity,
		Cover:       &cover,
		AuthorId:    authorId,
	}

	book = entity.Book{
		Id:          1,
		Title:       "Test Book",
		Description: "Cool book",
		Quantity:    quantity,
		Cover:       &cover,
	}

	bookWithAuthor = entity.Book{
		Id:          1,
		Author:      &author,
		Title:       "Test Book",
		Description: "Cool book",
		Quantity:    quantity,
		Cover:       &cover,
	}
	booksWithAuthor = []entity.Book{bookWithAuthor}

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

func TestListBooksUsecase(t *testing.T) {
	t.Run("should return list of books when no error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockBookRepo := new(mocks.BookRepo)
		mockBookRepo.On("ListBooks", ctx).Return(booksWithAuthor, nil)
		bookUsecase := usecase.NewBookUsecase(mockBookRepo)

		actualBooksResponse, _ := bookUsecase.ListBooks(ctx)

		assert.Equal(t, booksWithAuthorResponse, actualBooksResponse)
	})

	t.Run("should return error when list books encounters error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockBookRepo := new(mocks.BookRepo)
		mockBookRepo.On("ListBooks", ctx).Return(nil, errors.New("server error"))
		bookUsecase := usecase.NewBookUsecase(mockBookRepo)

		_, err := bookUsecase.ListBooks(ctx)

		assert.NotNil(t, err)
	})
}

func TestGetBooksByTitleUsecase(t *testing.T) {
	t.Run("should return list of books matching the input title when no error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockBookRepo := new(mocks.BookRepo)
		mockBookRepo.On("GetBooksByTitle", ctx, "tes").Return(booksWithAuthor, nil)
		bookUsecase := usecase.NewBookUsecase(mockBookRepo)

		actualBooksResponse, _ := bookUsecase.GetBooksByTitle(ctx, "tes")

		assert.Equal(t, booksWithAuthorResponse, actualBooksResponse)
	})

	t.Run("should return error when search books encounters error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockBookRepo := new(mocks.BookRepo)
		mockBookRepo.On("GetBooksByTitle", ctx, "any").Return(nil, errors.New("server error"))
		bookUsecase := usecase.NewBookUsecase(mockBookRepo)

		_, err := bookUsecase.GetBooksByTitle(ctx, "any")

		assert.NotNil(t, err)
	})
}

func TestAddBookUsecase(t *testing.T) {
	t.Run("should return the added book when no error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockBookRepo := new(mocks.BookRepo)
		mockBookRepo.On("IsTitleExisted", ctx, "Test Book").Return(false, nil)
		mockBookRepo.On("IsAuthorExisted", ctx, authorId).Return(true, nil)
		mockBookRepo.On("AddBook", ctx, &bookPost).Return(&book, nil)
		bookUsecase := usecase.NewBookUsecase(mockBookRepo)

		actualBookResponse, _ := bookUsecase.AddBook(ctx, bookRequest)

		assert.Equal(t, bookResponse, actualBookResponse)
	})

	t.Run("should return ErrDuplicateTitle when adding book with existing title", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockBookRepo := new(mocks.BookRepo)
		mockBookRepo.On("IsTitleExisted", ctx, "Test Book").Return(true, nil)
		bookUsecase := usecase.NewBookUsecase(mockBookRepo)

		_, err := bookUsecase.AddBook(ctx, bookRequest)

		assert.Equal(t, apperror.ErrDuplicateTitle{}, err)
	})

	t.Run("should return ErrAuthorNotFound when adding book with non-existing author", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockBookRepo := new(mocks.BookRepo)
		mockBookRepo.On("IsTitleExisted", ctx, "Test Book").Return(false, nil)
		mockBookRepo.On("IsAuthorExisted", ctx, authorId).Return(false, nil)
		bookUsecase := usecase.NewBookUsecase(mockBookRepo)

		_, err := bookUsecase.AddBook(ctx, bookRequest)

		assert.Equal(t, apperror.ErrAuthorNotFound{}, err)
	})

	t.Run("should return error when the checking of duplicate title encounters error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockBookRepo := new(mocks.BookRepo)
		mockBookRepo.On("IsTitleExisted", ctx, "Test Book").Return(false, errors.New("server error"))
		bookUsecase := usecase.NewBookUsecase(mockBookRepo)

		_, err := bookUsecase.AddBook(ctx, bookRequest)

		assert.NotNil(t, err)
	})

	t.Run("should return error when the checking of author encounters error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockBookRepo := new(mocks.BookRepo)
		mockBookRepo.On("IsTitleExisted", ctx, "Test Book").Return(false, nil)
		mockBookRepo.On("IsAuthorExisted", ctx, authorId).Return(false, errors.New("server error"))
		bookUsecase := usecase.NewBookUsecase(mockBookRepo)

		_, err := bookUsecase.AddBook(ctx, bookRequest)

		assert.NotNil(t, err)
	})

	t.Run("should return error when add book encounters server error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockBookRepo := new(mocks.BookRepo)
		mockBookRepo.On("IsTitleExisted", ctx, "Test Book").Return(false, nil)
		mockBookRepo.On("IsAuthorExisted", ctx, authorId).Return(true, nil)
		mockBookRepo.On("AddBook", ctx, &bookPost).Return(nil, errors.New("server error"))
		bookUsecase := usecase.NewBookUsecase(mockBookRepo)

		_, err := bookUsecase.AddBook(ctx, bookRequest)

		assert.NotNil(t, err)
	})
}
