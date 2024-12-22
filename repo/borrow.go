package repo

import (
	"archive_lib/entity"
	"context"
	"database/sql"
)

type BorrowRepo interface {
	Record(ctx context.Context, borrow *entity.Borrow) (*entity.Borrow, error)
	Return(ctx context.Context, borrow *entity.Borrow) (*entity.Borrow, error)
	GetBookByBorrowId(ctx context.Context, id int) (int, error)
	IsUserAuthorized(ctx context.Context, record_id int, user_id int) (bool, error)
	IsBorrowExisted(ctx context.Context, id int) (bool, error)
	IsReturned(ctx context.Context, id int) (bool, error)
}

type borrowRepoImpl struct {
	db *sql.DB
}

func NewBorrowRepo(db *sql.DB) borrowRepoImpl {
	return borrowRepoImpl{
		db: db,
	}
}

func (repo borrowRepoImpl) Record(ctx context.Context, borrow *entity.Borrow) (*entity.Borrow, error) {
	const status = "borrowed"
	sql := `INSERT INTO 
				borrowing_records (user_id, book_id, status, borrowing_date) 
			VALUES 
				($1, $2, $3, NOW()) 
			RETURNING 
				id, status, borrowing_date`

	tx := extractTx(ctx)
	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, sql, borrow.UserId, borrow.BookId, status).Scan(&borrow.Id, &borrow.Status, &borrow.BorrowingDate)
	} else {
		err = repo.db.QueryRowContext(ctx, sql, borrow.UserId, borrow.BookId, status).Scan(&borrow.Id, &borrow.Status, &borrow.BorrowingDate)
	}

	if err != nil {
		return nil, err
	}

	return borrow, nil
}

func (repo borrowRepoImpl) IsBorrowExisted(ctx context.Context, id int) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM borrowing_records WHERE id = $1);`

	tx := extractTx(ctx)
	var err error
	var found bool

	if tx != nil {
		err = tx.QueryRowContext(ctx, sql, id).Scan(&found)
	} else {
		err = repo.db.QueryRowContext(ctx, sql, id).Scan(&found)
	}

	if err != nil {
		return false, err
	}

	return found, nil
}

func (repo borrowRepoImpl) GetBookByBorrowId(ctx context.Context, id int) (int, error) {
	sql := `SELECT book_id FROM borrowing_records WHERE id = $1;`

	tx := extractTx(ctx)
	var err error
	var bookId int

	if tx != nil {
		err = tx.QueryRowContext(ctx, sql, id).Scan(&bookId)
	} else {
		err = repo.db.QueryRowContext(ctx, sql, id).Scan(&bookId)
	}

	if err != nil {
		return bookId, err
	}

	return bookId, nil
}

func (repo borrowRepoImpl) IsUserAuthorized(ctx context.Context, record_id int, user_id int) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM borrowing_records WHERE id = $1 AND user_id = $2);`

	tx := extractTx(ctx)
	var err error
	var found bool

	if tx != nil {
		err = tx.QueryRowContext(ctx, sql, record_id, user_id).Scan(&found)
	} else {
		err = repo.db.QueryRowContext(ctx, sql, record_id, user_id).Scan(&found)
	}

	if err != nil {
		return false, err
	}

	return found, nil
}

func (repo borrowRepoImpl) IsReturned(ctx context.Context, id int) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM borrowing_records WHERE id = $1 AND returning_date IS NOT NULL);`

	tx := extractTx(ctx)
	var err error
	var found bool

	if tx != nil {
		err = tx.QueryRowContext(ctx, sql, id).Scan(&found)
	} else {
		err = repo.db.QueryRowContext(ctx, sql, id).Scan(&found)
	}

	if err != nil {
		return false, err
	}

	return found, nil
}

func (repo borrowRepoImpl) Return(ctx context.Context, borrow *entity.Borrow) (*entity.Borrow, error) {
	const status = "returned"
	sql := `UPDATE 
				borrowing_records 
			SET 
				status = $3, 
				returning_date = NOW(), 
				updated_at = NOW() 
			WHERE 
				id = $1 AND user_id = $2
			RETURNING 
				book_id, status, borrowing_date, returning_date;`

	tx := extractTx(ctx)
	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, sql, borrow.Id, borrow.UserId, status).Scan(
			&borrow.BookId,
			&borrow.Status,
			&borrow.BorrowingDate,
			&borrow.ReturningDate,
		)
	} else {
		err = repo.db.QueryRowContext(ctx, sql, borrow.Id, borrow.UserId, status).Scan(
			&borrow.BookId,
			&borrow.Status,
			&borrow.BorrowingDate,
			&borrow.ReturningDate,
		)
	}

	if err != nil {
		return nil, err
	}

	return borrow, nil
}
