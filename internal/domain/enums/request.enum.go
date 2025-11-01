package enums

type RequestType string

const (
	RequestTypeLost    RequestType = "LOST"
	RequestTypeBroken  RequestType = "BROKEN"
	RequestTypeRequest RequestType = "REQUEST"
)

func StringToRequestType(s string) *RequestType {
	var reqType RequestType
	switch s {
	case "LOST":
		reqType = RequestTypeLost
		return &reqType
	case "BROKEN":
		reqType = RequestTypeBroken
		return &reqType
	case "REQUEST":
		reqType = RequestTypeRequest
		return &reqType
	default:
		return nil
	}
}

type RequestStatus string

const (
	RequestStatusPending  RequestStatus = "PENDING"
	RequestStatusCancel   RequestStatus = "CANCELED"
	RequestStatusAccept   RequestStatus = "ACCEPTED"
	RequestStatusReject   RequestStatus = "REJECTED"
	RequestStatusComplete RequestStatus = "COMPLETED"
)

func StringToRequestStatus(s string) *RequestStatus {
	var reqStatus RequestStatus
	switch s {
	case "PENDING":
		reqStatus = RequestStatusPending
		return &reqStatus
	case "CANCELED":
		reqStatus = RequestStatusCancel
		return &reqStatus
	case "ACCEPTED":
		reqStatus = RequestStatusAccept
		return &reqStatus
	case "REJECTED":
		reqStatus = RequestStatusReject
		return &reqStatus
	case "COMPLETED":
		reqStatus = RequestStatusComplete
		return &reqStatus
	default:
		return nil
	}
}

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
	ExportTypeXLS  ExportType = "XLS"
	ExportTypePDF  ExportType = "PDF"
	ExportTypeCSV  ExportType = "CSV"
	ExportTypeJSON ExportType = "JSON"
)
