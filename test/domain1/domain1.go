package gocqrs_test_domain1

import (
	"github.com/vizidrix/gocqrs/cqrs"
	"fmt"
)

var DOMAIN int32 = 0

const ( // Aggregates
	_ = 0 + iota
	Person
)
const (
	_ = 0 + iota
	DoSomething // Overlapping with main test domain commands
	DoAnotherThing
	DoSomeThirdThing // Unique to this test domain
)

type DoSomethingCommand struct {
	cqrs.Command
	Data string `json:"data"`
}

type DoAnotherThingCommand struct {
	cqrs.Command
	Data int32 `json:"data"`
}

type DoSomeThirdThingCommand struct {
	cqrs.Command
	Value1 int32 `json:"value1"`
	Value2 int32 `json:"value2"`
}

func NewDoSomethingCommand(data string) *DoSomethingCommand {
	return &DoSomethingCommand {
		cqrs.NewCommand(DOMAIN, DoSomething, 0),
		data,
	}
}

func NewDoAnotherThingCommand(data int32) *DoAnotherThingCommand {
	return &DoAnotherThingCommand {
		cqrs.NewCommand(DOMAIN, DoAnotherThing, 0),
		data,
	}
}

func NewDoSomeThirdThingCommand(value1 int32, value2 int32) *DoSomeThirdThingCommand {
	return &DoSomeThirdThingCommand {
		cqrs.NewCommand(DOMAIN, DoSomeThirdThing, 0),
		value1,
		value2,
	}
}

func HandleDoSomethingCommand(cmd *DoSomethingCommand) {
	if cmd.Data == "" {
		fmt.Printf("Empty Something\n")
	} else {
		fmt.Printf("Got Something: [ %s ]\n", cmd.Data)
	}
}

func HandleDoAnotherThingCommand(cmd *DoAnotherThingCommand) {
	if cmd.Data == 0{
		fmt.Printf("Another empty\n")
	} else {
		fmt.Printf("Got Another [ %s ]\n", cmd.Data)
	}
}