package gocqrs_test

import (
	//cqrs "github.com/vizidrix/gocqrs"
	example "github.com/vizidrix/gocqrs/example"
	. "github.com/vizidrix/gocqrs/test_utils"
	"log"
	"testing"
	//"time"
)

func ignore() { log.Println("") }

// Publish Command
// - Store Command
// Command Handler
// - Hydrate Aggregate
// Publish Event
// - Store Event
// Event Handler(s)
// - Update View(s)

func Test_Should_return_nil_if_no_events_found(t *testing.T) {
	aggregate := AggregateTest(t).Given(example.LoadPerson).With().When().ThenFailOnError()

	person := (<-aggregate).(*example.Person)

	if person != nil {
		t.Fail()
	}
}

func Test_Should_load_registered_person(t *testing.T) {
	aggregate := AggregateTest(t).Given(example.LoadPerson).With().When(
		example.NewPersonRegistered(1, 1, "John", "Wayne", 987654321, "First"),
	).ThenFailOnError()

	person := (<-aggregate).(*example.Person)

	if person.Profile != "First" {
		log.Printf("Profile not set: %s", person.Profile)
		t.Fail()
	}
	if person.Aggregate.GetVersion() != 1 {
		log.Printf("Version not updated: %s", person.Aggregate.GetVersion())
		t.Fail()
	}
}

func Test_Should_update_profile_when_NewProfileUpdated_event_received(t *testing.T) {
	aggregate := AggregateTest(t).Given(example.LoadPerson).With(
		example.NewPersonRegistered(1, 1, "John", "Wayne", 987654321, "First"),
	).When(
		example.NewProfileUpdated(1, 2, "Stuff"),
	).ThenFailOnError()

	person := (<-aggregate).(*example.Person)

	if person.Profile != "Stuff" {
		log.Printf("Profile not updated: %s", person.Profile)
		t.Fail()
	}
	if person.Aggregate.GetVersion() != 2 {
		log.Printf("Version not updated: %s", person.Aggregate.GetVersion())
		t.Fail()
	}
}

func Test_Should_return_err_if_nil_is_sent(t *testing.T) {
	result, err := AggregateTest(t).Given(example.LoadPerson).With(
		example.NewPersonRegistered(1, 1, "John", "Wayne", 987654321, "First"),
	).When(
		nil,
	).Then()

	select {
	case p := <-result:
		{
			log.Printf("Should have sent error due to nil on event channel", p)
			t.Fail()
		}
	case <-err:
		{
		}
	}
}

func Test_Should_return_err_if_PersonRegistered_twice(t *testing.T) {
	result, err := AggregateTest(t).Given(example.LoadPerson).With(
		example.NewPersonRegistered(1, 1, "John", "Wayne", 987654321, "First"),
	).When(
		example.NewPersonRegistered(1, 1, "Jack", "Johnson", 987654321, "JJ"),
	).Then()

	select {
	case p := <-result:
		{
			log.Printf("Should have sent error due to nil on event channel", p)
			t.Fail()
		}
	case <-err:
		{
		}
	}
}

func Test_Should_return_err_if_first_event_version_out_of_sync(t *testing.T) {
	result, err := AggregateTest(t).Given(example.LoadPerson).With().When(
		example.NewPersonRegistered(1, 2, "Jack", "Johnson", 987654321, "JJ"),
	).Then()

	select {
	case p := <-result:
		{
			log.Printf("Should have sent error due to event version being wrong", p)
			t.Fail()
		}
	case <-err:
		{
		}
	}
}

func Test_Should_return_err_if_new_event_version_out_of_sync(t *testing.T) {
	result, err := AggregateTest(t).Given(example.LoadPerson).With(
		example.NewPersonRegistered(1, 1, "John", "Wayne", 987654321, "First"),
		example.NewProfileUpdated(1, 2, "Stuff"),
		example.NewProfileUpdated(1, 3, "Stuff"),
	).When(
		example.NewProfileUpdated(1, 3, "Stuff"),
	).Then()

	select {
	case p := <-result:
		{
			log.Printf("Should have sent error due to event version being wrong", p)
			t.Fail()
		}
	case <-err:
		{
		}
	}
}
