package responses

import (
	"github.com/google/uuid"
)

type TagResponse struct {
	TagID    uuid.UUID `json:"id"`
	TagName  string    `json:"name"`
	TagColor string    `json:"color"`
}
