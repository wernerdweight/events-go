package events

import (
	"log"
	"sort"
	"sync"
)

type EventKey string
type EventPayload interface{}
type subscriberPool map[EventKey]map[int][]EventSubscriber[Event[EventPayload]]

type Event[T EventPayload] interface {
	GetKey() EventKey
	GetPayload() T
}

type EventHub struct {
	subscribers  subscriberPool
	priorities   []int
	syncChannel  chan Event[EventPayload]
	asyncChannel chan Event[EventPayload]
	lock         sync.Mutex
}

type EventSubscriber[T Event[EventPayload]] interface {
	Handle(event T)
	GetKey() EventKey
	GetPriority() int
}

func (h *EventHub) Subscribe(subscriber EventSubscriber[Event[EventPayload]]) {
	h.lock.Lock()
	if _, ok := h.subscribers[subscriber.GetKey()]; !ok {
		h.subscribers[subscriber.GetKey()] = make(map[int][]EventSubscriber[Event[EventPayload]])
	}
	if _, ok := h.subscribers[subscriber.GetKey()][subscriber.GetPriority()]; !ok {
		h.subscribers[subscriber.GetKey()][subscriber.GetPriority()] = make([]EventSubscriber[Event[EventPayload]], 0)
	}
	h.subscribers[subscriber.GetKey()][subscriber.GetPriority()] = append(h.subscribers[subscriber.GetKey()][subscriber.GetPriority()], subscriber)
	h.priorities = append(h.priorities, subscriber.GetPriority())
	sort.Ints(h.priorities)
	h.lock.Unlock()
}

func (h *EventHub) DispatchSync(event Event[EventPayload]) {
	h.syncChannel <- event
}

func (h *EventHub) DispatchAsync(event Event[EventPayload]) {
	h.asyncChannel <- event
}

func (h *EventHub) run(channel chan Event[EventPayload], async bool) {
	for {
		event := <-channel
		if event == nil {
			log.Println("event hub closed")
			return
		}
		if _, ok := h.subscribers[event.GetKey()]; ok {
			subscribers := h.subscribers[event.GetKey()]
			for _, priority := range h.priorities {
				for _, subscriber := range subscribers[priority] {
					if async {
						go subscriber.Handle(event)
						continue
					}
					subscriber.Handle(event)
				}
			}
		}
	}
}

func (h *EventHub) Close() {
	close(h.syncChannel)
	close(h.asyncChannel)
}

var hub *EventHub

func GetEventHub() *EventHub {
	if hub == nil {
		hub = &EventHub{
			subscribers:  make(subscriberPool),
			syncChannel:  make(chan Event[EventPayload]),
			asyncChannel: make(chan Event[EventPayload]),
		}
		go hub.run(hub.syncChannel, false)
		go hub.run(hub.asyncChannel, true)
	}
	return hub
}

func ResetEventHub() {
	hub = nil
}
