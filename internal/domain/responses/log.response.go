package responses

import "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"

type AdminLogResponse struct {
	LogID    string        `json:"log_id"`
	UserID   string        `json:"user_id"`
	UserName string        `json:"user_name"`
	LogType  enums.LogType `json:"log_type"`
	Message  string        `json:"message,omitempty"`
	Created  string        `json:"created_at"`
}

/*
type Log struct {
	LogID      uuid.UUID     `db:"log_id" gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID     `db:"user_id" gorm:"type:uuid;foreignKey:UserID;references:UserID"`
	LogType    enums.LogType `db:"log_type"`
	LogMessage *string       `db:"log_message"`
	CreatedAt  time.Time     `db:"created_at"`
}
*/
