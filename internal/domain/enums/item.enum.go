package enums

type ItemStatus string

const (
	ItemStatusAvailable  ItemStatus = "AVAILABLE"
	ItemStatusUnBorrowAble ItemStatus = "UNBORROWABLE"
	ItemStatusOutOfStock ItemStatus = "UNAVAILABLE"
)
