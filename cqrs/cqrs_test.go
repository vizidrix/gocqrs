package cqrs_test

import (
	"fmt"
	. "github.com/vizidrix/gocqrs/cqrs"
	"testing"
)

var (
	TEST_DOMAIN   uint32 = 0x11111111
	B_TestCommand uint64 = C(TEST_DOMAIN, 1, 1)
	B_TestEvent   uint64 = E(TEST_DOMAIN, 1, 1)
)

type TestCommand struct {
	CommandMemento
	Value string
}

func NewTestCommand(id uint64, version uint32, value string) TestCommand {
	return TestCommand{
		CommandMemento: NewCommand(TEST_DOMAIN, id, version, C_TestCommand),
		Value:          value,
	}
}

type TestEvent struct {
	EventMemento
	Value string
}

func NewTestEvent(id uint64, version uint32, value string) TestEvent {
	return TestEvent{
		EventMemento: NewEvent(TEST_DOMAIN, id, version, E_TestEvent),
		Value:        value,
	}
}

func Test_Should_calculate_correct_command_key(t *testing.T) {
	var expected = "0x1111111180010001"
	var key = C(TEST_DOMAIN, 1, 1)

	if fmt.Sprintf("%#x", key) != expected {
		t.Errorf("Expected [ %s ] but received [ %#x ]\n", expected, key)
	}
}

func Test_Should_calculate_correct_event_key(t *testing.T) {
	var expected = "0x1111111100010001"
	var key = E(TEST_DOMAIN, 1, 1)

	if fmt.Sprintf("%#x", key) != expected {
		t.Errorf("Expected [ %s ] but received [ %#x ]\n", expected, key)
	}
}

func Test_Should_return_correct_command_domain(t *testing.T) {
	var expected = "0x11111111"
	var command = NewCommand(TEST_DOMAIN, 0, 0, B_TestCommand)
	var result = command.GetDomain()

	if fmt.Sprintf("%#x", result) != expected {
		t.Errorf("Expected [ %s ] but received [ %#x ]\n", expected, result)
	}
}

func Test_Should_return_correct_event_domain(t *testing.T) {
	var expected = "0x11111111"
	var event = NewEvent(TEST_DOMAIN, 0, 0, B_TestEvent)
	var result = event.GetDomain()

	if fmt.Sprintf("%#x", result) != expected {
		t.Errorf("Expected [ %s ] but received [ %#x ]\n", expected, result)
	}
}
