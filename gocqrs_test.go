package gocqrs_test

import (
	cqrs "github.com/vizidrix/gocqrs"
	example "github.com/vizidrix/gocqrs/example"
	"log"
	"testing"
	//"time"
)

func ignore() {
	log.Println("")
}

// Publish Command
// - Store Command
// Command Handler
// - Hydrate Aggregate
// Publish Event
// - Store Event
// Event Handler(s)
// - Update View(s)

func Test_Should_load_aggregate_from_event_channel_sync_when_channel_closes_immediately(t *testing.T) {
	// Arrange
	eventChannel := make(chan cqrs.IEvent)
	close(eventChannel)

	// Act
	if person, err := example.LoadPersonSync(eventChannel); err != nil {
		log.Printf("Error: %s", err)
		t.Fail()
	} else {
		log.Printf("Person: %s", person)
		if person != nil {
			t.Fail()
		}
	}
}

func Test_Should_return_err_if_nil_is_sent(t *testing.T) {
	// Arrange
	eventChannel := make(chan cqrs.IEvent, 1)
	eventChannel <- nil
	close(eventChannel)

	// Act
	if person, err := example.LoadPersonSync(eventChannel); err == nil {
		log.Printf("Should have returned err: %s, %s", person, err)
		t.Fail()
	}
}

func Test_Should_load_aggregate_from_event_channel(t *testing.T) {
	// Arrange
	version := int64(1)
	eventChannel := make(chan cqrs.IEvent, 1)
	eventChannel <- example.NewPersonRegistered(1, version, "John", "Wayne", 987654321, "")
	close(eventChannel)

	// Act
	if person, err := example.LoadPersonSync(eventChannel); err != nil {
		log.Printf("Error: %s", err)
		t.Fail()
	} else {
		log.Printf("Person: %s", person)
		if person == nil {
			t.Fail()
		}
		if person.Aggregate.GetVersion() != 1 {
			log.Printf("Version not updated: %s", person.Aggregate.GetVersion())
			t.Fail()
		}
	}
}

func Test_Should_update_profile(t *testing.T) {
	// Arrange
	version := int64(1)
	eventChannel := make(chan cqrs.IEvent, 2)
	eventChannel <- example.NewPersonRegistered(1, version, "John", "Wayne", 987654321, "")
	version++
	eventChannel <- example.NewProfileUpdated(1, version, "Stuff")
	close(eventChannel)

	// Act
	if person, err := example.LoadPersonSync(eventChannel); err != nil {
		log.Printf("Error: %s", err)
		t.Fail()
	} else {
		log.Printf("Person: %s", person)
		if person == nil {
			t.Fail()
		}
		if person.Profile != "Stuff" {
			log.Printf("Profile not updated: %s", person.Profile)
			t.Fail()
		}
		if person.Aggregate.GetVersion() != 2 {
			log.Printf("Version not updated: %s", person.Aggregate.GetVersion())
			t.Fail()
		}
	}
}

func Test_Given_an_empty_memory_event_store(t *testing.T) {
	t.Skipf("Skipping", "because")
	// Arrange
	//commands := GetCommandPublisher()
	//eventStore := cqrs.NewMemoryEventStore()

	//eventSet := eventStore.Of("Person").WithId(1)

	// Act
	//count := <-eventSet.Count()

	//for event := range eventSet.All() {

	// Assert
	//if count != 0 {
	//	t.Fail()
	//}

	//log.Printf("Aggregate Store: %s", aggregateStore)
	//log.Printf("Event Set: %s", eventSet)
	// Act
	//t.Fail()
}
