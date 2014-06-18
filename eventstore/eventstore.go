// Eventstore provides interface and implementation used to store
// and retrieve event streams related to domain aggregates
package eventstore

import (
	"fmt"
	//"errors"
	"os"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/boltdb/bolt"
	"github.com/vizidrix/gocqrs/cqrs"
)

type EventChanFactory func() chan cqrs.Event

type EventStore interface {
	Put(cqrs.Event) error
	GetStreamById(domain uint32, id uint64) (<-chan cqrs.Event, error)
	GetStreamByDomain(domain uint32) (<-chan cqrs.Event, error)
}

type boltES struct {
	DB *bolt.DB
	eventChanFactory EventChanFactory
}

func NewBoltES(
	path string, 
	mode os.FileMode, 
	eventChanFactory EventChanFactory,
	) (es *boltES, err error) {
	var db *bolt.DB
	if db, err = bolt.Open(path, mode); err != nil {
		return // Return with the error
	}
	es = &boltES {
		DB: db,
		eventChanFactory: eventChanFactory,
	}
	return
}

// Appends the provided event to the set by it's domain and aggregate id
func (es *boltES) Put(event cqrs.Event) error {
	err := es.DB.Update(func(tx *bolt.Tx) error {
		domainKey := EncodeDomain(event.GetDomain())
		bucket := tx.Bucket(domainKey)
		if bucket == nil {
			fmt.Printf("Nil bucket\n")
			if _, err := tx.CreateBucketIfNotExists(domainKey); err != nil {
				panic(fmt.Sprintf("Error creating bucket [ %s ]\n", err))
			}
			bucket = tx.Bucket(domainKey)
			fmt.Printf("Made a bucket\n")
		}
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		if err := enc.Encode(event); err != nil {
			return err
		}
		fmt.Printf("Put Event Id: [ %d ] -> [ %v ]\n", event.GetId(), EncodeAggregateId(event.GetId()))
		//var data []byte = make([]byte, 0, 0)
		//EncodeAggregateId(event.GetId())
		bucket.Put(EncodeAggregateId(event.GetId()), buffer.Bytes())//buffer.Bytes())
		return nil
		})
	return err // Need to test retrieve of value and make data a slice of event
}

// Loads all events for the specified aggregate by it's id
func (es *boltES) GetStreamById(domain uint32, id uint64) (<-chan cqrs.Event, error) {
	stream := es.eventChanFactory() // Build an event chan to handle distribution
	//make(chan cqrs.Event, 0) // Don't need buffered chan for no events
	//stream<-nil
	return stream, nil
}

// Loads all events within the specified domain bucket
func (es *boltES) GetStreamByDomain(domain uint32) (<-chan cqrs.Event, error) {
	stream := es.eventChanFactory() // Build an event chan to handle distribution
	return stream, nil
}

func EncodeDomain(domain uint32) []byte {
	temp := make([]byte, 4, 4)
	binary.BigEndian.PutUint32(temp, domain)
	//fmt.Printf("Encoded domain [ %v ]\n", temp)
	return temp
}

func EncodeAggregateId(aggregateId uint64) []byte {
	temp := make([]byte, 8, 8)
	binary.BigEndian.PutUint64(temp, aggregateId)
	//fmt.Printf("Encoded aggregate id [ %v ]\n", temp)
	return temp
}

/*
var (
	BUCKET_EVENTS []byte = []byte("events")
)

// Initialize the database if necessary
func OpenEventStore(path string, mode os.FileMode) (db *bolt.DB, err error) {
	if db, err = bolt.Open(path, mode); err != nil {
		return
	}
	// Setup default buckets
	if err = makeBuckets(db, BUCKET_EVENTS); err != nil {
		db.Close()
		db = nil
	}
	return
}

func makeBuckets(db *bolt.DB, names ...[]byte) error {
	for _, name := range names {
		if err := db.Update(makeBucket(name)); err != nil {
			return err
		}
	}
	return nil
}

func makeBucket(name []byte) func(*bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(name)); err != nil {
			panic(fmt.Sprintf("Error creating bucket [ %s ]\n", err))
		}
		return nil
	}
}
*/