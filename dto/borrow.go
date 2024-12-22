package dto

import "time"

type BorrowResponse struct {
	Id            int        `json:"id"`
	UserId        int        `json:"user_id"`
	BookId        int        `json:"book_id"`
	Status        string     `json:"status"`
	BorrowingDate time.Time  `json:"borrowing_date"`
	ReturningDate *time.Time `json:"returning_date,omitempty"`
}

type BorrowRequest struct {
	BookId *int `json:"book_id" binding:"required,gt=0"`
	UserId int  `json:"user_id"`
}

type ReturnRequest struct {
	Id     *int `json:"id" binding:"required,gt=0"`
	UserId int  `json:"user_id"`
}
