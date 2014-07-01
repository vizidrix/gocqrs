package cqrs

/* View / View Handler Interface Type? */

type View interface {
	HandleStream(eventchan chan Event)
	HandleEvent(event Event)
}
