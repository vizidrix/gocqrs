package eventstore_test

import (
	//"fmt"
	"testing"
	"os"
	"github.com/vizidrix/gocqrs/cqrs"
	"github.com/vizidrix/gocqrs/eventstore"
)

var filename string = "test.db"
var fileperms os.FileMode = 0600

func makeEventStore(bufferSize int) eventstore.EventStore {
	es, err := eventstore.NewBoltES(filename, fileperms, func() chan cqrs.Event {
		return make(chan cqrs.Event, bufferSize)
	})
	if err != nil {
		panic("Unable to make event store")
	}
	return es
}

func cleanEventStore() {
	os.Remove(filename)
}

var (
	TEST_DOMAIN   uint32 = 0x11111111
	B_TestEvent   uint64 = cqrs.E(TEST_DOMAIN, 1, 1)
)

type TestEvent struct {
	cqrs.EventMemento
	Value string
}

func NewTestEvent(id uint64, version uint32, value string) TestEvent {
	return TestEvent{
		EventMemento: cqrs.NewEvent(id, version, B_TestEvent),
		Value:        value,
	}
}

func Test_Should_make_db(t *testing.T) {
	defer cleanEventStore()
	if _, err := eventstore.NewBoltES(filename, fileperms, func() chan cqrs.Event { return nil }); err != nil {
		t.Errorf("Error creating event store [ %s ]\n", err)
	}
}

func Test_Should_put_event(t *testing.T) {
	es := makeEventStore(0)
	defer cleanEventStore()

	event := NewTestEvent(10, 1, "stuff")

	if err := es.Put(event); err != nil {
		t.Errorf("Error putting event in store [ %s ]\n", err)
	}
}

func Test_Should_not_error_getting_empty_stream(t *testing.T) {
	es := makeEventStore(0)
	defer cleanEventStore()

	if _, err := es.GetStreamById(TEST_DOMAIN, 10); err != nil {
		t.Errorf("Error loading empty stream [ %s ]\n", err)
	}
}

func Test_Should_return_zero_events_from_empty_eventstore_by_id(t *testing.T) {
	es := makeEventStore(0)
	defer cleanEventStore()

	stream, _ := es.GetStreamById(TEST_DOMAIN, 10)
	if len(stream) != 0 {
		t.Errorf("Event stream should have zero entries for empty event store by id\n")
	}
}

func Test_Should_return_single_event_by_id(t *testing.T) {
	es := makeEventStore(0)
	defer cleanEventStore()

	expected := NewTestEvent(10, 1, "stuff")
	es.Put(expected)

	stream, _ := es.GetStreamById(TEST_DOMAIN, 10)
	if len(stream) != 1 {
		t.Errorf("Event stream should have one entry for id but len was [ %d ]\n", len(stream))
	}
}

func Test_Should_return_zero_events_from_empty_eventstore_by_domain(t *testing.T) {
	es := makeEventStore(0)
	defer cleanEventStore()

	stream, _ := es.GetStreamByDomain(TEST_DOMAIN)
	if len(stream) != 0 {
		t.Errorf("Event stream should have zero entries for empty event store by domain")
	}
}

func Test_Should_return_single_event_by_domain(t *testing.T) {
	es := makeEventStore(0)
	defer cleanEventStore()

	expected := NewTestEvent(10, 1, "stuff")
	es.Put(expected)

	stream, _ := es.GetStreamByDomain(TEST_DOMAIN)
	if len(stream) != 1 {
		t.Errorf("Event stream should have one entry for domain but len was [ %d ]\n", len(stream))
	}
}

/*
db, err := es.OpenEventStore("test.db", 0600)
if err != nil {
	return
}
defer db.Close()

db.Update(func(tx *bolt.Tx) error {
	bucket := tx.Bucket([]byte("events"))
	bucket.Put([]byte("key"), []byte("values"))
	return nil
	// tx.Bucket([]byte("bucketname")).Put([]byte("name"))
	})

db.View(func(tx *bolt.Tx) error {
	bucket := tx.Bucket([]byte("events"))
	data := bucket.Get([]byte("key"))
	fmt.Printf("Data: [ %s ]\n", data)
	return nil
	})
//data := bucket.Get([]byte("value"))
//fmt.Printf("Data: [ %s ]\n", data)
*/

