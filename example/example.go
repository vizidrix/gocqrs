package gocqrs_example

import (
	"errors"
	"fmt"
	"github.com/vizidrix/gocqrs/cqrs"
)

var DOMAIN int32 = 1

const ( // Aggregates
	_ = 0 + iota
	Person
)

const ( // Commands
	_ = 0 + iota
	RegisterPerson
	UpdateProfile
)

const ( // Events
	_ = 0 + iota
	PersonRegistered
	ProfileUpdated
)

type Person struct {
	cqrs.Aggregate
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Profile string `json:"profile"`
}

type RegisterPerson struct {
	cqrs.Command
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Profile string `json:"profile"`
}

type UpdateProfile struct {
	cqrs.Command
	Profile string `json:"profile"`
}

type PersonRegistered struct {
	cqrs.Event
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Profile string `json:"profile"`
}

type ProfileUpdated struct {
	cqrs.Event
	Profile string `json:"profile"`
}

type PersonSummary struct { // View
	Id        int64  `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	FullName  string `json:"fullname"`
}

func NewPerson(id int64, firstName string, lastName string, profile string) *Person {
	return &Person{
		Aggregate: cqrs.NewAggregate(DOMAIN, Person, id),
		FirstName: firstName,
		LastName:  lastName,
		Profile:   profile,
	}
}

func NewRegisterPerson(id int64, version int64, firstName string, lastName string, profile string) *RegisterPerson {
	return &RegisterPerson{
		Command:   cqrs.NewCommand(id, version),
		FirstName: firstName,
		LastName:  lastName,
		Profile:   profile,
	}
}

func NewUpdateProfile(id int64, version int64, profile string) *UpdateProfile {
	return &UpdateProfile{
		Command: cqrs.NewCommand(id, version),
		Profile: profile,
	}
}

func NewPersonRegistered(id int64, version int64, firstName string, lastName string, profile string) *PersonRegistered {
	return &PersonRegistered{
		Event:     cqrs.NewEvent(id, version),
		FirstName: firstName,
		LastName:  lastName,
		Profile:   profile,
	}
}

func NewProfileUpdated(id int64, version int64, profile string) *ProfileUpdated {
	return &ProfileUpdated{
		Event:   cqrs.NewEvent(id, version),
		Profile: profile,
	}
}

func LoadPerson(aggregate interface{}, eventChan <-chan cqrs.IEvent) (result interface{}, err error) {
	var person *Person
	// A nill means we're loading from baseline, otherwise we're appending events
	if aggregate != nil {
		person = aggregate.(*Person)
	}
	for e := range eventChan {
		if e == nil {
			err = errors.New("Invalid event on channel")
			return
		}
		// Validate aggregate version
		version := e.(cqrs.IHasVersion).GetVersion()
		if person == nil {
			// New aggregate must have a version of 1
			if version != int64(1) {
				err = errors.New(fmt.Sprintf("Event version should be 1 but was %s", version))
				return
			}
		} else {
			// Existing aggregate must have an incremented version
			if version != person.GetVersion()+1 {
				err = errors.New(fmt.Sprintf("Event version should be %s but was %s", person.GetVersion(), (e).(cqrs.IHasVersion).GetVersion()))
				return
			}
		}
		// For each event in the channel get it's type and apply it to the aggregate
		switch event := (e).(type) {
		default:
			{
				err = errors.New(fmt.Sprintf("Invalid event type on channel: %s", e))
				return
			}
		case *PersonRegistered:
			{
				if person != nil {
					err = errors.New("Cannot register person with conflicting id")
					return
				}
				person = NewPerson(event.GetId(), event.FirstName, event.LastName, event.SSN, event.Profile)
			}
		case *ProfileUpdated:
			{
				person.Profile = event.Profile
			}
		}
		person.IncrementVersion()
	}
	result = person
	return
}
