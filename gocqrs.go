package gocqrs

import (
	//"errors"
	//"log"
)

//func ignore() { log.Println("") }

type Command struct {
	//Id int64 `json:"_id"`
	Type int64 `json:"__type"`
	//Version int64 `json:"__version"`
	
}

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