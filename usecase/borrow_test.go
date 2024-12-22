package usecase_test

import (
	"archive_lib/apperror"
	"archive_lib/dto"
	"archive_lib/entity"
	"archive_lib/mocks"
	"archive_lib/usecase"
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	recordId      = 1
	bookId        = 1
	borrowingDate = time.Now()

	borrowRequest = &dto.BorrowRequest{
		BookId: &bookId,
		UserId: 1,
	}

	borrow = &entity.Borrow{
		UserId: 1,
		BookId: bookId,
	}

	borrowed = &entity.Borrow{
		Id:            recordId,
		UserId:        1,
		BookId:        bookId,
		Status:        "borrowed",
		BorrowingDate: borrowingDate,
	}

	borrowResponse = &dto.BorrowResponse{
		Id:            recordId,
		UserId:        1,
		BookId:        bookId,
		Status:        "borrowed",
		BorrowingDate: borrowingDate,
	}
)

func TestRecordBorrowUsecase(t *testing.T) {
	t.Run("should return borrowed book when there is no error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockBookRepo := new(mocks.BookRepo)
		mockBorrowRepo := new(mocks.BorrowRepo)
		mockTxRepo := new(mocks.TransactionRepo)
		mockTxRepo.On(
			"WithinTransaction",
			ctx,
			mock.MatchedBy(func(txFn func(context.Context) error) bool {
				err := txFn(ctx)
				return err == nil
			}),
		).Return(nil)
		mockBookRepo.On("IsBookExisted", ctx, bookId).Return(true, nil)
		mockBookRepo.On("IsStockAvailable", ctx, bookId).Return(true, nil)
		mockBorrowRepo.On("Record", ctx, borrow).Return(borrowed, nil)
		mockBookRepo.On("DecrementStock", ctx, bookId).Return(nil)
		borrowUsecase := usecase.NewBorrowUsecase(mockBorrowRepo, mockBookRepo, mockTxRepo)

		borrowRecord, _ := borrowUsecase.Record(ctx, borrowRequest)

		assert.Equal(t, borrowResponse, borrowRecord)
	})

	t.Run("should return error when book checking encounters error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		errIsBookExisted := errors.New("error")
		mockBookRepo := new(mocks.BookRepo)
		mockBorrowRepo := new(mocks.BorrowRepo)
		mockTxRepo := new(mocks.TransactionRepo)
		mockTxRepo.On(
			"WithinTransaction",
			ctx,
			mock.MatchedBy(func(txFn func(context.Context) error) bool {
				err := txFn(ctx)
				return err == errIsBookExisted
			}),
		).Return(errIsBookExisted)
		mockBookRepo.On("IsBookExisted", ctx, bookId).Return(false, errIsBookExisted)
		borrowUsecase := usecase.NewBorrowUsecase(mockBorrowRepo, mockBookRepo, mockTxRepo)

		_, err := borrowUsecase.Record(ctx, borrowRequest)

		assert.Equal(t, err, errIsBookExisted)
	})

	t.Run("should return error when borrowing non-existent book", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		errBookNotFound := apperror.ErrBookNotFound{}
		mockBookRepo := new(mocks.BookRepo)
		mockBorrowRepo := new(mocks.BorrowRepo)
		mockTxRepo := new(mocks.TransactionRepo)
		mockTxRepo.On(
			"WithinTransaction",
			ctx,
			mock.MatchedBy(func(txFn func(context.Context) error) bool {
				err := txFn(ctx)
				return err == errBookNotFound
			}),
		).Return(errBookNotFound)
		mockBookRepo.On("IsBookExisted", ctx, 1).Return(false, nil)
		borrowUsecase := usecase.NewBorrowUsecase(mockBorrowRepo, mockBookRepo, mockTxRepo)

		_, err := borrowUsecase.Record(ctx, borrowRequest)

		assert.Equal(t, err, errBookNotFound)
	})

	t.Run("should return error when stock checking encounters error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		errIsStockAvailable := errors.New("error")
		mockBookRepo := new(mocks.BookRepo)
		mockBorrowRepo := new(mocks.BorrowRepo)
		mockTxRepo := new(mocks.TransactionRepo)
		mockTxRepo.On(
			"WithinTransaction",
			ctx,
			mock.MatchedBy(func(txFn func(context.Context) error) bool {
				err := txFn(ctx)
				return err == errIsStockAvailable
			}),
		).Return(errIsStockAvailable)
		mockBookRepo.On("IsBookExisted", ctx, 1).Return(true, nil)
		mockBookRepo.On("IsStockAvailable", ctx, bookId).Return(false, errIsStockAvailable)
		borrowUsecase := usecase.NewBorrowUsecase(mockBorrowRepo, mockBookRepo, mockTxRepo)

		_, err := borrowUsecase.Record(ctx, borrowRequest)

		assert.Equal(t, err, errIsStockAvailable)
	})

	t.Run("should return error when borrowing book with empty stock", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		errEmptyStock := apperror.ErrEmptyStock{}
		mockBookRepo := new(mocks.BookRepo)
		mockBorrowRepo := new(mocks.BorrowRepo)
		mockTxRepo := new(mocks.TransactionRepo)
		mockTxRepo.On(
			"WithinTransaction",
			ctx,
			mock.MatchedBy(func(txFn func(context.Context) error) bool {
				err := txFn(ctx)
				return err == errEmptyStock
			}),
		).Return(errEmptyStock)
		mockBookRepo.On("IsBookExisted", ctx, 1).Return(true, nil)
		mockBookRepo.On("IsStockAvailable", ctx, bookId).Return(false, nil)
		borrowUsecase := usecase.NewBorrowUsecase(mockBorrowRepo, mockBookRepo, mockTxRepo)

		_, err := borrowUsecase.Record(ctx, borrowRequest)

		assert.Equal(t, err, errEmptyStock)
	})

	t.Run("should return error when record borrow encounters error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		errRecord := errors.New("error")
		mockBookRepo := new(mocks.BookRepo)
		mockBorrowRepo := new(mocks.BorrowRepo)
		mockTxRepo := new(mocks.TransactionRepo)
		mockTxRepo.On(
			"WithinTransaction",
			ctx,
			mock.MatchedBy(func(txFn func(context.Context) error) bool {
				err := txFn(ctx)
				return err == errRecord
			}),
		).Return(errRecord)
		mockBookRepo.On("IsBookExisted", ctx, 1).Return(true, nil)
		mockBookRepo.On("IsStockAvailable", ctx, bookId).Return(true, nil)
		mockBorrowRepo.On("Record", ctx, borrow).Return(nil, errRecord)
		borrowUsecase := usecase.NewBorrowUsecase(mockBorrowRepo, mockBookRepo, mockTxRepo)

		_, err := borrowUsecase.Record(ctx, borrowRequest)

		assert.Equal(t, err, errRecord)
	})

	t.Run("should return error when decrement stock of the borrowed book encounters error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		errDecrementStock := errors.New("error")
		mockBookRepo := new(mocks.BookRepo)
		mockBorrowRepo := new(mocks.BorrowRepo)
		mockTxRepo := new(mocks.TransactionRepo)
		mockTxRepo.On(
			"WithinTransaction",
			ctx,
			mock.MatchedBy(func(txFn func(context.Context) error) bool {
				err := txFn(ctx)
				return err == errDecrementStock
			}),
		).Return(errDecrementStock)
		mockBookRepo.On("IsBookExisted", ctx, 1).Return(true, nil)
		mockBookRepo.On("IsStockAvailable", ctx, bookId).Return(true, nil)
		mockBorrowRepo.On("Record", ctx, borrow).Return(borrowed, nil)
		mockBookRepo.On("DecrementStock", ctx, bookId).Return(errDecrementStock)
		borrowUsecase := usecase.NewBorrowUsecase(mockBorrowRepo, mockBookRepo, mockTxRepo)

		_, err := borrowUsecase.Record(ctx, borrowRequest)

		assert.Equal(t, err, errDecrementStock)
	})
}
