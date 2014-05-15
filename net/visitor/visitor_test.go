package visitor_test

import(
	"github.com/vizidrix/gocqrs/cqrs"
	"github.com/vizidrix/gocqrs/web/visitor"
	"testing"
)

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
				if (e.IPV4Address != 302) {
					t.Errorf("Incorrect ip address")
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