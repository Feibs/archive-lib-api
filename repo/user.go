package repo

import (
	"archive_lib/entity"
	"context"
	"database/sql"
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
	sql := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);`

	var found bool
	err := repo.db.QueryRowContext(ctx, sql, email).Scan(&found)
	if err != nil {
		return false, err
	}

	return found, nil
}

func (repo userRepoImpl) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	sql := `SELECT id, pass FROM users WHERE email = $1;`

	var user entity.User
	err := repo.db.QueryRowContext(ctx, sql, email).Scan(&user.Id, &user.Password)
	if err != nil {
		return nil, err
	}
	user.Email = email

	return &user, nil
}
