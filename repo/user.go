package repo

import (
	"archive_lib/entity"
	"context"
	"database/sql"

	"github.com/opentracing/opentracing-go"
)

type UserRepo interface {
	IsEmailExisted(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

type userRepoImpl struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) userRepoImpl {
	return userRepoImpl{
		db: db,
	}
}

func (repo userRepoImpl) IsEmailExisted(ctx context.Context, email string) (bool, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "archive-lib-service repo Login func")
	defer span.Finish()
	ctxTracer := opentracing.ContextWithSpan(context.Background(), span)

	sql := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);`

	var found bool
	err := repo.db.QueryRowContext(ctxTracer, sql, email).Scan(&found)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("error", err)
		return false, err
	}

	span.LogKV("info", "repo Login IsEmailExisted success")

	return found, nil
}

func (repo userRepoImpl) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "archive-lib-service repo Login func")
	defer span.Finish()
	ctxTracer := opentracing.ContextWithSpan(context.Background(), span)

	sql := `SELECT id, pass FROM users WHERE email = $1;`

	var user entity.User
	err := repo.db.QueryRowContext(ctxTracer, sql, email).Scan(&user.Id, &user.Password)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("error", err)
		return nil, err
	}
	user.Email = email

	span.LogKV("info", "repo Login GetUserByEmail success")

	return &user, nil
}
