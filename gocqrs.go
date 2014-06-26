package gocqrs

import (
	//"errors"
	//"log"
)

/*

REFER TO cqrs/cqrs.go  for the main library

Working on refactoring it from the sub folder back into root

*/














/*
type Command struct {
	CommandDomain int32 `json:"__domain"`
	CommandType int32 `json:"__type"`
	CommandId int64 `json:"__id"`
	CommandVersion int64 `json:"__version"`
}

func NewCommand(commandDomain int32, commandType int32, commandId int64) Command {
	return Command {
		CommandDomain: commandDomain,
		CommandType: commandType,
		CommandId: commandId,
		CommandVersion: 0,
	}
}

func NewVersionedCommand(commandDomain int32, commandType int32, commandId int64, commandVersion int64) Command {
	return Command {
		CommandDomain: commandDomain,
		CommandType: commandType,
		CommandId: commandId,
		CommandVersion: commandVersion,
	}
}
*/
/*
func ExtractCommand(data string) *Command {
	command := &Command{}
	return json.DeMarshal(data, &command)
	
	return &Command {
		Id: 10,
		Version: 1,
		Type: 64,
	}
}

// Below this line should go into a different library
package gocqrsauth

type Commands int64

const (
	_ Commands = 0x00000000
	AuthenticateUser = 0x0000001
	TerminateSession = 0x0000002
)

func Handle(data string) {
	cmd := cqrs.ExtractCommand(data)
	switch cmd.Type {
		case Commands.AuthenticateUser: {

		}
	}
}

type RegisterPersonCommand struct {
	cqrs.Command
	FirstName string
	LastName string
}

// Somewhere in yet another package

package coolbeans

import(
	auth "github.com/vizidrix/gocqrsauth"
)

func Main() {
	fmt.Printf("Auth command id: %s", auth.Commands.AuthenticateUser)
	correlationToken, err := cqrs.CommandBus.Send(&RegisterPersonCommand {
		Command.Id: 10,
		Command.Version: 0,
		Command.Type: auth.Commands.RegisterPersonCommand,
		FirstName: "Perry",
		LastName: "Birch",
		})

}

*/
