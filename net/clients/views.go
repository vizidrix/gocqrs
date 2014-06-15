package clients

import (
	"errors"
)

type ClientSessionView struct {
	Clients map[string]uint64
}

func NewClientSessionView() ClientSessionView {
	return ClientSessionView{
		Clients: make(map[string]uint64),
	}
}

/*
func (view *ClientView) NewBySession(id string) (uint64, error) {
	for clientid := 1; clientid < len(view.Clients)+2; clientid++ {
		if _, inuse := view.Clients[id]; !inuse {
			view.Clients[id] =
		}
	}
}
*/

func (view *ClientSessionView) GetBySession(session string) (uint64, error) {
	if client, valid := view.Clients[session]; !valid {
		return 0, errors.New("invalid_session_id")
	} else {
		return client, nil
	}
}
