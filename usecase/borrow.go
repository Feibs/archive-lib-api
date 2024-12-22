package usecase

import (
	"archive_lib/apperror"
	"archive_lib/dto"
	"archive_lib/entity"
	"archive_lib/repo"
	"context"
)

type BorrowUsecase interface {
	Record(ctx context.Context, borrowRequest *dto.BorrowRequest) (*dto.BorrowResponse, error)
	Return(ctx context.Context, returnRequest *dto.ReturnRequest) (*dto.BorrowResponse, error)
}

type borrowUsecaseImpl struct {
	borrowRepo repo.BorrowRepo
	bookRepo   repo.BookRepo
	txRepo     repo.TransactionRepo
}

func NewBorrowUsecase(borrowRepo repo.BorrowRepo, bookRepo repo.BookRepo, txRepo repo.TransactionRepo) borrowUsecaseImpl {
	return borrowUsecaseImpl{
		borrowRepo: borrowRepo,
		bookRepo:   bookRepo,
		txRepo:     txRepo,
	}
}

func (uc borrowUsecaseImpl) convertBorrowReqToBorrow(dto *dto.BorrowRequest) *entity.Borrow {
	return &entity.Borrow{
		BookId: *dto.BookId,
		UserId: dto.UserId,
	}
}

func (uc borrowUsecaseImpl) convertReturnReqToBorrow(dto *dto.ReturnRequest) *entity.Borrow {
	return &entity.Borrow{
		Id:     *dto.Id,
		UserId: dto.UserId,
	}
}

func (uc borrowUsecaseImpl) convertBorrowToBorrowRes(borrow *entity.Borrow) *dto.BorrowResponse {
	return &dto.BorrowResponse{
		Id:            borrow.Id,
		UserId:        borrow.UserId,
		BookId:        borrow.BookId,
		Status:        borrow.Status,
		BorrowingDate: borrow.BorrowingDate,
	}
}

func (uc borrowUsecaseImpl) convertBorrowToReturnRes(borrow *entity.Borrow) *dto.BorrowResponse {
	return &dto.BorrowResponse{
		Id:            borrow.Id,
		UserId:        borrow.UserId,
		BookId:        borrow.BookId,
		Status:        borrow.Status,
		BorrowingDate: borrow.BorrowingDate,
		ReturningDate: &borrow.ReturningDate,
	}
}

func (uc borrowUsecaseImpl) Record(ctx context.Context, borrowRequest *dto.BorrowRequest) (*dto.BorrowResponse, error) {
	var borrow *entity.Borrow
	var borrowed *entity.Borrow
	bookId := *borrowRequest.BookId

	err := uc.txRepo.WithinTransaction(ctx, func(txCtx context.Context) error {
		found, err := uc.bookRepo.IsBookExisted(txCtx, bookId)
		if err != nil {
			return err
		}
		if !found {
			return apperror.ErrBookNotFound{}
		}

		ok, err := uc.bookRepo.IsStockAvailable(txCtx, bookId)
		if err != nil {
			return err
		}
		if !ok {
			return apperror.ErrEmptyStock{}
		}

		borrow = uc.convertBorrowReqToBorrow(borrowRequest)
		borrowed, err = uc.borrowRepo.Record(txCtx, borrow)
		if err != nil {
			return err
		}
		err = uc.bookRepo.DecrementStock(txCtx, bookId)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return uc.convertBorrowToBorrowRes(borrowed), nil
}

func (uc borrowUsecaseImpl) Return(ctx context.Context, returnRequest *dto.ReturnRequest) (*dto.BorrowResponse, error) {
	var borrow *entity.Borrow
	var borrowed *entity.Borrow
	id := *returnRequest.Id
	user_id := returnRequest.UserId

	err := uc.txRepo.WithinTransaction(ctx, func(txCtx context.Context) error {
		authorized, err := uc.borrowRepo.IsUserAuthorized(txCtx, id, user_id)
		if err != nil {
			return err
		}
		if !authorized {
			return apperror.ErrReturnUnauthorized{}
		}

		found, err := uc.borrowRepo.IsBorrowExisted(txCtx, id)
		if err != nil {
			return err
		}
		if !found {
			return apperror.ErrBorrowNotFound{}
		}

		hasReturned, err := uc.borrowRepo.IsReturned(txCtx, id)
		if err != nil {
			return err
		}
		if hasReturned {
			return apperror.ErrAlreadyReturned{}
		}

		borrow = uc.convertReturnReqToBorrow(returnRequest)
		borrowed, err = uc.borrowRepo.Return(txCtx, borrow)
		if err != nil {
			return err
		}
		bookId, err := uc.borrowRepo.GetBookByBorrowId(txCtx, id)
		if err != nil {
			return err
		}
		err = uc.bookRepo.IncrementStock(txCtx, bookId)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return uc.convertBorrowToReturnRes(borrowed), nil
}
