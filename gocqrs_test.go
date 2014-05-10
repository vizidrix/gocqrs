package gocqrs_test

import (
	cqrs "github.com/vizidrix/gocqrs"
	//example "github.com/vizidrix/gocqrs/example"
	//. "github.com/vizidrix/gocqrs/test_utils"
	//"log"
	"testing"
	"fmt"
	//"time"
	//"errors"
	"encoding/json"
)

func HandleJson(data []byte) error {
	cmd := &cqrs.Command{}
	if err := json.Unmarshal(data, cmd); err != nil {
		return err
	}
	switch cmd.Type {
		case 1: {
			_cmd := &DoSomethingCommand{}
			if err := json.Unmarshal(data, _cmd); err != nil {
				return err
			}
			HandleDoSomethingCommand(_cmd)
		}
		case 2: {
			_cmd := &DoAnotherThingCommand{}
			if err := json.Unmarshal(data, _cmd); err != nil {
				return err
			}
			HandleDoAnotherThingCommand(_cmd)
		}
	}
	return nil
}

type DoSomethingCommand struct {
	cqrs.Command
	StringValue string `json:"stringalue"`
	IntValue int64 `json:"intvalue"`
}

func NewDoSomethingCommand(stringValue string, intValue int64) *DoSomethingCommand {
	return &DoSomethingCommand {
		cqrs.Command { Type: 1 },
		stringValue,
		intValue,
	}
}

func HandleDoSomethingCommand(cmd *DoSomethingCommand) {
	if cmd.StringValue == "" {
		fmt.Printf("Empty Something\n")
	} else {
		fmt.Printf("Got something: [ %s ]\n", cmd.StringValue)
	}
}

type DoAnotherThingCommand struct {
	cqrs.Command
	StringValue string `json:"stringvalue"`
	IntValue int64 `json:"intvalue"`
}

func NewDoAnotherThingCommand(stringValue string, intValue int64) *DoAnotherThingCommand {
	return &DoAnotherThingCommand {
		cqrs.Command { Type: 2 },
		stringValue,
		intValue,
	}
}

func HandleDoAnotherThingCommand(cmd *DoAnotherThingCommand) {
	if cmd.StringValue == "" {
		fmt.Printf("Another empty\n")
	} else {
		fmt.Printf("Got Another [ %s ]\n", cmd.StringValue)
	}
}

func Test_Should_create_command(t *testing.T) {
	/*cmd := &TestCommandOneCommand {
		cqrs.Command { Type: 10 },
		"one",
		1,
	}*/
	cmd1 := NewDoSomethingCommand("one", 1)
	cmd2 := NewDoAnotherThingCommand("two", 2)
	HandleDoSomethingCommand(cmd1)
	HandleDoAnotherThingCommand(cmd2)
	
	data, err := json.Marshal(cmd1)
	if err != nil {
		fmt.Printf("JSON Err: [ %s ]\n", err)
	} else {
		fmt.Printf("JSON [ %s ]\n", data)
	}
	cmd1_ := &DoSomethingCommand{}
	err = json.Unmarshal(data, cmd1_)
	if err != nil {
		fmt.Printf("JSON Un Err: [ %s ]\n", err)
	} else {
		fmt.Printf("Unmarshalled: [ %v ]\n", cmd1_)
	}
	base := &cqrs.Command{}
	err = json.Unmarshal(data, base)
	if err != nil {
		fmt.Printf("JSON Base Err: [ %s ]\n", err)
	} else {
		fmt.Printf("Base: [ %v ]\n", base)
	}
	fmt.Printf("cmd:\n\t%v\n\t%v\n", cmd1, cmd2)

	HandleJson(data)

}

//func ignore() { log.Println("") }

// Publish Command
// - Store Command
// Command Handler
// - Hydrate Aggregate
// Publish Event
// - Store Event
// Event Handler(s)
// - Update View(s)

/*
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
*/
