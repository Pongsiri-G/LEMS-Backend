package enums

type AuthProvider string

const (
	Local  AuthProvider = "LOCAL"
	Google AuthProvider = "GOOGLE"
)

type UserStatus string

const (
	Pending     UserStatus = "PENDING"
	Active      UserStatus = "ACTIVE"
	Rejected    UserStatus = "REJECTED"
	Deactivated UserStatus = "DEACTIVATED"
)

type UsesrRole string

const (
	user  UsesrRole = "USER"
	admin UsesrRole = "ADMIN"
)
