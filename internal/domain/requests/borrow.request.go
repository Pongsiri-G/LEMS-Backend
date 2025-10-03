package requests

type ReturnRequest struct {
	UserID  string `json:"user_id" validate:"required"`
	StoreID string `json:"store_id" validate:"required"`
}
