package enums

type BorrowStatus string

const (
	StatusBorrowed BorrowStatus = "BORROWED"
	StatusReturned BorrowStatus = "RETURNED"
	StatusWaiting  BorrowStatus = "WAITING"
	StatusCanceled BorrowStatus = "CANCELED"
)
