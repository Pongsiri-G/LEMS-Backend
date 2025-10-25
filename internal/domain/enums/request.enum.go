package enums

type RequestType string

const (
	RequestTypeLost    RequestType = "LOST"
	RequestTypeBroken  RequestType = "BROKEN"
	RequestTypeRequest RequestType = "REQUEST"
)

type RequestStatus string

const (
	RequestStatusPending  RequestStatus = "PENDING"
	RequestStatusCancel   RequestStatus = "CANCELLED"
	RequestStatusAccept   RequestStatus = "ACCEPTED"
	RequestStatusReject   RequestStatus = "REJECTED"
	RequestStatusComplete RequestStatus = "COMPLETED"
)

func IsValidRequestType(requestType RequestType) bool {
	switch requestType {
	case RequestTypeLost, RequestTypeBroken, RequestTypeRequest:
		return true
	default:
		return false
	}
}

type ExportType string

const (
	ExportTypeXLS ExportType = "XLS"
	ExportTypePDF ExportType = "PDF"
)
