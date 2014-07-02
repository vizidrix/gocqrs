package clientsockets

import (
	"testing"
)

var (
	testsession string = "test"
	testclient  uint64 = 1
)

func NewTestConnectionService() ConnectionService {
	return ConnectionService{
		connections:      make(map[uint64]*ConnectionMemento),
		addChan:          make(chan *ConnectionMemento, 1),
		removeChan:       make(chan *ConnectionMemento, 1),
		subscriptionChan: make(chan ClientConnection, 1),
	}
}

func Test_Should_add_connection_to_connection_service(t *testing.T) {
	var testservice = NewTestConnectionService()
	var testconn = NewConnection(testsession, testclient)
	var expected = &testconn

	AddConnection(&testservice, &testconn)
	actual, active := testservice.connections[testclient]

	if !active {
		t.Errorf("Should have added connection for client [ %v ] but failed\n", testclient)
	}
	if actual != expected {
		t.Errorf("Should have added connection for client [ %v ] but instead added for client [ %v ]\n", expected, actual)
	}
}

func Test_Should_add_connection_to_subscription_channel(t *testing.T) {
	var testservice = NewTestConnectionService()
	var testconn = NewConnection(testsession, testclient)
	var expected = &testconn

	AddConnection(&testservice, &testconn)

	actual := <-testservice.subscriptionChan

	if actual != expected {
		t.Errorf("Should have subscribed connection [ %v ] but instead subscribed connection [ %v ]\n", expected, actual)
	}
}

func Test_Should_not_add_connection_to_connection_service(t *testing.T) {
	var testservice = NewTestConnectionService()
	var startconn = NewConnection(testsession, testclient)
	var testconn = NewConnection(testsession, testclient)
	var expected = &startconn

	testservice.connections[testclient] = &startconn
	AddConnection(&testservice, &testconn)
	actual, active := testservice.connections[testclient]

	if !active {
		t.Errorf("Should have retained connection for client [ %v ] but failed\n", testclient)
	}
	if actual != expected {
		t.Errorf("Should not have retained connection [ %v ] but changed to connection [ %v ]\n", expected, actual)
	}
}

func Test_Should_add_connection_to_remove_channel_instead_of_subscription_channel(t *testing.T) {
	var testservice = NewTestConnectionService()
	var startconn = NewConnection(testsession, testclient)
	var testconn = NewConnection(testsession, testclient)
	var expected = &startconn

	testservice.connections[testclient] = &startconn
	AddConnection(&testservice, &testconn)

	select {
	case <-testservice.subscriptionChan:
		t.Errorf("Should not have added connection to subscription channel\n")
	case actual := <-testservice.removeChan:
		if actual != expected {
			t.Errorf("Should added connection [ %v ] to remove channel but instead added connection [ %v ]\n", expected, actual)
		}
	}
}

func Test_Should_remove_connection_from_connection_service(t *testing.T) {
	var testservice = NewTestConnectionService()
	var testconn = NewConnection(testsession, testclient)

	RemoveConnection(&testservice, &testconn)
	_, active := testservice.connections[testclient]

	if active {
		t.Errorf("Should have removed connection for client [ %v ] but failed\n", testclient)
	}
}

func Test_Should_handle_redundant_removals(t *testing.T) {
	var testservice = NewTestConnectionService()
	var testconn = NewConnection(testsession, testclient)

	RemoveConnection(&testservice, &testconn)
	RemoveConnection(&testservice, &testconn)
	_, active := testservice.connections[testclient]

	if active {
		t.Errorf("Should have removed connection for client [ %v ] but failed\n", testclient)
	}
}
