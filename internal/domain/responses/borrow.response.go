package responses

import "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"

type UserBorrrowResponse struct {
	BorrowID     string             `json:"borrow_id"`
	ItemName     string             `json:"item_name"`
	BorrowDate   string             `json:"borrow_date"`
	ReturnDate   *string            `json:"return_date"`
	BorrowStatus enums.BorrowStatus `json:"borrow_status"`
}

type AdminBorrowResponse struct {
	BorrowID       string             `json:"borrow_id"`
	BorrowParentID *string            `json:"borrow_parent_id"`
	UserID         string             `json:"user_id"`
	UserName       string             `json:"user_name"`
	ItemID         string             `json:"item_id"`
	ItemName       string             `json:"item_name"`
	BorrowDate     string             `json:"borrow_date"`
	ReturnURL      *string            `json:"return_image_url"`
	ReturnDate     *string            `json:"return_date"`
	BorrowStatus   enums.BorrowStatus `json:"borrow_status"`
}
