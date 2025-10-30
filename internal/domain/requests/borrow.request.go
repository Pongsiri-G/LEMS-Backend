package requests

type BorrowRequest struct {
	ItemID string `json:"item_id" validate:"required"`
}

type ReturnRequest struct {
	BorrowID string `json:"borrow_id" validate:"required"`
	ImageURL string `json:"image_url" validate:"required"`
}
