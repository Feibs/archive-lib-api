package handler_test

import (
	"archive_lib/apperror"
	"archive_lib/dto"
	"archive_lib/handler"
	"archive_lib/middleware"
	"archive_lib/mocks"
	"archive_lib/util"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	authRequest = &dto.AuthRequest{
		Email:    "dokja@mail.com",
		Password: "$2y$10$k/GScY6fCv4kEW5iT6irG.dUBmzUYtscBINENXygnRBxN/ht04bhG",
	}
)

func TestLoginHandler(t *testing.T) {
	t.Run("should return StatusOK with access token when no error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "jwt secret for test")
		jwtImpl := util.NewJWT()
		token, _ := jwtImpl.GenerateJWT("1")
		authResponse := dto.AuthResponse{AccessToken: token}
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockUserUsecase := new(mocks.UserUsecase)
		mockUserUsecase.On("Login", ctx, authRequest).Return(&authResponse, nil)
		userHandler := handler.NewUserHandler(mockUserUsecase)
		router.POST("/login", userHandler.Login)
		authRequestJSON, _ := json.Marshal(*authRequest)
		body := strings.NewReader(string(authRequestJSON))
		expectedResponse, _ := json.Marshal(gin.H{"data": authResponse})

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/login", body)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return StatusBadRequest when auth request encounters decode error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockUserUsecase := new(mocks.UserUsecase)
		userHandler := handler.NewUserHandler(mockUserUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.POST("/login", userHandler.Login)
		errDecodeResponse := []util.FieldError{{Field: "", Message: "Mismatch data type or malformed request"}}
		expectedResponse, _ := json.Marshal(gin.H{"message": errDecodeResponse})
		authRequestJSON, _ := json.Marshal("")
		body := strings.NewReader(string(authRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/login", body)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, string(expectedResponse), w.Body.String())
	})

	t.Run("should return error when login encounters error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "jwt secret for test")
		w := httptest.NewRecorder()
		ctx, router := gin.CreateTestContext(w)
		mockUserUsecase := new(mocks.UserUsecase)
		mockUserUsecase.On("Login", ctx, authRequest).Return(nil, apperror.ErrLoginFailed{})
		userHandler := handler.NewUserHandler(mockUserUsecase)
		router.Use(middleware.ErrorMiddleware)
		router.POST("/login", userHandler.Login)
		authRequestJSON, _ := json.Marshal(authRequest)
		body := strings.NewReader(string(authRequestJSON))

		ctx.Request, _ = http.NewRequest(http.MethodPost, "/login", body)
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
