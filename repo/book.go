package repo

import (
	"archive_lib/entity"
	"context"
	"database/sql"
	"strings"
)

type BookRepo interface {
	ListBooks(ctx context.Context) ([]entity.Book, error)
	GetBooksByTitle(ctx context.Context, title string) ([]entity.Book, error)
	AddBook(ctx context.Context, bookPost *entity.BookPost) (*entity.Book, error)
	IsTitleExisted(ctx context.Context, title string) (bool, error)
	IsAuthorExisted(ctx context.Context, id int) (bool, error)
	IsBookExisted(ctx context.Context, id int) (bool, error)
	IsStockAvailable(ctx context.Context, id int) (bool, error)
	DecrementStock(ctx context.Context, id int) error
	IncrementStock(ctx context.Context, id int) error
}

type bookRepoImpl struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) bookRepoImpl {
	return bookRepoImpl{
		db: db,
	}
}

func (repo bookRepoImpl) ListBooks(ctx context.Context) ([]entity.Book, error) {
	books := []entity.Book{}
	sql := `SELECT 
				b.id, b.author_id, a.name, b.title, b.description, b.quantity, b.cover
			FROM 
				books b JOIN authors a ON a.id = b.author_id 
			WHERE 
				b.deleted_at IS NULL;`

	rows, err := repo.db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var book entity.Book
		var author entity.Author
		err := rows.Scan(
			&book.Id,
			&author.Id,
			&author.Name,
			&book.Title,
			&book.Description,
			&book.Quantity,
			&book.Cover,
		)
		if err != nil {
			return nil, err
		}
		book.Author = &author
		books = append(books, book)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (repo bookRepoImpl) GetBooksByTitle(ctx context.Context, title string) ([]entity.Book, error) {
	books := []entity.Book{}

	sql := `SELECT 
				b.id, b.author_id, a.name, b.title, b.description, b.quantity, b.cover
			FROM 
				books b JOIN authors a ON a.id = b.author_id 
			WHERE 
				b.title ILIKE $1 AND b.deleted_at IS NULL`

	rows, err := repo.db.QueryContext(ctx, sql, "%"+title+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var book entity.Book
		var author entity.Author
		err := rows.Scan(
			&book.Id,
			&author.Id,
			&author.Name,
			&book.Title,
			&book.Description,
			&book.Quantity,
			&book.Cover,
		)
		if err != nil {
			return nil, err
		}
		book.Author = &author
		books = append(books, book)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (repo bookRepoImpl) IsTitleExisted(ctx context.Context, title string) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM books WHERE title = $1);`

	var found bool
	err := repo.db.QueryRowContext(ctx, sql, title).Scan(&found)
	if err != nil {
		return false, err
	}

	return found, nil
}

func (repo bookRepoImpl) buildAddBookQuery(book *entity.BookPost, inputs *[]any) string {
	var query strings.Builder

	columns := `INSERT INTO books (title, author_id, description, quantity`
	values := `) VALUES ($1, $2, $3, $4`
	returning := `) RETURNING id`

	if book.Cover != nil {
		columns += ", cover"
		values += ", $5"
		*inputs = append(*inputs, *book.Cover)
	}

	query.WriteString(columns)
	query.WriteString(values)
	query.WriteString(returning)

	return query.String()
}

func (repo bookRepoImpl) IsAuthorExisted(ctx context.Context, id int) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM authors WHERE id = $1);`

	var found bool
	err := repo.db.QueryRowContext(ctx, sql, id).Scan(&found)
	if err != nil {
		return false, err
	}

	return found, nil
}

func (repo bookRepoImpl) AddBook(ctx context.Context, bookPost *entity.BookPost) (*entity.Book, error) {
	inputs := make([]any, 0)
	inputs = append(inputs, bookPost.Title, bookPost.AuthorId, bookPost.Description, bookPost.Quantity)

	sql := repo.buildAddBookQuery(bookPost, &inputs)
	err := repo.db.QueryRowContext(ctx, sql, inputs...).Scan(&bookPost.Id)
	if err != nil {
		return nil, err
	}

	return bookPost.ConvertToBook(), nil
}

func (repo bookRepoImpl) IsBookExisted(ctx context.Context, id int) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM books WHERE id = $1) FOR UPDATE;`

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

func (repo bookRepoImpl) IsStockAvailable(ctx context.Context, id int) (bool, error) {
	sql := `SELECT quantity FROM books WHERE id = $1;`

	tx := extractTx(ctx)
	var err error
	var stock int

	if tx != nil {
		err = tx.QueryRowContext(ctx, sql, id).Scan(&stock)
	} else {
		err = repo.db.QueryRowContext(ctx, sql, id).Scan(&stock)
	}

	if err != nil {
		return false, err
	}

	return (stock > 0), nil
}

func (repo bookRepoImpl) DecrementStock(ctx context.Context, id int) error {
	sql := `UPDATE 
				books
			SET 
				quantity = quantity - 1, 
				updated_at = NOW()
			WHERE 
				id=$1`

	tx := extractTx(ctx)
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, sql, id)
	} else {
		_, err = repo.db.ExecContext(ctx, sql, id)
	}

	if err != nil {
		return err
	}

	return nil
}

func (repo bookRepoImpl) IncrementStock(ctx context.Context, id int) error {
	sql := `UPDATE 
				books
			SET 
				quantity = quantity + 1, 
				updated_at = NOW()
			WHERE 
				id=$1`

	tx := extractTx(ctx)
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, sql, id)
	} else {
		_, err = repo.db.ExecContext(ctx, sql, id)
	}

	if err != nil {
		return err
	}

	return nil
}
