package gocqrs_example

import (
	"errors"
	"fmt"
	cqrs "github.com/vizidrix/gocqrs"
	"log"
)

func ignore() { log.Println("") }

type Person struct {
	cqrs.Aggregate
	FirstName string `datastore:",noindex"`
	LastName  string `datastore:",noindex"`
	SSN       int64  `datastore:",noindex"`
	Profile   string `datastore:",noindex"`
}

type RegisterPerson struct {
	cqrs.Command
	FirstName string `datastore:",noindex"`
	LastName  string `datastore:",noindex"`
	SSN       int64  `datastore:",noindex"`
	Profile   string `datastore:",noindex"`
}

type UpdateProfile struct {
	cqrs.Command
	Profile string `datastore:",noindex"`
}

type PersonRegistered struct {
	cqrs.Event
	FirstName string `datastore:",noindex"`
	LastName  string `datastore:",noindex"`
	SSN       int64  `datastore:",noindex"`
	Profile   string `datastore:",noindex"`
}

type ProfileUpdated struct {
	cqrs.Event
	Profile string `datastore:",noindex"`
}

type PersonSummary struct { // View
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
}

func NewPerson(id int64, firstName string, lastName string, ssn int64, profile string) *Person {
	return &Person{
		Aggregate: cqrs.NewAggregate("person", id),
		FirstName: firstName,
		LastName:  lastName,
		SSN:       ssn,
		Profile:   profile,
	}
}

func NewRegisterPerson(id int64, version int64, firstName string, lastName string, ssn int64, profile string) *RegisterPerson {
	return &RegisterPerson{
		Command:   cqrs.NewCommand(id, version),
		FirstName: firstName,
		LastName:  lastName,
		SSN:       ssn,
		Profile:   profile,
	}
}

func NewUpdateProfile(id int64, version int64, profile string) *UpdateProfile {
	return &UpdateProfile{
		Command: cqrs.NewCommand(id, version),
		Profile: profile,
	}
}

func NewPersonRegistered(id int64, version int64, firstName string, lastName string, ssn int64, profile string) *PersonRegistered {
	return &PersonRegistered{
		Event:     cqrs.NewEvent(id, version),
		FirstName: firstName,
		LastName:  lastName,
		SSN:       ssn,
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
