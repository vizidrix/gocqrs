package clients

import (
	"errors"
	"fmt"
	"github.com/vizidrix/gocqrs/cqrs"
)

var (
	ErrDuplicateAddedId = errors.New("Cannot add an added id to state")
	ErrInvalidRemoveId  = errors.New("Cannot remove invalid id from state")
	ErrRemoveRemovedId  = errors.New("Cannot remove a removed id from state")
)

type ClientCommandHandlerState struct {
	nextId    uint64
	activeIds map[uint64]bool
}

func (state *ClientCommandHandlerState) AvailableId() uint64 {
	id := state.nextId
	state.nextId++
	return id
}

/*
func (state *ClientCommandHandlerState) IncrementAvailable() {
	state.nextId++
}
*/

func (state *ClientCommandHandlerState) AddId(id uint64) error {
	if _, added := state.activeIds[id]; added {
		return ErrDuplicateAddedId
	} else {
		state.activeIds[id] = true
	}

	return nil
}

func (state *ClientCommandHandlerState) RemoveId(id uint64) error {
	if active, added := state.activeIds[id]; !added {
		return ErrInvalidRemoveId
	} else {
		if !active {
			return ErrRemoveRemovedId
		}
		state.activeIds[id] = false
	}
	return nil
}

func NewClientCommandHandlerState(availableId uint64, activeIds map[uint64]bool) *ClientCommandHandlerState {
	return &ClientCommandHandlerState{
		nextId:    availableId,
		activeIds: activeIds,
	}
}

func LoadClientCommandHandlerInitialState() *ClientCommandHandlerState {
	// TODO: Read initial state from data store
	used := make(map[uint64]bool)
	used[0] = true
	return NewClientCommandHandlerState(1, used)
}

func NewClientCommandHandler(eventbus cqrs.EventRouter, state *ClientCommandHandlerState) (func(command cqrs.Command), error) {

	return func(command cqrs.Command) {
		switch cmd := command.(type) {
		case AddClient:
			{
				newid := state.AvailableId()               // Retrieve available id from state
				if err := state.AddId(newid); err != nil { // Update state to include this id as active
					fmt.Errorf("%#v", err)
					break
				}
				// state.IncrementAvailable() // Update state to

				fmt.Printf("\nCommand Handler: Trying to add client->\n\t%v", command)
				eventbus.Publish(NewClientAdded(newid, cmd.Session))
			}
		case UpdateClientSession:
			{
				fmt.Printf("\nCommand Handler: Trying to update entire client->\n\t%v", command)
				eventbus.Publish(NewClientSessionUpdated(cmd.GetId(), cmd.GetVersion(), cmd.Session))
			}
		case RemoveClient:
			{
				if err := state.RemoveId(cmd.GetId()); err != nil {
					fmt.Errorf("%#v", err)
					break
				}

				fmt.Printf("\nCommand Handler: Trying to delete client->\n\t%v", command)
				eventbus.Publish(NewClientRemoved(cmd.GetId(), cmd.GetVersion()))
			}
		default:
			{
				fmt.Printf("\nCommand Handler: Unable to handle command: [ %v ]", cmd)
			}
		}
	}, nil
}
