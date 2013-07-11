package gocqrs

import ()

/*
type MemoryEventStore struct {
	// A map by Kind which further maps by Aggregate Id to a slice of events
	aggregateKindMap map[string]*AggregateStore
}

func NewMemoryEventStore() *MemoryEventStore {
	return &MemoryEventStore{
		aggregateKindMap: make(map[string]*AggregateStore),
	}
}

func (eventStore *MemoryEventStore) Of(kind string) AggregateIdMapper {
	// Do a quick check to ensure that the inner map exsits
	if _, ok := eventStore.aggregateKindMap[kind]; !ok {
		eventStore.aggregateKindMap[kind] = &AggregateStore{
			aggregateKeyMap: make(map[int64]*EventSet),
		}
	}
	return eventStore.aggregateKindMap[kind]

	//return func(event *Event) error {
}
*/

/*
	return func(id int64) EventWriter {
		if _, ok :=
		return func(event *Event) error {
			if aggregateMap, ok := eventStore[kind]; ok {

			} else {
				return errors.New("text")
			}
		}
		return nil
	}
*/
