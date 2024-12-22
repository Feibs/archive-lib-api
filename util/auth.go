package util

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt interface {
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

type JWT interface {
	GenerateJWT(userId string) (string, error)
}

type bcryptImpl struct{}

type jwtImpl struct{}

func NewBcrypt() bcryptImpl {
	return bcryptImpl{}
}

func NewJWT() jwtImpl {
	return jwtImpl{}
}

func (b bcryptImpl) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func (j jwtImpl) GenerateJWT(userId string) (string, error) {
	now := time.Now()
	registeredClaims := jwt.RegisteredClaims{
		Issuer: "archive_lib",
		IssuedAt: &jwt.NumericDate{
			Time: now,
		},
		Subject: userId,
		ExpiresAt: &jwt.NumericDate{
			Time: now.Add(24 * time.Hour),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
