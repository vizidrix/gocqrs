package cqrs_test

import (
	"fmt"
	. "github.com/vizidrix/gocqrs/cqrs"
	"testing"
)

var (
	TEST_DOMAIN uint32 = 0x11111111
)

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

func Test_Should_return_correct_domain(t *testing.T) {
	var expected = "0x11111111"
	var commandtype = C(TEST_DOMAIN, 1, 1)
	var command = NewCommand(commandtype, 0, 0)

	if fmt.Sprintf("%#x", command.GetDomain()) != expected {
		t.Errorf("Expected [ %s ] but received [ %#x ]\n", expected, command.GetDomain())
	}
}
