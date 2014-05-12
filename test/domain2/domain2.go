package gocqrs_test_domain2

import (
	cqrs "github.com/vizidrix/gocqrs"
)

var DOMAIN int32 = 1

const (
	_ = 0 + iota
	DoSomething
	DoAnotherThing
)

type DoSomethingCommand struct {
	cqrs.Command
	StringValue string `json:"stringalue"`
	IntValue int64 `json:"intvalue"`
}

type DoAnotherThingCommand struct {
	cqrs.Command
	StringValue string `json:"stringvalue"`
	IntValue int64 `json:"intvalue"`
}

func NewDoSomethingCommand(stringValue string, intValue int64) *DoSomethingCommand {
	return &DoSomethingCommand {
		//cqrs.Command { Domain: 1, Type: 1 },
		cqrs.NewCommand(CQRS_DOMAIN, DoSomething, 0),
		stringValue,
		intValue,
	}
}

func NewDoAnotherThingCommand(stringValue string, intValue int64) *DoAnotherThingCommand {
	return &DoAnotherThingCommand {
		cqrs.NewCommand(CQRS_DOMAIN, DoAnotherThing, 0),
		//cqrs.Command { Domain: 1, Type: 2 },
		stringValue,
		intValue,
	}
}

func HandleDoSomethingCommand(cmd *DoSomethingCommand) {
	if cmd.StringValue == "" {
		fmt.Printf("Empty Something\n")
	} else {
		fmt.Printf("Got Something: [ %s ]\n", cmd.StringValue)
	}
}

func HandleDoAnotherThingCommand(cmd *DoAnotherThingCommand) {
	if cmd.StringValue == "" {
		fmt.Printf("Another empty\n")
	} else {
		fmt.Printf("Got Another [ %s ]\n", cmd.StringValue)
	}
}