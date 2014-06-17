package eventstore

import (
	//"fmt"
	//"errors"
	"os"
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

func (es *boltES) Put(event cqrs.Event) error {
	return nil
}

func (es *boltES) GetStreamById(domain uint32, id uint64) (<-chan cqrs.Event, error) {
	stream := es.eventChanFactory() // Build an event chan to handle distribution
	//make(chan cqrs.Event, 0) // Don't need buffered chan for no events
	//stream<-nil
	return stream, nil
}

func (es *boltES) GetStreamByDomain(domain uint32) (<-chan cqrs.Event, error) {
	stream := es.eventChanFactory() // Build an event chan to handle distribution
	return stream, nil
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