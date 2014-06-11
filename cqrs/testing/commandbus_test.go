package cqrs_test

import (
	. "github.com/vizidrix/gocqrs/cqrs"
	"testing"
)

var (
	DOMAIN        uint32 = 0x11111111
	C_TestCommand uint32 = C(1, 1)
)

type MockCommandBus struct {
	Registrations map[uint32]CommandHandler
	RegisterChan  chan CommandHandler
	PublishChan   chan Command
	CommandChan   chan Command
	CancelChan    chan struct{}
}

func NewMockCommandBus() *MockCommandBus {
	return &MockCommandBus{
		//		Registrations: make(map[uint32]CommandHandler),
		RegisterChan: make(chan CommandHandler, 1),
		PublishChan:  make(chan Command, 1),
		CommandChan:  make(chan Command, 1),
		CancelChan:   make(chan struct{}),
	}
}

func (mock *MockCommandBus) Create() CommandRouter {
	return NewChannelCommandBus(
		mock.RegisterChan,
		mock.PublishChan,
		func() chan Command { return mock.CommandChan },
	)
}

type MockHandler struct{}

func (mock *MockHandler) CommandChan() <-chan Command {
	return nil
}

func (mock *MockHandler) Publish(command Command) {
	return
}

func (mock *MockHandler) Domain() uint32 {
	return DOMAIN
}

type TestCommand struct {
	CommandMemento
	Value string
}

func NewTestCommand(id uint64, version uint32, value string) TestCommand {
	return TestCommand{
		CommandMemento: NewCommand(DOMAIN, id, version, C_TestCommand),
		Value:          value,
	}
}

func Test_Should_return_nil_for_empty_domain_registration(t *testing.T) {
	commandbus := NewMockCommandBus().Create()
	_, err := commandbus.Register(0)

	if err != ErrInvalidDomainRegistered {
		t.Errorf("Should have returned an error for invalid domain but was [ %v ]\n", err)
	}

}
