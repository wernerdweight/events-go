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
	subscribers subscriberPool
	priorities  []int
	channel     chan Event[EventPayload]
	lock        sync.Mutex
}

type EventSubscriber[T Event[EventPayload]] interface {
	Handle(event T) error
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

func (h *EventHub) handleEvent(event Event[EventPayload], async bool) error {
	if _, ok := h.subscribers[event.GetKey()]; ok {
		subscribers := h.subscribers[event.GetKey()]
		for _, priority := range h.priorities {
			for _, subscriber := range subscribers[priority] {
				if async {
					go func() {
						err := subscriber.Handle(event)
						if err != nil {
							log.Println(err)
						}
					}()
					continue
				}
				err := subscriber.Handle(event)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (h *EventHub) DispatchSync(event Event[EventPayload]) error {
	return h.handleEvent(event, false)
}

func (h *EventHub) DispatchAsync(event Event[EventPayload]) {
	h.channel <- event
}

func (h *EventHub) run(channel chan Event[EventPayload]) {
	for {
		event := <-channel
		if event == nil {
			log.Println("event hub closed")
			return
		}
		err := h.handleEvent(event, true)
		if err != nil {
			log.Println(err)
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
		go hub.run(hub.channel)
	}
	return hub
}

func ResetEventHub() {
	hub = nil
}
