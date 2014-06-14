package clientsockets

import (
	"github.com/vizidrix/gocqrs/cqrs"
)

type ClientConn struct {
	Session     string
	Client      uint64
	EventChan   chan cqrs.Event
	MessageChan chan []byte
	ExitChan    chan struct{}
}

func NewClientConn(session string, client uint64) ClientConn {
	return ClientConn{
		Session:     session,
		Client:      client,
		EventChan:   make(chan cqrs.Event),
		MessageChan: make(chan []byte),
		ExitChan:    make(chan struct{}),
	}
}
