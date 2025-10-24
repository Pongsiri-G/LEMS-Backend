package enums

type RequestType string

const (
	RequestTypeLost    RequestType = "LOST"
	RequestTypeBroken  RequestType = "BROKEN"
	RequestTypeRequest RequestType = "REQUEST"
)

type RequestStatus string

const (
	RequestStatusPending RequestStatus = "PENDING"
	RequestStatusAccept  RequestStatus = "ACCEPTED"
	RequestStatusReject  RequestStatus = "REJECTED"
	RequestStatusClosed  RequestStatus = "CLOSED"
)

func IsValidRequestType(requestType RequestType) bool {
	switch requestType {
	case RequestTypeLost, RequestTypeBroken, RequestTypeRequest:
		return true
	default:
		return false
	}
}
