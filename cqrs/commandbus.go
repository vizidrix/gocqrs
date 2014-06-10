package cqrs

import (
	"errors"
)

var (
	ErrInvalidCommandBusState = errors.New("Command bus not properly initialized")
	ErrInvalidNilRegister = errors.New("Cannot register nil command type set")
	ErrDuplicateRegistration = errors.New("Cannot register duplicate command type handlers")
)

var CommandBus CommandRouter

func init() {
	// Auto init the command bus as a global service
	CommandBus = NewDefaultedCommandBus()
	CommandBus.Listen()
}

type CommandChanFactory func() chan Command

type CommandHandler interface {
	CommandChan() <-chan Command
	Publish(command Command)
	Cancel()
}

type CommandRouter interface {
	Listen() // Iterates across the Step function in a goroutine loop
	Step() // Grabs the next operation from teh queue and processes it
	Publish(command Command) (error) // Pushes the command to the registerd command handler or err
	Register(commandTypes []uint32) (CommandHandler, error) // Registers a command handler for the set of commands provided
}

// CommandRouter implementation that uses Go chans to provide routing
type channelCommandBus struct {
	registrations map[uint32]chan Command
	registerChan chan []uint32
	publishChan chan Command
	commandChanFactory CommandChanFactory
}

func NewDefaultedCommandBus() *channelCommandBus {
	return NewChannelCommandBus(
		make(chan []uint32),
		make(chan Command),
		func() chan Command { return make(chan Command)},
		)
}

func NewChannelCommandBus(
	registerChan chan []uint32,
	publishChan chan Command,
	commandChanFactory CommandChanFactory,
	) *channelCommandBus {
	bus := &channelCommandBus {
		registerChan: registerChan,
		publishChan: publishChan,
		commandChanFactory: commandChanFactory,
	}
	return bus
}

func (c *channelCommandBus) Step() {
	select { // Synchronized select for command bus mutable actions
	case newRegistrations := <-c.registerChan: {
		for _, commandType := range newRegistrations {
			if _, ok := c.registrations[commandType]; ok {
				panic(ErrDuplicateRegistration)
			}
		}
	}
	}
}
func (c *channelCommandBus) Listen() {

}

func (c *channelCommandBus) Publish(command Command) (error) {
	return nil
}

func (c *channelCommandBus) Register(commandTypes ...uint32) (CommandHandler, error) {
	return nil, nil
}
