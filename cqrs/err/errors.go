package err

import (
	"github.com/vizidrix/gocqrs/cqrs"
)

var DOMAIN_NAME = "err"
var DOMAIN = 0x125FDBBD

type Error interface {
	cqrs.Event
	GetError() string
}

type ErrorMemento struct {
	cqrs.EventMemento
	Error error `json:"__error"`
}

func NewError(err error) ErrorMemento {
	return ErrorMemento{
		EventMemento: cqrs.NewEvent(0, 0, 0),
		Error:        err,
	}
}
