package clientsockets

import (
	"github.com/vizidrix/gocqrs/cqrs"
)

type Connection struct {
	session     string
	client      uint64
	eventChan   chan cqrs.Event
	messageChan chan []byte
	exitChan    chan struct{}
}

func NewConnection(session string, client uint64) Connection {
	return Connection{
		session:     session,
		client:      client,
		eventChan:   make(chan cqrs.Event),
		messageChan: make(chan []byte),
		exitChan:    make(chan struct{}),
	}
}

func (connection *Connection) Client() uint64 {
	return connection.client
}

func (connection *Connection) EventChan() chan cqrs.Event {
	return connection.eventChan
}

func (connection *Connection) MessageChan() chan []byte {
	return connection.messageChan
}

func (connection *Connection) ExitChan() chan struct{} {
	return connection.exitChan
}
