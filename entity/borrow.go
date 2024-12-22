package entity

import "time"

type Borrow struct {
	Id            int
	UserId        int
	BookId        int
	Status        string
	BorrowingDate time.Time
	ReturningDate time.Time
}
