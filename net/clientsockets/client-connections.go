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
		messageChan: make(chan []byte, 1),
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

type ConnectionService struct {
	connections      map[uint64]*Connection
	addChan          chan *Connection
	removeChan       chan *Connection
	subscriptionChan chan *Connection
}

func NewConnectionService(subscriptionchan chan *Connection) ConnectionService {
	return ConnectionService{
		connections:      make(map[uint64]*Connection),
		addChan:          make(chan *Connection),
		removeChan:       make(chan *Connection),
		subscriptionChan: subscriptionchan,
	}
}
