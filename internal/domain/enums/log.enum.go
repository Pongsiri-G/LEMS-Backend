package enums

type LogType string

const (
	LogTypeBorrow   LogType = "BORROW"
	LogTypeReturn   LogType = "RETURN"
	LogTypeLogin    LogType = "LOGIN"
	LogTypeRegister LogType = "REGISTER"
)
