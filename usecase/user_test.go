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
	authRequest = &dto.AuthRequest{
		Email:    "dokja@mail.com",
		Password: "benchmark",
	}

	authResponse = &dto.AuthResponse{
		AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJsaWJyYXJ5Iiwic3ViIjoiMyIsImV4cCI6MTczMDc4ODM4MSwiaWF0IjoxNzMwNzAxOTgxfQ.aEiNVN7rjmzZEU6qz0ksQ1pRHed8RBAGGGD8zGAYnmM",
	}

	user = &entity.User{
		Id:       1,
		Email:    "dokja@mail.com",
		Password: "$2y$10$pVht.RZGVnCI1CPoSzGPZe2GMSADwiTMftYRK/CUhEsicu/KBhyc.",
	}
)

func TestLoginUsecase(t *testing.T) {
	t.Run("should return access token when no error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockUserRepo := new(mocks.UserRepo)
		mockBcrypt := new(mocks.Bcrypt)
		mockJWT := new(mocks.JWT)
		mockUserRepo.On("IsEmailExisted", ctx, "dokja@mail.com").Return(true, nil)
		mockUserRepo.On("GetUserByEmail", ctx, "dokja@mail.com").Return(user, nil)
		mockBcrypt.On("CompareHashAndPassword", []byte(user.Password), []byte(authRequest.Password)).Return(nil)
		mockJWT.On("GenerateJWT", "1").Return(authResponse.AccessToken, nil)
		userUsecase := usecase.NewUserUsecase(mockUserRepo, mockBcrypt, mockJWT)

		actualAuthResponse, _ := userUsecase.Login(ctx, authRequest)

		assert.Equal(t, authResponse, actualAuthResponse)
	})

	t.Run("should return ErrEmailNotFound when login with non-existing email", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockUserRepo := new(mocks.UserRepo)
		mockBcrypt := new(mocks.Bcrypt)
		mockJWT := new(mocks.JWT)
		mockUserRepo.On("IsEmailExisted", ctx, "dokja@mail.com").Return(false, nil)
		userUsecase := usecase.NewUserUsecase(mockUserRepo, mockBcrypt, mockJWT)

		_, err := userUsecase.Login(ctx, authRequest)

		assert.Equal(t, apperror.ErrEmailNotFound{}, err)
	})

	t.Run("should return error when email checking encounters error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockUserRepo := new(mocks.UserRepo)
		mockBcrypt := new(mocks.Bcrypt)
		mockJWT := new(mocks.JWT)
		mockUserRepo.On("IsEmailExisted", ctx, "dokja@mail.com").Return(false, errors.New("error"))
		userUsecase := usecase.NewUserUsecase(mockUserRepo, mockBcrypt, mockJWT)

		_, err := userUsecase.Login(ctx, authRequest)

		assert.NotNil(t, err)
	})

	t.Run("should return error when get user by email encounters error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockUserRepo := new(mocks.UserRepo)
		mockBcrypt := new(mocks.Bcrypt)
		mockJWT := new(mocks.JWT)
		mockUserRepo.On("IsEmailExisted", ctx, "dokja@mail.com").Return(true, nil)
		mockUserRepo.On("GetUserByEmail", ctx, "dokja@mail.com").Return(nil, errors.New("error"))
		userUsecase := usecase.NewUserUsecase(mockUserRepo, mockBcrypt, mockJWT)

		_, err := userUsecase.Login(ctx, authRequest)

		assert.NotNil(t, err)
	})

	t.Run("should return ErrWrongPassword when login with wrong password", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockUserRepo := new(mocks.UserRepo)
		mockBcrypt := new(mocks.Bcrypt)
		mockJWT := new(mocks.JWT)
		mockUserRepo.On("IsEmailExisted", ctx, "dokja@mail.com").Return(true, nil)
		mockUserRepo.On("GetUserByEmail", ctx, "dokja@mail.com").Return(user, nil)
		mockBcrypt.On("CompareHashAndPassword", []byte(user.Password), []byte(authRequest.Password)).Return(apperror.ErrWrongPassword{})
		userUsecase := usecase.NewUserUsecase(mockUserRepo, mockBcrypt, mockJWT)

		_, err := userUsecase.Login(ctx, authRequest)

		assert.Equal(t, apperror.ErrWrongPassword{}, err)
	})

	t.Run("should return ErrLoginFailed when generate jwt encounters error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockUserRepo := new(mocks.UserRepo)
		mockBcrypt := new(mocks.Bcrypt)
		mockJWT := new(mocks.JWT)
		mockUserRepo.On("IsEmailExisted", ctx, "dokja@mail.com").Return(true, nil)
		mockUserRepo.On("GetUserByEmail", ctx, "dokja@mail.com").Return(user, nil)
		mockBcrypt.On("CompareHashAndPassword", []byte(user.Password), []byte(authRequest.Password)).Return(nil)
		mockJWT.On("GenerateJWT", "1").Return("", apperror.ErrLoginFailed{})
		userUsecase := usecase.NewUserUsecase(mockUserRepo, mockBcrypt, mockJWT)

		_, err := userUsecase.Login(ctx, authRequest)

		assert.Equal(t, apperror.ErrLoginFailed{}, err)
	})
}
