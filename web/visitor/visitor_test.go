package visitor_test

import(
	"github.com/vizidrix/gocqrs/cqrs"
	"github.com/vizidrix/gocqrs/web/visitor"
	"testing"
)

func init() {
	//visitor.DOMAIN = 10
}

func Test_Should_not_allow_Blacklist_without_register(t *testing.T) {
	// Given
	var visitorId uint64 = 1024
	//var visitorIP int32 = 1
	//var visitorRequest []byte = make([]byte, 0)
	var eventBus chan interface{} = make(chan interface{}, 2)
	var eventStream []interface{} = []interface{}{
		//web.NewVisitorRequestReceived(visitorId, visitorIP, visitorRequest),
	}
	var es cqrs.EventStorer = &cqrs.MemoryEventStore {
		Snapshot: nil,
		Data: eventStream,
	}
	command := visitor.NewBlacklist(visitorId, 0)

	// When
	visitor.Handle(eventBus, es, command)
	
	// Then
	select {
		case event := <-eventBus:
			switch e := event.(type) {
				case visitor.Blacklisted: {
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