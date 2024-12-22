package handler_test

import (
	"archive_lib/apperror"
	"archive_lib/dto"
	"archive_lib/handler"
	"archive_lib/middleware"
	"archive_lib/mocks"
	"archive_lib/util"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	recordId      = 1
	bookId        = 1
	borrowingDate = time.Now()
	returningDate = time.Now()

	borrowRequest = &dto.BorrowRequest{
		BookId: &bookId,
	}

	returnRequest = &dto.ReturnRequest{
		Id: &recordId,
	}

	borrowResponse = &dto.BorrowResponse{
		Id:            recordId,
		UserId:        1,
		BookId:        bookId,
		Status:        "borrowed",
		BorrowingDate: borrowingDate,
	}

	returnResponse = &dto.BorrowResponse{
		Id:            recordId,
		UserId:        1,
		BookId:        bookId,
		Status:        "returned",
		BorrowingDate: borrowingDate,
		ReturningDate: &returningDate,
	}
)

func TestBorrowHandler(t *testing.T) {
	t.Run("should return StatusCreated with borrowed book when no error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "jwt secret for test")
		jwt := util.NewJWT()
		token, _ := jwt.GenerateJWT("1")
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBorrowUsecase := new(mocks.BorrowUsecase)
		mockBorrowUsecase.On("Record", ctx, borrowRequest).Return(borrowResponse, nil)
		borrowHandler := handler.NewBorrowHandler(mockBorrowUsecase)
		router.POST("/borrowing-records", middleware.AuthMiddleware, borrowHandler.BorrowBookHandler)
		expectedResponse, _ := json.Marshal(gin.H{"data": borrowResponse})
		borrowRequestJSON, _ := json.Marshal(*borrowRequest)
		body := strings.NewReader(string(borrowRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/borrowing-records", body)
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return error when get subject from context encounters error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "jwt secret for test")
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBorrowUsecase := new(mocks.BorrowUsecase)
		borrowHandler := handler.NewBorrowHandler(mockBorrowUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.POST("/borrowing-records", middleware.AuthMiddleware, borrowHandler.BorrowBookHandler)
		expectedResponse, _ := json.Marshal(gin.H{"message": "Invalid token"})
		borrowRequestJSON, _ := json.Marshal(*borrowRequest)
		body := strings.NewReader(string(borrowRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/borrowing-records", body)
		ctx.Request.Header.Set("Authorization", "Bearer ")
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return error when user id conversion to string encounters error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "jwt secret for test")
		jwt := util.NewJWT()
		token, _ := jwt.GenerateJWT("")
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBorrowUsecase := new(mocks.BorrowUsecase)
		borrowHandler := handler.NewBorrowHandler(mockBorrowUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.POST("/borrowing-records", middleware.AuthMiddleware, borrowHandler.BorrowBookHandler)
		expectedResponse, _ := json.Marshal(gin.H{"message": "Unrecognized user id"})
		borrowRequestJSON, _ := json.Marshal(*borrowRequest)
		body := strings.NewReader(string(borrowRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/borrowing-records", body)
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return StatusBadRequest when borrow request encounters decode error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "jwt secret for test")
		jwt := util.NewJWT()
		token, _ := jwt.GenerateJWT("1")
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBorrowUsecase := new(mocks.BorrowUsecase)
		mockBorrowUsecase.On("Record", ctx, borrowRequest).Return(borrowResponse, nil)
		borrowHandler := handler.NewBorrowHandler(mockBorrowUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.POST("/borrowing-records", middleware.AuthMiddleware, borrowHandler.BorrowBookHandler)
		errDecodeResponse := []util.FieldError{{Field: "", Message: "Mismatch data type or malformed request"}}
		expectedResponse, _ := json.Marshal(gin.H{"message": errDecodeResponse})
		borrowRequestJSON, _ := json.Marshal("")
		body := strings.NewReader(string(borrowRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/borrowing-records", body)
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return error when record borrow encounters error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "jwt secret for test")
		jwt := util.NewJWT()
		token, _ := jwt.GenerateJWT("0")
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBorrowUsecase := new(mocks.BorrowUsecase)
		mockBorrowUsecase.On("Record", ctx, borrowRequest).Return(nil, apperror.ErrRequestUnrecognized{})
		borrowHandler := handler.NewBorrowHandler(mockBorrowUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.POST("/borrowing-records", middleware.AuthMiddleware, borrowHandler.BorrowBookHandler)
		expectedResponse, _ := json.Marshal(gin.H{"message": "Unrecognized user id"})
		borrowRequestJSON, _ := json.Marshal(*borrowRequest)
		body := strings.NewReader(string(borrowRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/borrowing-records", body)
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})
}

func TestReturnHandler(t *testing.T) {
	t.Run("should return StatusOK with returned book when no error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "jwt secret for test")
		jwt := util.NewJWT()
		token, _ := jwt.GenerateJWT("1")
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBorrowUsecase := new(mocks.BorrowUsecase)
		mockBorrowUsecase.On("Return", ctx, returnRequest).Return(returnResponse, nil)
		borrowHandler := handler.NewBorrowHandler(mockBorrowUsecase)
		router.PATCH("/borrowing-records", middleware.AuthMiddleware, borrowHandler.ReturnBookHandler)
		expectedResponse, _ := json.Marshal(gin.H{"data": returnResponse})
		returnRequestJSON, _ := json.Marshal(*returnRequest)
		body := strings.NewReader(string(returnRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPatch, "/borrowing-records", body)
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return error when get subject from context encounters error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "jwt secret for test")
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBorrowUsecase := new(mocks.BorrowUsecase)
		borrowHandler := handler.NewBorrowHandler(mockBorrowUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.PATCH("/borrowing-records", middleware.AuthMiddleware, borrowHandler.ReturnBookHandler)
		expectedResponse, _ := json.Marshal(gin.H{"message": "Invalid token"})
		returnRequestJSON, _ := json.Marshal(*returnRequest)
		body := strings.NewReader(string(returnRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPatch, "/borrowing-records", body)
		ctx.Request.Header.Set("Authorization", "Bearer ")
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return error when user id conversion to string encounters error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "jwt secret for test")
		jwt := util.NewJWT()
		token, _ := jwt.GenerateJWT("")
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBorrowUsecase := new(mocks.BorrowUsecase)
		borrowHandler := handler.NewBorrowHandler(mockBorrowUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.PATCH("/borrowing-records", middleware.AuthMiddleware, borrowHandler.ReturnBookHandler)
		expectedResponse, _ := json.Marshal(gin.H{"message": "Unrecognized user id"})
		returnRequestJSON, _ := json.Marshal(*returnRequest)
		body := strings.NewReader(string(returnRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPatch, "/borrowing-records", body)
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return StatusBadRequest when return request encounters decode error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "jwt secret for test")
		jwt := util.NewJWT()
		token, _ := jwt.GenerateJWT("1")
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBorrowUsecase := new(mocks.BorrowUsecase)
		mockBorrowUsecase.On("Return", ctx, returnRequest).Return(returnResponse, nil)
		borrowHandler := handler.NewBorrowHandler(mockBorrowUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.PATCH("/borrowing-records", middleware.AuthMiddleware, borrowHandler.ReturnBookHandler)
		errDecodeResponse := []util.FieldError{{Field: "", Message: "Mismatch data type or malformed request"}}
		expectedResponse, _ := json.Marshal(gin.H{"message": errDecodeResponse})
		borrowRequestJSON, _ := json.Marshal("")
		body := strings.NewReader(string(borrowRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPatch, "/borrowing-records", body)
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return error when return borrowed encounters error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "jwt secret for test")
		jwt := util.NewJWT()
		token, _ := jwt.GenerateJWT("0")
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockBorrowUsecase := new(mocks.BorrowUsecase)
		mockBorrowUsecase.On("Return", ctx, returnRequest).Return(nil, apperror.ErrRequestUnrecognized{})
		borrowHandler := handler.NewBorrowHandler(mockBorrowUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.PATCH("/borrowing-records", middleware.AuthMiddleware, borrowHandler.ReturnBookHandler)
		expectedResponse, _ := json.Marshal(gin.H{"message": "Unrecognized user id"})
		returnRequestJSON, _ := json.Marshal(*returnRequest)
		body := strings.NewReader(string(returnRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPatch, "/borrowing-records", body)
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})
}
