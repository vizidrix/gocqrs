package webclients

import (
	"testing"
)

func Test_Should_return_error_for_active_session(t *testing.T) {
	webclients := NewWebClientSessionsView()
	session := "test"
	var webclient uint64 = 1
	webclients.WebClients[session] = 0

	err := webclients.RegisterWebClientBySession(session, webclient)

	registered := webclients.WebClients[session]

	if err != ErrActiveSession {
		t.Errorf("Should have returned an error for  but returned [ %v ]\n", err)
	}
	if registered != webclient {
		t.Errorf("Should have registered webclient [ %v ] but registered [ %v ]\n", webclient, registered)
	}
}

func Test_Should_register_webclient_for_valid_session(t *testing.T) {
	webclients := NewWebClientSessionsView()
	session := "test"
	var expected uint64 = 1

	err := webclients.RegisterWebClientBySession(session, expected)

	actual := webclients.WebClients[session]

	if err != nil {
		t.Errorf("Should not have returned an error but returned [ %v ]\n", err)
	}
	if actual != expected {
		t.Errorf("Should have registered webclient [ %v ] but registered [ %v ]\n", expected, actual)
	}
}

func Test_Should_return_error_for_invalid_session(t *testing.T) {
	webclients := NewWebClientSessionsView()
	session := "test"

	actual, err := webclients.GetBySession(session)

	if err != ErrInvalidSession {
		t.Errorf("Should have returned an error for invalid session but was [ %v ]\n", err)
	}
	if actual != 0 {
		t.Errorf("Should have returned nil webclient but returned [ %v ]\n", actual)
	}
}

func Test_Should_return_webclient_for_valid_session(t *testing.T) {
	webclients := NewWebClientSessionsView()
	var webclient uint64 = 1
	session := "test"
	webclients.WebClients[session] = webclient

	var expected uint64 = 1

	actual, err := webclients.GetBySession(session)

	if err != nil {
		t.Errorf("Should not have returned an error but returned [ %v ]\n", err)
	}
	if actual != expected {
		t.Errorf("Should have returned webclient [ %v ] but returned [ %v ]\n", expected, actual)
	}
}

func Test_Should_return_error_for_invalid_webclient(t *testing.T) {
	webclients := NewWebClientSessionsView()
	var webclient uint64 = 1

	err := webclients.DeleteByWebClient(webclient)

	if err != ErrInactiveWebClient {
		t.Errorf("Should not have returned an error but returned [ %v ]\n", err)
	}
}

func Test_Should_delete_webclient_for_valid_session(t *testing.T) {
	webclients := NewWebClientSessionsView()
	var webclient uint64 = 1
	session := "test"
	webclients.WebClients[session] = webclient

	err := webclients.DeleteByWebClient(webclient)
	_, active := webclients.WebClients[session]

	if err != nil {
		t.Errorf("Should not have returned an error but returned [ %v ]\n", err)
	}
	if active {
		t.Errorf("Should have deleted webclient by session [ %v ] but failed\n", session)
	}
}

func Test_Should_add_webclient_to_view(t *testing.T) {
	webclients := NewWebClientSessionsView()
	var webclient uint64 = 1
	session := "session"

	event := NewWebClientAdded(webclient, session)
	var expected uint64 = 1

	webclients.HandleEvent(event)
	actual, err := webclients.GetBySession(session)

	if err != nil {
		t.Errorf("Should not have returned an error but returned [ %v ]\n", err)
	}
	if actual != expected {
		t.Errorf("Should have added webclient [ %v ] but returned [ %v ]\n", expected, actual)
	}
}

func Test_Should_overwrite_webclient_in_view(t *testing.T) {
	webclients := NewWebClientSessionsView()
	webclients.RegisterWebClientBySession("session", 2)
	var webclient uint64 = 1
	session := "session"

	event := NewWebClientAdded(webclient, session)
	var expected uint64 = 1

	webclients.HandleEvent(event)
	actual, err := webclients.GetBySession(session)

	if err != nil {
		t.Errorf("Should not have returned an error but returned [ %v ]\n", err)
	}
	if actual != expected {
		t.Errorf("Should have returned webclient [ %v ] but returned [ %v ]\n", expected, actual)
	}
}

func Test_Should_update_webclient_in_view(t *testing.T) {
	webclients := NewWebClientSessionsView()
	var previous uint64 = 2
	var webclient uint64 = 1
	session := "session"

	webclients.RegisterWebClientBySession(session, previous)
	event := NewWebClientSessionUpdated(webclient, 0, session)
	var expected uint64 = 1

	webclients.HandleEvent(event)
	actual, err := webclients.GetBySession(session)

	if err != nil {
		t.Errorf("Should not have returned an error but returned [ %v ]\n", err)
	}
	if actual != expected {
		t.Errorf("Should have returned webclient [ %v ] but returned [ %v ]\n", expected, actual)
	}
}

func Test_Should_remove_webclient_in_view(t *testing.T) {
	webclients := NewWebClientSessionsView()
	var webclient uint64 = 1
	var expected uint64 = 0
	var session = "session"

	webclients.RegisterWebClientBySession(session, webclient)
	event := NewWebClientRemoved(webclient, 0, session)

	webclients.HandleEvent(event)
	actual, err := webclients.GetBySession(session)

	if err != ErrInvalidSession {
		t.Errorf("Should not have returned an error but returned [ %v ]\n", err)
	}

	if actual != expected {
		t.Errorf("Should have returned webclient [ %v ] but returned [ %v ]\n", expected, actual)
	}
}
