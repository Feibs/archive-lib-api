package usecase

import (
	"archive_lib/apperror"
	"archive_lib/dto"
	"archive_lib/repo"
	"archive_lib/util"
	"context"
	"strconv"

	"github.com/opentracing/opentracing-go"
)

type UserUsecase interface {
	Login(ctx context.Context, authRequest *dto.AuthRequest) (*dto.AuthResponse, error)
}

type userUsecaseImpl struct {
	userRepo repo.UserRepo
	bcrypt   util.Bcrypt
	jwt      util.JWT
}

func NewUserUsecase(br repo.UserRepo, bcrypt util.Bcrypt, jwt util.JWT) userUsecaseImpl {
	return userUsecaseImpl{
		userRepo: br,
		bcrypt:   bcrypt,
		jwt:      jwt,
	}
}

func (uc userUsecaseImpl) Login(ctx context.Context, authRequest *dto.AuthRequest) (*dto.AuthResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "archive-lib-service usecase Login func")
	defer span.Finish()
	ctxTracer := opentracing.ContextWithSpan(context.Background(), span)

	email := authRequest.Email

	found, err := uc.userRepo.IsEmailExisted(ctxTracer, email)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("error", err)
		return nil, err
	}

	if !found {
		span.SetTag("error", true)
		span.LogKV("error", err)
		return nil, apperror.ErrEmailNotFound{}
	}

	span.LogKV("info", "usecase Login IsEmailExisted success")

	user, err := uc.userRepo.GetUserByEmail(ctxTracer, email)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("error", err)
		return nil, err
	}

	span.LogKV("info", "usecase Login GetUserByEmail success")

	err = uc.bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authRequest.Password))
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("error", err)
		return nil, apperror.ErrWrongPassword{}
	}

	userId := strconv.Itoa(user.Id)
	token, err := uc.jwt.GenerateJWT(userId)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("error", err)
		return nil, apperror.ErrLoginFailed{}
	}

	span.LogKV("info", "usecase Login success")

	return &dto.AuthResponse{AccessToken: token}, nil
}
