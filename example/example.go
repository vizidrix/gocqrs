package gocqrs_example

import (
	"errors"
	"fmt"
	cqrs "github.com/vizidrix/gocqrs"
	"log"
)

type Person struct {
	cqrs.Aggregate
	FirstName string `datastore:",noindex"`
	LastName  string `datastore:",noindex"`
	SSN       int64  `datastore:",noindex"`
	Profile   string `datastore:",noindex"`
}

type RegisterPerson struct { // Command
	FirstName string `datastore:",noindex"`
	LastName  string `datastore:",noindex"`
	SSN       int64  `datastore:",noindex"`
	Profile   string `datastore:",noindex"`
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

func LoadPersonSync(eventChan <-chan cqrs.IEvent) (person *Person, err error) {
	for e := range eventChan {
		if e == nil {
			err = errors.New("Invalid event on channel")
			return
		}
		switch event := (e).(type) {
		default:
			{
				err = errors.New(fmt.Sprintf("Invalid event type on channel: %s", e))
				return
			}
		case *PersonRegistered:
			{
				log.Printf("Received PersonRegistered: %s", event)
				person = NewPerson(event.GetId(), event.FirstName, event.LastName, event.SSN, event.Profile)
			}
		case *ProfileUpdated:
			{
				log.Printf("Received ProfileUpdated: %s", event)
				person.Profile = event.Profile
			}
		}
		person.IncrementVersion()
	}
	return
}

func LoadPerson(eventChan <-chan *cqrs.IEvent) (<-chan *Person, <-chan error) {
	result := make(chan *Person)
	err := make(chan error)
	go func() {
		var person *Person
		defer func() {
			if person == nil {
				return
			}
			result <- person
		}()
		select {
		case e := <-eventChan:
			{
				if e == nil {
					return
				}
				switch event := (*e).(type) {
				case PersonRegistered:
					{
						person = NewPerson(event.GetId(), event.FirstName, event.LastName, event.SSN, event.Profile)
					}
				}
				person.IncrementVersion()
			}
		}
	}()
	return result, err
}

//func PersonLoader_Handle_PersonRegistered(*Person person, )
