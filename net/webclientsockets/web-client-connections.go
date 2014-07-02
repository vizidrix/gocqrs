package webclientsockets

import (
	"github.com/vizidrix/gocqrs/cqrs"
)

type WebClientConnection interface {
	WebClient() uint64
	EventChan() chan cqrs.Event
	MessageChan() chan []byte
	ExitChan() chan struct{}
}

type ConnectionMemento struct {
	session     string
	webclient   uint64
	eventChan   chan cqrs.Event
	messageChan chan []byte
	exitChan    chan struct{}
}

func NewConnection(session string, webclient uint64) ConnectionMemento {
	return ConnectionMemento{
		session:     session,
		webclient:   webclient,
		eventChan:   make(chan cqrs.Event),
		messageChan: make(chan []byte),
		exitChan:    make(chan struct{}),
	}
}

func (connection *ConnectionMemento) WebClient() uint64 {
	return connection.webclient
}

func (connection *ConnectionMemento) EventChan() chan cqrs.Event {
	return connection.eventChan
}

func (connection *ConnectionMemento) MessageChan() chan []byte {
	return connection.messageChan
}

func (connection *ConnectionMemento) ExitChan() chan struct{} {
	return connection.exitChan
}

type ConnectionService struct {
	connections      map[uint64]*ConnectionMemento
	addChan          chan *ConnectionMemento
	removeChan       chan *ConnectionMemento
	subscriptionChan chan WebClientConnection
}

func NewConnectionService(subscriptionchan chan WebClientConnection) ConnectionService {
	return ConnectionService{
		connections:      make(map[uint64]*ConnectionMemento),
		addChan:          make(chan *ConnectionMemento),
		removeChan:       make(chan *ConnectionMemento),
		subscriptionChan: subscriptionchan,
	}
}
