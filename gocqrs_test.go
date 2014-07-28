package gocqrs_test

import (
	"testing"
	"runtime"
	"fmt"
	"reflect"
	"path/filepath"
	cqrs "github.com/vizidrix/gocqrs"
)

// Helper functions for test suites
// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func Test_Should_mask_valid_command_type_with_version(t *testing.T) {
	// Arrange
	var version uint8 = 1
	var typeID uint32 = 1
	var expected uint32 = 0x81000001

	// Act
	actual := cqrs.MakeVersionedCommandType(version, typeID)

	// Assert
	assert(t, expected == actual, "masked event type id expected [ %0x ] but was [ %0x ]", expected, actual)
}

func Test_Should_mask_valid_event_type_with_version(t *testing.T) {
	// Arrange
	var version uint8 = 1
	var typeID uint32 = 1
	var expected uint32 = 0x01000001

	// Act
	actual := cqrs.MakeVersionedEventType(version, typeID)

	// Assert
	assert(t, expected == actual, "masked event type id expected [ %0x ] but was [ %0x ]", expected, actual)
}

func Test_Should_populate_aggregate_from_New(t *testing.T) {
	// Arrange
	var app uint32 = 10
	var domain uint32 = 11
	var id uint64= 12
	var version uint32 = 13

	// Act
	aggregate := cqrs.NewAggregate(app, domain, id, version)

	// Assert
	assert(t, aggregate.GetApplication() == app, "aggregate app not set")
	assert(t, aggregate.GetDomain() == domain, "aggregate domain not set")
	assert(t, aggregate.GetId() == id, "aggregate ID not set")
	assert(t, aggregate.GetVersion() == version, "aggregate version not set")
}

func Test_Should_populate_command_from_New(t *testing.T) {
	// Arrange
	var app uint32 = 10
	var domain uint32 = 11
	var id uint64 = 12
	var version uint32 = 13
	var commandType uint32 = 14

	// Act
	command := cqrs.NewCommand(app, domain, id, version, commandType)

	// Assert
	assert(t, command.GetApplication() == app, "command app not set")
	assert(t, command.GetDomain() == domain, "command domain not set")
	assert(t, command.GetId() == id, "command id not set")
	assert(t, command.GetVersion() == version, "command version not set")
	assert(t, command.GetCommandType() == commandType, "command type not set")
}

func Test_Should_populate_event_from_New(t *testing.T) {
	// Arrange
	var app uint32 = 10
	var domain uint32 = 11
	var id uint64 = 12
	var version uint32 = 13
	var eventType uint32 = 14

	// Act
	event := cqrs.NewEvent(app, domain, id, version, eventType)

	// Assert
	assert(t, event.GetApplication() == app, "event app not set")
	assert(t, event.GetDomain() == domain, "event domain not set")
	assert(t, event.GetId() == id, "event id not set")
	assert(t, event.GetVersion() == version, "event version not set")
	assert(t, event.GetEventType() == eventType, "event type not set")
}


