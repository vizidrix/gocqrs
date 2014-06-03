package cqrs_test

import (
	"fmt"
	"testing"
	"github.com/vizidrix/gocqrs/cqrs"
)

func Test_Should_calculate_correct_command_key(t *testing.T) {
	var expected = "0x80010001"
	var key = cqrs.C(1,1)

	if fmt.Sprintf("%#x", key) != expected {
		t.Errorf("Expected [ %s ] but received [ %#x ]\n", expected, key)
	}
}

func Test_Should_calculate_correct_event_key(t *testing.T) {
	var expected = "0x10001"
	var key = cqrs.E(1,1)

	if fmt.Sprintf("%#x", key) != expected {
		t.Errorf("Expected [ %s ] but received [ %#x ]\n", expected, key)
	}
}