package apperror

type ErrDuplicateTitle struct{}

func (err ErrDuplicateTitle) Error() string {
	return "Already existed"
}

type ErrAuthorNotFound struct{}

func (err ErrAuthorNotFound) Error() string {
	return "Not found"
}

type ErrBookNotFound struct{}

func (err ErrBookNotFound) Error() string {
	return "Book not found"
}

type ErrEmptyStock struct{}

func (err ErrEmptyStock) Error() string {
	return "Empty book stock"
}

type ErrRequestUnrecognized struct{}

func (err ErrRequestUnrecognized) Error() string {
	return "Unrecognized user id"
}

type ErrReturnUnauthorized struct{}

func (err ErrReturnUnauthorized) Error() string {
	return "Unauthorized for this borrowing record"
}

type ErrBorrowNotFound struct{}

func (err ErrBorrowNotFound) Error() string {
	return "Borrowing record not found"
}

type ErrAlreadyReturned struct{}

func (err ErrAlreadyReturned) Error() string {
	return "Book has already been returned"
}

type ErrInvalidToken struct{}

func (err ErrInvalidToken) Error() string {
	return "Invalid token"
}

type ErrGetClaimsFailed struct{}

func (err ErrGetClaimsFailed) Error() string {
	return "Get claims failed"
}

type ErrWrongPassword struct{}

func (err ErrWrongPassword) Error() string {
	return "Password incorrect"
}

type ErrLoginFailed struct{}

func (err ErrLoginFailed) Error() string {
	return "Login failed"
}

type ErrEmailNotFound struct{}

func (err ErrEmailNotFound) Error() string {
	return "Email not registered yet"
}
