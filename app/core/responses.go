package core

// SimpleTransaction is used to retrreive transaction from request
type SimpleTransaction struct {
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}
