package enums

type RequestItemStatus string

const (
	RequestItemStatusLost   RequestItemStatus = "LOST"
	RequestItemStatusBroken RequestItemStatus = "BROKEN"
)

type RequestType string

const (
	RequestTypeLost    RequestType = "LOST"
	RequestTypeBroken  RequestType = "BROKEN"
	RequestTypeRequest RequestType = "REQUEST"
)
