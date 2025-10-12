package enums

type LogType string

const (
	LogRegister        LogType = "REGISTER"
	LogLogin           LogType = "LOGIN"
	LogBorrow          LogType = "BORROW"
	LogReturn          LogType = "RETURN"
	LogRequest         LogType = "REQUEST"
	LogAdminAccept     LogType = "ADMIN_ACCEPT"
	LogAdminReject     LogType = "ADMIN_REJECT"
	LogAdminDeactivate LogType = "ADMIN_DEACTIVATE"
	LogAdminDelete     LogType = "ADMIN_DELETE"
	LogAdminGrant      LogType = "ADMIN_GRANT"
	LogAdminRevoke     LogType = "ADMIN_REVOKE"
)
