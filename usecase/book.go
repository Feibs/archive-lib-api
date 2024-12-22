package usecase

import (
	"archive_lib/apperror"
	"archive_lib/dto"
	"archive_lib/entity"
	"archive_lib/repo"
	"context"
)

type BookUsecase interface {
	ListBooks(ctx context.Context) ([]*dto.BookResponse, error)
	GetBooksByTitle(ctx context.Context, title string) ([]*dto.BookResponse, error)
	AddBook(ctx context.Context, bookRequest *dto.BookRequest) (*dto.BookResponse, error)
}

type bookUsecaseImpl struct {
	bookRepo repo.BookRepo
}

func NewBookUsecase(br repo.BookRepo) bookUsecaseImpl {
	return bookUsecaseImpl{
		bookRepo: br,
	}
}

func (uc bookUsecaseImpl) convertBookToRes(book *entity.Book) *dto.BookResponse {
	var cover string
	if book.Cover != nil {
		cover = *book.Cover
	}

	if book.Author != nil {
		return &dto.BookResponse{
			Id: book.Id,
			Author: &dto.AuthorResponse{
				Id:   book.Author.Id,
				Name: book.Author.Name,
			},
			Title:       book.Title,
			Description: book.Description,
			Quantity:    book.Quantity,
			Cover:       cover,
		}
	}

	return &dto.BookResponse{
		Id:          book.Id,
		Title:       book.Title,
		Description: book.Description,
		Quantity:    book.Quantity,
		Cover:       cover,
	}
}

func (uc bookUsecaseImpl) convertReqToBookPost(dto *dto.BookRequest) *entity.BookPost {
	return &entity.BookPost{
		Title:       dto.Title,
		Description: dto.Description,
		Quantity:    *dto.Quantity,
		Cover:       dto.Cover,
		AuthorId:    *dto.AuthorId,
	}
}

func (uc bookUsecaseImpl) ListBooks(ctx context.Context) ([]*dto.BookResponse, error) {
	books, err := uc.bookRepo.ListBooks(ctx)
	if err != nil {
		return nil, err
	}
	booksResponse := []*dto.BookResponse{}
	for _, book := range books {
		booksResponse = append(booksResponse, uc.convertBookToRes(&book))
	}
	return booksResponse, nil
}

func (uc bookUsecaseImpl) GetBooksByTitle(ctx context.Context, title string) ([]*dto.BookResponse, error) {
	books, err := uc.bookRepo.GetBooksByTitle(ctx, title)
	if err != nil {
		return nil, err
	}
	booksResponse := []*dto.BookResponse{}
	for _, book := range books {
		booksResponse = append(booksResponse, uc.convertBookToRes(&book))
	}
	return booksResponse, nil
}

func (uc bookUsecaseImpl) AddBook(ctx context.Context, bookRequest *dto.BookRequest) (*dto.BookResponse, error) {
	duplicate, err := uc.bookRepo.IsTitleExisted(ctx, bookRequest.Title)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, apperror.ErrDuplicateTitle{}
	}

	found, err := uc.bookRepo.IsAuthorExisted(ctx, *bookRequest.AuthorId)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, apperror.ErrAuthorNotFound{}
	}

	addedBook, err := uc.bookRepo.AddBook(ctx, uc.convertReqToBookPost(bookRequest))
	if err != nil {
		return nil, err
	}
	return uc.convertBookToRes(addedBook), nil
}
