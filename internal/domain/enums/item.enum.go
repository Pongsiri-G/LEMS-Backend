package enums

type ItemStatus string

const (
	ItemStatusAvailable  ItemStatus = "AVAILABLE"
	ItemStatusOutOfStock ItemStatus = "UNAVAILABLE"
)
