package enums

type LogType string

const (
	LogTypeBorrow   LogType = "BORROW"
	LogTypeReturn   LogType = "RETURN"
	LogTypeLogin    LogType = "LOGIN"
	LogTypeRegister LogType = "REGISTER"

	// Admin Actions
	LogTypeAccept      LogType = "ACCEPT"
	LogTypeReject      LogType = "REJECT"
	LogTypeActivate    LogType = "ACTIVATE"
	LogTypeDeactivate  LogType = "DEACTIVATE"
	LogTypeDelete      LogType = "DELETE"
	LogTypeGrantAdmin  LogType = "GRANT_ADMIN"
	LogTypeRevokeAdmin LogType = "REVOKE_ADMIN"
)
