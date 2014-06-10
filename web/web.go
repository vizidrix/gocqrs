package web

type CorrelationMemento struct {
	Client      uint64 `json:"__clientid"`
	Correlation uint64 `json:"__correlation"`
}