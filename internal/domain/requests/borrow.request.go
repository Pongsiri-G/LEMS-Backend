package requests

type BorrowRequest struct {
	UserID string `json:"user_id" validate:"required"`
	ItemID string `json:"item_id" validate:"required"`
}

type ReturnRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	BorrowID string `json:"borrow_id" validate:"required"`
}
