package webclients

import (
	"fmt"
	"github.com/vizidrix/gocqrs/cqrs"
)

type WebClientState struct {
	cqrs.DomainState
}

func NewWebClientState(availableId uint64, activeIds map[uint64]bool) WebClientState {
	return WebClientState{
		DomainState: cqrs.NewDomainState(availableId, activeIds),
	}
}

func LoadWebClientState() WebClientState {
	// TODO: Read initial state from data store
	used := make(map[uint64]bool)
	used[0] = true
	return NewWebClientState(1, used)
}

func LoadWebClientInitialState() WebClientState {
	// TODO: Read initial state from data store
	used := make(map[uint64]bool)
	used[0] = true
	return NewWebClientState(1, used)
}

func NewWebClientCommandHandler(eventbus cqrs.EventRouter, state *WebClientState) (func(command cqrs.Command), error) {

	return func(command cqrs.Command) {
		switch cmd := command.(type) {
		case AddWebClient:
			{
				newid := state.AvailableId()               // Retrieve available id from state
				if err := state.AddId(newid); err != nil { // Update state to include this id as active
					fmt.Errorf("%#v", err)
					break
				}
				// state.IncrementAvailable() // Update state to

				fmt.Printf("\nCommand Handler: Trying to add webclient->\n\t%v", command)
				eventbus.Publish(NewWebClientAdded(newid, cmd.Session))
			}
		case UpdateWebClientSession:
			{
				fmt.Printf("\nCommand Handler: Trying to update entire webclient->\n\t%v", command)
				eventbus.Publish(NewWebClientSessionUpdated(cmd.GetId(), cmd.GetVersion(), cmd.Session))
			}
		case RemoveWebClient:
			{
				if err := state.RemoveId(cmd.GetId()); err != nil {
					fmt.Errorf("%#v", err)
					break
				}

				fmt.Printf("\nCommand Handler: Trying to delete webclient->\n\t%v", command)
				eventbus.Publish(NewWebClientRemoved(cmd.GetId(), cmd.GetVersion(), cmd.Session))
			}
		default:
			{
				fmt.Printf("\nCommand Handler: Unable to handle command: [ %v ]", cmd)
			}
		}
	}, nil
}
