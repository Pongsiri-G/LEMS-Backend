package responses

import "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"

type UserBorrrowResponse struct {
	BorrowID     string             `json:"borrow_id"`
	ItemName     string             `json:"item_name"`
	BorrowDate   string             `json:"borrow_date"`
	ReturnDate   *string            `json:"return_date"`
	BorrowStatus enums.BorrowStatus `json:"borrow_status"`
}
