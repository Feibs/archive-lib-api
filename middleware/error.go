package middleware

import (
	"archive_lib/apperror"
	"archive_lib/util"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorMiddleware(ctx *gin.Context) {
	ctx.Next()

	var fieldErrors []util.FieldError

	if len(ctx.Errors) > 0 {
		var ve validator.ValidationErrors
		err := ctx.Errors[0]
		if errors.As(err, &ve) {
			for _, fe := range ve {
				fieldErrors = append(fieldErrors, util.FieldError{
					Field:   fe.Field(),
					Message: util.ExtractValidationError(fe),
				})
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": fieldErrors})
			return
		}

		var je *json.UnmarshalTypeError
		if errors.As(err, &je) {
			fieldErrors = append(fieldErrors, util.FieldError{
				Field:   je.Field,
				Message: util.ExtractUnmarshalError(je),
			})
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": fieldErrors})
			return
		}

		var errAuthorNotFound apperror.ErrAuthorNotFound
		if errors.As(err, &errAuthorNotFound) {
			fieldErrors = append(fieldErrors, util.FieldError{
				Field:   "author_id",
				Message: err.Error(),
			})
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": fieldErrors})
			return
		}

		var errBookNotFound apperror.ErrBookNotFound
		if errors.As(err, &errBookNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		var errDuplicateTitle apperror.ErrDuplicateTitle
		if errors.As(err, &errDuplicateTitle) {
			fieldErrors = append(fieldErrors, util.FieldError{
				Field:   "title",
				Message: err.Error(),
			})
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": fieldErrors})
			return
		}

		var errEmailNotFound apperror.ErrEmailNotFound
		if errors.As(err, &errEmailNotFound) {
			fieldErrors = append(fieldErrors, util.FieldError{
				Field:   "email",
				Message: err.Error(),
			})
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": fieldErrors})
			return
		}

		var errWrongPassword apperror.ErrWrongPassword
		if errors.As(err, &errWrongPassword) {
			fieldErrors = append(fieldErrors, util.FieldError{
				Field:   "password",
				Message: err.Error(),
			})
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": fieldErrors})
			return
		}

		var errInvalidToken apperror.ErrInvalidToken
		if errors.As(err, &errInvalidToken) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		var errGetClaimsFailed apperror.ErrGetClaimsFailed
		if errors.As(err, &errGetClaimsFailed) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		var errLoginFailed apperror.ErrLoginFailed
		if errors.As(err, &errLoginFailed) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		var errEmptyStock apperror.ErrEmptyStock
		if errors.As(err, &errEmptyStock) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var errRequestUnrecognized apperror.ErrRequestUnrecognized
		if errors.As(err, &errRequestUnrecognized) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		var errReturnUnauthorized apperror.ErrReturnUnauthorized
		if errors.As(err, &errReturnUnauthorized) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		var errBorrowNotFound apperror.ErrBorrowNotFound
		if errors.As(err, &errBorrowNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		var errAlreadyReturned apperror.ErrAlreadyReturned
		if errors.As(err, &errAlreadyReturned) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
}
