package responses

import "time"

type BorrowQueueResponse struct {
	QueueID    string     `json:"queue_id"`
	UserID     string     `json:"user_id"`
	ItemID     string     `json:"item_id"`
	IsBorrow   bool       `json:"is_borrow"`
	CreatedAt  time.Time  `json:"created_at"`
	BorrowedAt *time.Time `json:"borrowed_at,omitempty"`
}
