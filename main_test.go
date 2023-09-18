package events

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TestSuite struct {
	suite.Suite
	testingHub *EventHub
}

func (s *TestSuite) SetupTest() {
	ResetEventHub()
	s.testingHub = GetEventHub()
}

func (s *TestSuite) TearDownTest() {
	s.testingHub.Close()
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestGetEventHub() {
	testingHub := GetEventHub()
	s.NotNil(testingHub, "hub should not be nil")
}

func (s *TestSuite) TestEventHub_Subscribe() {
	testingHub := GetEventHub()
	subscriber := &TestSubscriber{}
	testingHub.Subscribe(subscriber)
	s.Len(testingHub.subscribers, 1, "hub subscribers should have 1 subscriber")
	s.Len(testingHub.subscribers[subscriber.GetKey()], 1, "hub subscribers should have 1 subscriber for the key")
	s.Len(testingHub.subscribers[subscriber.GetKey()][subscriber.GetPriority()], 1, "hub subscribers should have 1 subscriber for the key and priority")
}

func (s *TestSuite) TestEventHub_DispatchAsync() {
	testingHub := GetEventHub()
	subscriber := &TestSubscriber{}
	testingHub.Subscribe(subscriber)
	testingHub.DispatchAsync(&TestEvent{
		payload: "A",
		delay:   100 * time.Millisecond,
	})
	testingHub.DispatchAsync(&TestEvent{
		payload: "B",
		delay:   0,
	})
	time.Sleep(120 * time.Millisecond)
	s.True(subscriber.isHandled, "event should be handled")
	s.Equal("BA", subscriber.buffer, "async event should not be handled in order")
}

func (s *TestSuite) TestEventHub_DispatchSync() {
	testingHub := GetEventHub()
	subscriber := &TestSubscriber{}
	testingHub.Subscribe(subscriber)
	err := testingHub.DispatchSync(&TestEvent{
		payload: "A",
		delay:   100 * time.Millisecond,
	})
	s.Nil(err, "event should be handled")
	err = testingHub.DispatchSync(&TestEvent{
		payload: "B",
		delay:   0,
	})
	s.Nil(err, "event should be handled")
	s.True(subscriber.isHandled, "event should be handled")
	s.Equal("AB", subscriber.buffer, "sync event should be handled in order")
}

func (s *TestSuite) TestEventHub_Priority() {
	testingHub := GetEventHub()
	buffer := ""
	subscriber1 := &TestPrioritySubscriber{buffer: &buffer, priority: 1}
	subscriber2 := &TestPrioritySubscriber{buffer: &buffer, priority: 2}
	subscriber_1 := &TestPrioritySubscriber{buffer: &buffer, priority: -1}
	subscriber0 := &TestPrioritySubscriber{buffer: &buffer}
	testingHub.Subscribe(subscriber1)
	testingHub.Subscribe(subscriber2)
	testingHub.Subscribe(subscriber_1)
	testingHub.Subscribe(subscriber0)
	err := testingHub.DispatchSync(&TestPriorityEvent{})
	s.Nil(err, "event should be handled")
	err = testingHub.DispatchSync(&TestPriorityEvent{})
	s.Nil(err, "event should be handled")
	s.Equal("-1012-1012", buffer, "sync event should be handled by priority")
}

func (s *TestSuite) TestEventHub_Close() {
	testingHub := GetEventHub()
	subscriber := &TestSubscriber{}
	testingHub.Subscribe(subscriber)
	s.Len(testingHub.channel, 0, "channel should be closed")
}

type TestEvent struct {
	payload string
	delay   time.Duration
}

func (e *TestEvent) GetKey() EventKey {
	return "test"
}

func (e *TestEvent) GetPayload() EventPayload {
	return e.payload
}

type TestSubscriber struct {
	isHandled bool
	buffer    string
	priority  int
}

func (s *TestSubscriber) Handle(event Event[EventPayload]) error {
	time.Sleep(event.(*TestEvent).delay)
	s.isHandled = true
	s.buffer += event.GetPayload().(string)
	return nil
}

func (s *TestSubscriber) GetKey() EventKey {
	return "test"
}

func (s *TestSubscriber) GetPriority() int {
	return s.priority
}

type TestPriorityEvent struct{}

func (e *TestPriorityEvent) GetKey() EventKey {
	return "test-priority"
}

func (e *TestPriorityEvent) GetPayload() EventPayload {
	return ""
}

type TestPrioritySubscriber struct {
	buffer   *string
	priority int
}

func (s *TestPrioritySubscriber) Handle(event Event[EventPayload]) error {
	*s.buffer += fmt.Sprintf("%d", s.priority)
	return nil
}

func (s *TestPrioritySubscriber) GetKey() EventKey {
	return "test-priority"
}

func (s *TestPrioritySubscriber) GetPriority() int {
	return s.priority
}