/*
// Mock Event Store

type mockEventStore struct {
	eventstore.EventStoreReaderWriter
	generateAggregateId        func(uint32, uint32) (uint64, error)
	appendEvent                func(eventstore.Event) (time.Time, error)
	loadEventStreamByAggregate func(uint32, uint32, uint64) ([]eventstore.Event, error)
	loadEventStreamByDomain    func(uint32, uint32) ([]eventstore.Event, error)
}

func newMockEventStore() *mockEventStore {
	return &mockEventStore{
		generateAggregateId:        func(uint32, uint32) (uint64, error) { return 0, nil },
		appendEvent:                func(eventstore.Event) (time.Time, error) { return time.Time{}, nil },
		loadEventStreamByAggregate: func(uint32, uint32, uint64) ([]eventstore.Event, error) { return make([]eventstore.Event, 0), nil },
		loadEventStreamByDomain:    func(uint32, uint32) ([]eventstore.Event, error) { return make([]eventstore.Event, 0), nil },
	}
}

func (es *mockEventStore) GenerateAggregateId(application uint32, domain uint32) (uint64, error) {
	return es.generateAggregateId(application, domain)
}

func (es *mockEventStore) AppendEvent(event eventstore.Event) (time.Time, error) {
	return es.appendEvent(event)
}

func (es *mockEventStore) LoadEventStreamByAggregate(application uint32, domain uint32, id uint64) ([]eventstore.Event, error) {
	return es.loadEventStreamByAggregate(application, domain, id)
}

func (es *mockEventStore) LoadEventStreamByDomain(application uint32, domain uint32) ([]eventstore.Event, error) {
	return es.loadEventStreamByDomain(application, domain)
}

func (es *mockEventStore) mockGenerateAggregateId(id uint64, err error) {
	es.generateAggregateId = func(uint32, uint32) (uint64, error) { return id, err }
}

func (es *mockEventStore) mockAppendEvent(timestamp time.Time, err error) {
	es.appendEvent = func(eventstore.Event) (time.Time, error) { return timestamp, err }
}

func (es *mockEventStore) mockLoadEventStreamByAggregate(events []eventstore.Event, err error) {
	es.loadEventStreamByAggregate = func(uint32, uint32, uint64) ([]eventstore.Event, error) { return events, err }
}

func (es *mockEventStore) mockLoadEventStreamByDomain(events []eventstore.Event, err error) {
	es.loadEventStreamByDomain = func(uint32, uint32) ([]eventstore.Event, error) { return events, err }
}

// Mock Hydrator

type mockHydrator struct {
	cqrs.Hydrator
	hydrate func([]eventstore.Event) (eventstore.Aggregate, error)
}

func newMockHydrator() *mockHydrator {
	return &mockHydrator{
		hydrate: func([]eventstore.Event) (eventstore.Aggregate, error) { return nil, nil },
	}
}

func (h *mockHydrator) Hydrate(events []eventstore.Event) (eventstore.Aggregate, error) {
	return h.hydrate(events)
}

func (h *mockHydrator) mockHydrate(aggregate eventstore.Aggregate, err error) {
	h.hydrate = func([]eventstore.Event) (eventstore.Aggregate, error) { return aggregate, err }
}

// Mock Commander

type mockCommander struct {
	cqrs.Commander
	runCommand func(eventstore.Aggregate, cqrs.Command) (eventstore.Event, error)
}

func newMockCommander() *mockCommander {
	return &mockCommander{
		runCommand: func(eventstore.Aggregate, cqrs.Command) (eventstore.Event, error) { return nil, nil },
	}
}

func (c *mockCommander) RunCommand(state eventstore.Aggregate, command cqrs.Command) (eventstore.Event, error) {
	return c.runCommand(state, command)
}

func (c *mockCommander) mockRunCommand(event eventstore.Event, err error) {
	c.runCommand = func(eventstore.Aggregate, cqrs.Command) (eventstore.Event, error) { return event, err }
}

//
// Testing Suites
//

//
// Test LoadAggregateState() method
//

func Test_Should_mock_valid_state_from_load_aggregate_state(t *testing.T) {
	eventStore := newMockEventStore()
	hydrator := newMockHydrator()
	commander := newMockCommander()
	commandhandler := cqrs.NewCommandHandler(
		eventStore, hydrator, commander, CQRS_TEST_APPLICATION, CQRS_TEST_DOMAIN,
	)
	eventStore.mockLoadEventStreamByAggregate([]eventstore.Event{NewTestEvent(1, 1)}, nil)
	hydrator.mockHydrate(NewTestAggregate(1, 1), nil)

	state, err := commandhandler.LoadAggregateState(1)
	ok(t, err)
	assert(t, state != nil, "returned nil valued aggregate state")
}

func Test_Should_mock_nil_state_unhydrated_from_load_aggregate_state(t *testing.T) {
	eventStore := newMockEventStore()
	hydrator := newMockHydrator()
	commander := newMockCommander()
	commandhandler := cqrs.NewCommandHandler(
		eventStore, hydrator, commander, CQRS_TEST_APPLICATION, CQRS_TEST_DOMAIN,
	)
	eventStore.mockLoadEventStreamByAggregate([]eventstore.Event{NewTestEvent(1, 1)}, nil)

	state, err := commandhandler.LoadAggregateState(1)
	ok(t, err)
	assert(t, state == nil, "returned non-nil valued aggregate state")
}

func Test_Should_mock_error_loading_aggregate_events_with_nil_state_from_load_aggregate_state(t *testing.T) {
	eventStore := newMockEventStore()
	hydrator := newMockHydrator()
	commander := newMockCommander()
	commandhandler := cqrs.NewCommandHandler(
		eventStore, hydrator, commander, CQRS_TEST_APPLICATION, CQRS_TEST_DOMAIN,
	)
	loadError := errors.New("error loading aggregate state")
	eventStore.mockLoadEventStreamByAggregate([]eventstore.Event{NewTestEvent(1, 1)}, loadError)

	state, err := commandhandler.LoadAggregateState(1)
	assert(t, err == loadError, "wrong load aggregate state error: %v", err)
	assert(t, state == nil, "returned non-nil valued aggregate state")
}

func Test_Should_mock_error_hydrating_state_from_load_aggregate_state(t *testing.T) {
	eventStore := newMockEventStore()
	hydrator := newMockHydrator()
	commander := newMockCommander()
	commandhandler := cqrs.NewCommandHandler(
		eventStore, hydrator, commander, CQRS_TEST_APPLICATION, CQRS_TEST_DOMAIN,
	)
	hydrateError := errors.New("error loading aggregate state")
	hydrator.mockHydrate(nil, hydrateError)

	_, err := commandhandler.LoadAggregateState(1)
	assert(t, err == hydrateError, "wrong load aggregate state error: %v", err)
}

//
// Test ValidateCommand() method
//

func Test_Should_return_true_from_validate_command(t *testing.T) {
	eventStore := newMockEventStore()
	hydrator := newMockHydrator()
	commander := newMockCommander()
	commandhandler := cqrs.NewCommandHandler(
		eventStore, hydrator, commander, CQRS_TEST_APPLICATION, CQRS_TEST_DOMAIN,
	)

	validation := commandhandler.ValidateCommand(NewTestCommand(1, 1))
	assert(t, validation == true, "did not validate correctly")
}

func Test_Should_return_false_from_validate_command(t *testing.T) {
	eventStore := newMockEventStore()
	hydrator := newMockHydrator()
	commander := newMockCommander()
	commandhandler := cqrs.NewCommandHandler(
		eventStore, hydrator, commander, CQRS_TEST_APPLICATION, CQRS_TEST_DOMAIN,
	)

	validation := commandhandler.ValidateCommand(cqrs.NewCommand(1, 1, 1, 1))
	assert(t, validation == false, "did not validate correctly")
}

//
// Test Publish() method
//

func Test_Should_mock_valid_timestamp_from_valid_publish(t *testing.T) {
	eventStore := newMockEventStore()
	hydrator := newMockHydrator()
	commander := newMockCommander()
	commandhandler := cqrs.NewCommandHandler(
		eventStore, hydrator, commander, CQRS_TEST_APPLICATION, CQRS_TEST_DOMAIN,
	)
	eventStore.mockAppendEvent(time.Now().UTC(), nil)

	timestamp, err := commandhandler.Publish(NewTestEvent(1, 1))
	ok(t, err)
	assert(t, timestamp != time.Time{}, "return nil valued timestamp")
}

func Test_Should_mock_error_from_publish(t *testing.T) {
	eventStore := newMockEventStore()
	hydrator := newMockHydrator()
	commander := newMockCommander()
	commandhandler := cqrs.NewCommandHandler(
		eventStore, hydrator, commander, CQRS_TEST_APPLICATION, CQRS_TEST_DOMAIN,
	)
	publishError := errors.New("error publishing event")
	eventStore.mockAppendEvent(time.Time{}, publishError)

	_, err := commandhandler.Publish(NewTestEvent(1, 1))
	assert(t, err == publishError, "wrong publish event error: %v", err)
}

//
// Test HandleCommand() method
//

func Test_Should_mock_valid_timestamp_from_handle_command(t *testing.T) {
	eventStore := newMockEventStore()
	hydrator := newMockHydrator()
	commander := newMockCommander()
	commandhandler := cqrs.NewCommandHandler(
		eventStore, hydrator, commander, CQRS_TEST_APPLICATION, CQRS_TEST_DOMAIN,
	)
	eventStore.mockAppendEvent(time.Now().UTC(), nil)

	timestamp, err := commandhandler.HandleCommand(NewTestCommand(1, 1))
	ok(t, err)
	assert(t, timestamp != time.Time{}, "return nil valued timestamp")
}

func Test_Should_return_ErrWrongDomain_with_nil_valued_timestamp_from_handle_command(t *testing.T) {
	eventStore := newMockEventStore()
	hydrator := newMockHydrator()
	commander := newMockCommander()
	commandhandler := cqrs.NewCommandHandler(
		eventStore, hydrator, commander, CQRS_TEST_APPLICATION, CQRS_TEST_DOMAIN,
	)
	eventStore.mockAppendEvent(time.Now().UTC(), nil)

	timestamp, err := commandhandler.HandleCommand(cqrs.NewCommand(1, 1, 1, 1))
	assert(t, err == cqrs.ErrWrongDomain, "wrong handle command error: %v", err)
	assert(t, timestamp == time.Time{}, "returned non-nil valued timestamp")
}

func Test_Should_mock_error_loading_aggregate_events_with_nil_valued_timestamp_from_handle_command(t *testing.T) {
	eventStore := newMockEventStore()
	hydrator := newMockHydrator()
	commander := newMockCommander()
	commandhandler := cqrs.NewCommandHandler(
		eventStore, hydrator, commander, CQRS_TEST_APPLICATION, CQRS_TEST_DOMAIN,
	)
	loadError := errors.New("error loading aggregate events")
	eventStore.mockLoadEventStreamByAggregate([]eventstore.Event{NewTestEvent(1, 1)}, loadError)
	eventStore.mockAppendEvent(time.Now().UTC(), errors.New("invalid error"))

	timestamp, err := commandhandler.HandleCommand(NewTestCommand(1, 1))
	assert(t, err == loadError, "wrong handle command error: %v", err)
	assert(t, timestamp == time.Time{}, "returned non-nil valued timestamp")
}

func Test_Should_mock_error_hydrating_state_with_nil_valued_timestamp_from_handle_command(t *testing.T) {
	eventStore := newMockEventStore()
	hydrator := newMockHydrator()
	commander := newMockCommander()
	commandhandler := cqrs.NewCommandHandler(
		eventStore, hydrator, commander, CQRS_TEST_APPLICATION, CQRS_TEST_DOMAIN,
	)
	hydrateError := errors.New("error hydrating aggregate state")
	hydrator.mockHydrate(nil, hydrateError)
	eventStore.mockAppendEvent(time.Now().UTC(), errors.New("invalid error"))

	timestamp, err := commandhandler.HandleCommand(NewTestCommand(1, 1))
	assert(t, err == hydrateError, "wrong handle command error: %v", err)
	assert(t, timestamp == time.Time{}, "returned non-nil valued timestamp")
}

func Test_Should_mock_error_running_command_with_nil_valued_timestamp_from_handle_command(t *testing.T) {
	eventStore := newMockEventStore()
	hydrator := newMockHydrator()
	commander := newMockCommander()
	commandhandler := cqrs.NewCommandHandler(
		eventStore, hydrator, commander, CQRS_TEST_APPLICATION, CQRS_TEST_DOMAIN,
	)
	commandError := errors.New("error hydrating aggregate state")
	commander.mockRunCommand(nil, commandError)
	eventStore.mockAppendEvent(time.Now().UTC(), errors.New("invalid error"))

	timestamp, err := commandhandler.HandleCommand(NewTestCommand(1, 1))
	assert(t, err == commandError, "wrong handle command error: %v", err)
	assert(t, timestamp == time.Time{}, "returned non-nil valued timestamp")
}
*/