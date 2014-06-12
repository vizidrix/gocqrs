package cqrs

import (
	"errors"
)

var (
	ErrInvalidCommandBusState     = errors.New("Command bus not properly initialized")
	ErrInvalidDomainRegistered    = errors.New("Cannot register nil command type set")
	ErrDuplicateRegistration      = errors.New("Cannot register duplicate command type handlers")
	ErrInvalidNilPublishedCommand = errors.New("Cannot publish a nil command")
)

var CommandBus CommandRouter

/*
func init() {
	// Auto init the command bus as a global service
	CommandBus = NewDefaultedCommandBus()
	CommandBus.Listen()
}
*/

type CommandChanFactory func() chan Command

type CommandHandler interface {
	CommandChan() chan Command
	Domain() uint32
}

type CommandRouter interface {
	Listen()                                        // Iterates across the Step function in a goroutine loop
	Step()                                          // Grabs the next operation from teh queue and processes it
	Publish(command Command) error                  // Pushes the command to the registerd command handler or err
	Register(domain uint32) (CommandHandler, error) // Registers a command handler for the set of commands provided
}

type registration struct {
	commandBus  CommandRouter
	commandChan chan Command
	domain      uint32
}

func (r *registration) CommandChan() chan Command {
	return r.commandChan
}

func (r *registration) Domain() uint32 {
	return r.domain
}

func SendCommand(c CommandHandler, command Command) {
	c.CommandChan() <- command
}

// CommandRouter implementation that uses Go chans to provide routing
type channelCommandBus struct {
	registrations      map[uint32]CommandHandler
	registerChan       chan CommandHandler
	publishChan        chan Command
	commandChanFactory CommandChanFactory
}

func NewDefaultedCommandBus() *channelCommandBus {
	return NewChannelCommandBus(
		make(chan CommandHandler),
		make(chan Command),
		func() chan Command { return make(chan Command) },
	)
}

func NewChannelCommandBus(
	registerChan chan CommandHandler,
	publishChan chan Command,
	commandChanFactory CommandChanFactory,
) *channelCommandBus {
	bus := &channelCommandBus{
		registrations:      make(map[uint32]CommandHandler),
		registerChan:       registerChan,
		publishChan:        publishChan,
		commandChanFactory: commandChanFactory,
	}
	return bus
}

func (c *channelCommandBus) Step() {
	select { // Synchronized select for command bus mutable actions
	case newRegistration := <-c.registerChan:
		{
			if _, handled := c.registrations[newRegistration.Domain()]; handled {
				panic(ErrDuplicateRegistration)
			}

			c.registrations[newRegistration.Domain()] = newRegistration
		}
	case command := <-c.publishChan:
		{
			handler := c.registrations[command.GetDomain()]
			SendCommand(handler, command)
		}
	}
}

func (c *channelCommandBus) Listen() {
	go func() {
		for {
			c.Step()
		}
	}()
}

func (c *channelCommandBus) Publish(command Command) error {
	if command == nil {
		return ErrInvalidNilPublishedCommand
	}
	select {
	case c.publishChan <- command:
	default:
		return ErrInvalidCommandBusState
	}
	return nil
}

func (c *channelCommandBus) Register(domain uint32) (CommandHandler, error) {
	if domain == 0 {
		return nil, ErrInvalidDomainRegistered
	}
	handle := &registration{
		commandBus:  c,
		commandChan: c.commandChanFactory(),
		domain:      domain,
	}
	c.registerChan <- handle
	return handle, nil
}
