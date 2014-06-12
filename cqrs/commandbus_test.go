package cqrs_test

import (
	. "github.com/vizidrix/gocqrs/cqrs"
	"testing"
)

var (
	C_TEST_DOMAIN uint32 = 0x11111111
	C_TestCommand uint64 = C(C_TEST_DOMAIN, 1, 1)
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

func (mock *MockHandler) CommandChan() chan Command {
	return nil
}

func (mock *MockHandler) Domain() uint32 {
	return C_TEST_DOMAIN
}

func Test_Should_return_error_for_empty_domain(t *testing.T) {
	commandbus := NewMockCommandBus().Create()
	handle, err := commandbus.Register(0)

	if err != ErrInvalidDomainRegistered {
		t.Errorf("Should have returned an error for invalid domain but was [ %v ]\n", err)
	}
	if handle != nil {
		t.Errorf("Should have returned nil registration handle but returned [ %v ]\n", handle)
	}
}

func Test_Should_return_registration_token_for_valid_domain(t *testing.T) {
	commandbus := NewMockCommandBus().Create()
	handle, err := commandbus.Register(C_TEST_DOMAIN)

	if err != nil {
		t.Errorf("Should not have err but was [ %s ]\n", err)
		return
	}
	if handle == nil {
		t.Errorf("Should have returned a non nil handle\n")
		return
	}
}

func Test_Should_return_error_when_publishing_nil_command(t *testing.T) {
	commandbus := NewMockCommandBus().Create()

	if err := commandbus.Publish(nil); err != ErrInvalidNilPublishedCommand {
		t.Errorf("Should have returned an error for nil command in publish but was [ %v ]\n", err)
		return
	}
}

func Test_Should_not_return_error_from_valid_command_publish(t *testing.T) {
	commandbus := NewMockCommandBus().Create()
	command := NewTestCommand(1, 1, "publish test")

	if err := commandbus.Publish(command); err != nil {
		t.Errorf("Should not return error from valid publish")
	}
}

func Test_Should_receive_matching_command_in_handler(t *testing.T) {
	commandbus := NewMockCommandBus().Create()
	handle, _ := commandbus.Register(C_TEST_DOMAIN)
	commandbus.Step()
	expected := NewTestCommand(1, 1, "publish test")
	commandbus.Publish(expected)
	commandbus.Step()

	select {
	case actual := <-handle.CommandChan():
		{
			if expected != actual {
				t.Errorf("\nExpect\t[ %v ]\n\tbut was\t[ %v ]\n", expected, actual)
			}
		}
	default:
		{
			t.Errorf("Should not have hit default case\n")
		}
	}
}
