package events

type EventKey string
type EventPayload interface{}
type subscriberPool map[EventKey]map[int][]EventSubscriber[Event[EventPayload]]

type Event[T EventPayload] interface {
	GetKey() EventKey
	GetPayload() T
}

type EventHub struct {
	subscribers subscriberPool
	channel     chan Event[EventPayload]
}

type EventSubscriber[T Event[EventPayload]] interface {
	Handle(event T)
	GetKey() EventKey
	GetPriority() int
}

func (h *EventHub) Subscribe(subscriber EventSubscriber[Event[EventPayload]]) {
	if _, ok := h.subscribers[subscriber.GetKey()]; !ok {
		h.subscribers[subscriber.GetKey()] = make(map[int][]EventSubscriber[Event[EventPayload]])
	}
	if _, ok := h.subscribers[subscriber.GetKey()][subscriber.GetPriority()]; !ok {
		h.subscribers[subscriber.GetKey()][subscriber.GetPriority()] = make([]EventSubscriber[Event[EventPayload]], 0)
	}
	h.subscribers[subscriber.GetKey()][subscriber.GetPriority()] = append(h.subscribers[subscriber.GetKey()][subscriber.GetPriority()], subscriber)
}

func (h *EventHub) Dispatch(event Event[EventPayload]) {
	h.channel <- event
}

func (h *EventHub) run() {
	for {
		event := <-h.channel
		if _, ok := h.subscribers[event.GetKey()]; ok {
			for _, priority := range h.subscribers[event.GetKey()] {
				for _, subscriber := range priority {
					go subscriber.Handle(event)
				}
			}
		}
	}
}

func (h *EventHub) Close() {
	close(h.channel)
}

var hub *EventHub

func GetEventHub() *EventHub {
	if hub == nil {
		hub = &EventHub{
			subscribers: make(subscriberPool),
			channel:     make(chan Event[EventPayload]),
		}
		go hub.run()
	}
	return hub
}
