package gocqrs_test

import (
	"github.com/vizidrix/gocqrs/cqrs"
	"github.com/vizidrix/gocqrs/web"

	//example "github.com/vizidrix/gocqrs/example"
//	"github.com/vizidrix/gocqrs/web"
	//. "github.com/vizidrix/gocqrs/test_utils"
	//"log"
	"testing"
	"fmt"
	//"time"
	//"errors"
	//"encoding/json"
)

func init() {
	web.DOMAIN = 10
}

func Test_Should_allow_visitor_ban_without_prior_visit(t *testing.T) {
	// Given
	var visitorId int64 = 1
	//var visitorIP int32 = 1
	//var visitorRequest []byte = make([]byte, 0)
	var expectedEventCount = 1
	var eventBus chan interface{} = make(chan interface{}, expectedEventCount)
	var eventStream []interface{} = []interface{}{
		//web.NewVisitorRequestReceived(visitorId, visitorIP, visitorRequest),
	}
	var es cqrs.EventStorer = &cqrs.MemoryEventStore {
		Data: eventStream,
	}
	command := web.NewBanVisitor(visitorId)

	// When
	web.Handle(eventBus, es, command)
	
	// Then
	select {
		case event := <-eventBus:
			fmt.Printf("Event: %v\n", event)
			switch e := event.(type) {
				case web.VisitorBanned: {
					if e.Id != visitorId {
						t.Errorf("Invalid visitor id [ %s ]\n", e)
					}
					return
				}
				default: {
					t.Errorf("Incorrect event received [ %s ]\n", e)
				}
			}
			break
		default:
			t.Errorf("Nothing on the bus\n")
	}
}


/*
func HandleJson(data []byte) error {
	cmd := &cqrs.Command{}
	if err := json.Unmarshal(data, cmd); err != nil {
		return err
	}
	switch cmd.Domain {
		case web.DOMAIN: {
			switch cmd.CommandType {
				case web.BanVisitor: {
					_cmd := &BanVisitor{}
					if err := json.Unmarshal(data, _cmd); err != nil {
						return err
					}
					HandleDoSomethingCommand(_cmd)
				}
				case DoAnotherThing: {
					_cmd := &DoAnotherThingCommand{}
					if err := json.Unmarshal(data, _cmd); err != nil {
						return err
					}
					HandleDoAnotherThingCommand(_cmd)
				}
			}
			return nil
		}
	}
	return errors.New("Invalid command received")
}

func On_BanVisitorCommand()
*/

	/*cmd := &TestCommandOneCommand {
		cqrs.Command { Type: 10 },
		"one",
		1,
	}*/
	/*
	cmd1 := NewDoSomethingCommand("one")
	cmd2 := NewDoAnotherThingCommand(2)
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
	*/

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
