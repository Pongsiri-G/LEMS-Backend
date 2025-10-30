package events

type EventType string

const (
	ItemAvaliable  EventType = "Item is Avaliable"
	UserDeleted    EventType = "USER_DELETED"
)

type Event struct {
	Type    EventType
	Payload interface{}
}
