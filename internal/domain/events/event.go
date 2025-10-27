package events

type EventType string

const (
	ItemAvaliable  EventType = "ITEM_AVALIABLE"
	UserDeleted    EventType = "USER_DELETED"
)

type Event struct {
	Type    EventType
	Payload interface{}
}
