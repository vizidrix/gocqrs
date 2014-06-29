package clients

import (
	"fmt"
	"github.com/vizidrix/gocqrs/cqrs"
)

type ClientState struct {
	cqrs.DomainState
}

func NewClientState(availableId uint64, activeIds map[uint64]bool) ClientState {
	return ClientState{
		DomainState: cqrs.NewDomainState(availableId, activeIds),
	}
}

func LoadClientState() ClientState {
	// TODO: Read initial state from data store
	used := make(map[uint64]bool)
	used[0] = true
	return NewClientState(1, used)
}

func LoadClientInitialState() ClientState {
	// TODO: Read initial state from data store
	used := make(map[uint64]bool)
	used[0] = true
	return NewClientState(1, used)
}

func NewClientCommandHandler(eventbus cqrs.EventRouter, state *ClientState) (func(command cqrs.Command), error) {

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
				eventbus.Publish(NewClientRemoved(cmd.GetId(), cmd.GetVersion(), cmd.Session))
			}
		default:
			{
				fmt.Printf("\nCommand Handler: Unable to handle command: [ %v ]", cmd)
			}
		}
	}, nil
}
