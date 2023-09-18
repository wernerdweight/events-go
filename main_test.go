package events

import (
	"testing"
	"time"
)

func TestGetEventHub(t *testing.T) {
	testingHub := GetEventHub()
	if testingHub == nil {
		t.Error("hub should not be nil")
	}
}

func TestEventHub_Subscribe(t *testing.T) {
	testingHub := GetEventHub()
	subscriber := &TestSubscriber{}
	testingHub.Subscribe(subscriber)
	if len(testingHub.subscribers) != 1 {
		t.Error("hub subscribers should have 1 subscriber")
	}
	if len(testingHub.subscribers[subscriber.GetKey()]) != 1 {
		t.Error("hub subscribers should have 1 subscriber for the key")
	}
	if len(testingHub.subscribers[subscriber.GetKey()][subscriber.GetPriority()]) != 1 {
		t.Error("hub subscribers should have 1 subscriber for the key and priority")
	}
}

func TestEventHub_Dispatch(t *testing.T) {
	testingHub := GetEventHub()
	subscriber := &TestSubscriber{}
	testingHub.Subscribe(subscriber)
	testingHub.Dispatch(&TestEvent{})
	time.Sleep(100 * time.Millisecond)
	if !subscriber.isHandled {
		t.Error("event should be handled")
	}
}

func TestEventHub_Close(t *testing.T) {
	testingHub := GetEventHub()
	subscriber := &TestSubscriber{}
	testingHub.Subscribe(subscriber)
	testingHub.Close()
	if len(testingHub.channel) != 0 {
		t.Error("channel should be closed")
	}
}

type TestEvent struct {
}

func (e *TestEvent) GetKey() EventKey {
	return "test"
}

func (e *TestEvent) GetPayload() EventPayload {
	return nil
}

type TestSubscriber struct {
	isHandled bool
}

func (s *TestSubscriber) Handle(event Event[EventPayload]) {
	s.isHandled = true
}

func (s *TestSubscriber) GetKey() EventKey {
	return "test"
}

func (s *TestSubscriber) GetPriority() int {
	return 0
}
