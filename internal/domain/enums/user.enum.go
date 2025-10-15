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

type UserRole string

const (
	User  UserRole = "USER"
	Admin UserRole = "ADMIN"
)
