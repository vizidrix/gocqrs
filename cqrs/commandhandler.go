package cqrs

import (
	"errors"
)

/* TODO: Command Handler Interface Type */

var (
	ErrDuplicateAddedId = errors.New("Cannot add an added id to state")
	ErrInvalidRemoveId  = errors.New("Cannot remove invalid id from state")
	ErrRemoveRemovedId  = errors.New("Cannot remove a removed id from state")
)

type CommandHandlerState interface {
	AvailableId() uint64
	AddId(id uint64) error
	RemoveId(id uint64) error
}

type DomainState struct {
	nextId    uint64
	activeIds map[uint64]bool
}

func (state *DomainState) AvailableId() uint64 {
	id := state.nextId
	state.nextId++
	return id
}

/*
func (state *WorkflowCommandHandlerState) IncrementAvailable() {
	state.nextId++
}
*/

func (state *DomainState) AddId(id uint64) error {
	if _, added := state.activeIds[id]; added {
		return ErrDuplicateAddedId
	} else {
		state.activeIds[id] = true
	}

	return nil
}

func (state *DomainState) RemoveId(id uint64) error {
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

func NewDomainState(availableId uint64, activeIds map[uint64]bool) DomainState {
	return DomainState{
		nextId:    availableId,
		activeIds: activeIds,
	}
}
