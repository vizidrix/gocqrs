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
	var eventBus chan cqrs.Event = make(chan cqrs.Event, 1)

	// Given
	es := &cqrs.MemoryEventStore { Data: []cqrs.Event{} }

	// When
	visitor.Handle(eventBus, es, visitor.NewBlacklist(1024, 0))
	//var visitorId uint64 = 1024
	//var visitorIP int32 = 1
	//var visitorRequest []byte = make([]byte, 0)
	
	//var eventStream []cqrs.Event = []cqrs.Event{
		//web.NewVisitorRequestReceived(visitorId, visitorIP, visitorRequest),
	//}
	//var es cqrs.EventStorer = &cqrs.MemoryEventStore {
	//	Snapshot: nil,
	//	Data: eventStream,
	//}
	//command := visitor.NewBlacklist(visitorId, 0)

	// When
	//visitor.Handle(eventBus, es, command)
	
	// Then
	select {
		case event := <-eventBus:
			switch e := event.(type) {
				case visitor.Blacklisted: {
					if e.Id != 1023 {
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

func Test_Should_allow_valid_visitor_register(t *testing.T) {
	var eventBus chan cqrs.Event = make(chan cqrs.Event, 1)

	// Given
	es := &cqrs.MemoryEventStore { Data: []cqrs.Event{} }

	// When
	visitor.Handle(eventBus, es, visitor.NewRegisterIPV4(10, 20, 30))
	
	// Then
	select {
	case event := <-eventBus:
		switch e := event.(type) {
			case visitor.Registered: {
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