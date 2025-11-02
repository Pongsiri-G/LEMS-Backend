package enums

type ItemStatus string

const (
	ItemStatusAvailable  ItemStatus = "AVAILABLE"
	ItemStatusInLabOnly  ItemStatus = "IN-LAB ONLY"
	ItemStatusOutOfStock ItemStatus = "UNAVAILABLE"
)
