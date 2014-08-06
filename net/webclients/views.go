package webclients

import (
	"errors"
	"fmt"
	"github.com/vizidrix/gocqrs/cqrs"
)

var (
	ErrInvalidSession    = errors.New("Cannot locate webclient by invalid session")
	ErrActiveSession     = errors.New("Cannot register webclient by active session")
	ErrInactiveWebClient = errors.New("No sessions to delete for an inactive webclient")
)

type WebClientSessionsView struct {
	WebClients map[string]uint64
}

func NewWebClientSessionsView() WebClientSessionsView {
	return WebClientSessionsView{
		WebClients: make(map[string]uint64),
	}
}

func (view *WebClientSessionsView) RegisterWebClientBySession(session string, webclient uint64) error {
	if _, active := view.WebClients[session]; !active {
		view.WebClients[session] = webclient
		return nil
	} else {
		view.WebClients[session] = webclient
		return ErrActiveSession
	}
}

func (view *WebClientSessionsView) GetBySession(session string) (uint64, error) {
	if webclient, valid := view.WebClients[session]; !valid {
		return 0, ErrInvalidSession
	} else {
		return webclient, nil
	}
}

func (view *WebClientSessionsView) DeleteByWebClient(webclientid uint64) error {
	sessions := 0
	for session, webclient := range view.WebClients {
		if webclient == webclientid {
			delete(view.WebClients, session)
			sessions++
		}
	}

	if sessions > 0 {
		return nil
	} else {
		return ErrInactiveWebClient
	}
}

func (view *WebClientSessionsView) HandleEvent(newevent cqrs.Event) {
	switch event := newevent.(type) {
	case WebClientAdded:
		//		fmt.Printf("\nWebClient View Handler: Adding webclient %d to View", event.GetId())
		view.RegisterWebClientBySession(event.Session, event.GetId())
	case WebClientSessionUpdated:
		//		fmt.Printf("\nWebClient View Handler: Updating session for webclient %d in View", event.GetId())
		view.DeleteByWebClient(event.GetId())
		view.RegisterWebClientBySession(event.Session, event.GetId())
	case WebClientRemoved:
		//		fmt.Printf("\nWebClient View Handler: Removing webclient %d from View", event.GetId())
		view.DeleteByWebClient(event.GetId())
	default:
		fmt.Println(errors.New("Invalid webclient view event"))
	}
}

func (view *WebClientSessionsView) HandleStream(eventchan chan cqrs.Event) {
	for {
		select {
		case event := <-eventchan:
			view.HandleEvent(event)
		}
		fmt.Printf("\nWebClients ->\n\t %+v", view)
	}
}

func NewWebClientSessionsViewHandler(eventbus cqrs.EventRouter, view *WebClientSessionsView) error {
	subscription, err := eventbus.Subscribe(cqrs.ByEventTypes(
		E_WebClientAdded,
		E_WebClientSessionUpdated,
		E_WebClientRemoved,
	))
	if err != nil {
		return err
	}

	go view.HandleStream(subscription.EventChan())

	return nil
}
