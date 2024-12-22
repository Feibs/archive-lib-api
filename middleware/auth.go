package middleware

import (
	"archive_lib/apperror"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	headerSection := strings.Split(header, " ")

	if len(headerSection) != 2 {
		ctx.Error(apperror.ErrLoginFailed{})
		return
	}

	tokenString := headerSection[1]
	token, err := jwt.Parse(
		tokenString,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
		jwt.WithIssuedAt(),
		jwt.WithIssuer("archive_lib"),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		ctx.Error(apperror.ErrInvalidToken{})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.Error(apperror.ErrGetClaimsFailed{})
		return
	}

	subject, err := claims.GetSubject()
	if err != nil {
		ctx.Error(apperror.ErrGetClaimsFailed{})
		return
	}

	ctx.Set("subject", subject)

	ctx.Next()
}
