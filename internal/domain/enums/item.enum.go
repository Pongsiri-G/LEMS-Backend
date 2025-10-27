package enums

type ItemStatus string

const (
	ItemStatusAvailable  ItemStatus = "AVAILABLE"
	ItemStatusInLabOnly ItemStatus = "INLABONLY"
	ItemStatusOutOfStock ItemStatus = "UNAVAILABLE"
)
