package usecase

import (
	"archive_lib/apperror"
	"archive_lib/dto"
	"archive_lib/repo"
	"archive_lib/util"
	"context"
	"strconv"
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
	email := authRequest.Email

	found, err := uc.userRepo.IsEmailExisted(ctx, email)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, apperror.ErrEmailNotFound{}
	}

	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = uc.bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authRequest.Password))
	if err != nil {
		return nil, apperror.ErrWrongPassword{}
	}

	userId := strconv.Itoa(user.Id)
	token, err := uc.jwt.GenerateJWT(userId)
	if err != nil {
		return nil, apperror.ErrLoginFailed{}
	}

	return &dto.AuthResponse{AccessToken: token}, nil
}
