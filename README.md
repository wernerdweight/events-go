Events for Go
====================================

A simple event hub for Go.

[![Build Status](https://www.travis-ci.com/wernerdweight/events-go.svg?branch=master)](https://www.travis-ci.com/wernerdweight/events-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/wernerdweight/events-go)](https://goreportcard.com/report/github.com/wernerdweight/events-go)
[![GoDoc](https://godoc.org/github.com/wernerdweight/events-go?status.svg)](https://godoc.org/github.com/wernerdweight/events-go)
[![go.dev](https://img.shields.io/badge/go.dev-pkg-007d9c.svg?style=flat)](https://pkg.go.dev/github.com/wernerdweight/events-go)


Installation
------------

### 1. Installation

```bash
go get github.com/wernerdweight/events-go
```

Usage
------------

```go
package main

import (
	"fmt"
	events "github.com/wernerdweight/events-go"
	"log"
	"time"
)

type TestEvent struct {
	payload string
	delay   time.Duration
}

func (e *TestEvent) GetKey() events.EventKey {
	return "test"
}

func (e *TestEvent) GetPayload() events.EventPayload {
	return e.payload
}

type TestSubscriber struct {
	isHandled bool
	buffer    string
	priority  int
}

func (s *TestSubscriber) Handle(event events.Event[events.EventPayload]) {
	time.Sleep(event.(*TestEvent).delay)
	s.isHandled = true
	s.buffer += event.GetPayload().(string)
}

func (s *TestSubscriber) GetKey() events.EventKey {
	return "test"
}

func (s *TestSubscriber) GetPriority() int {
	return s.priority
}

func asyncExample() {
	hub := events.GetEventHub()
	subscriber := &TestSubscriber{}
	hub.Subscribe(subscriber)

	// async dispatch (events dispatched in goroutine)
	hub.DispatchAsync(&TestEvent{
		payload: "A",
		delay:   100 * time.Millisecond,
	})
	hub.DispatchAsync(&TestEvent{
		payload: "B",
		delay:   0,
	})
	time.Sleep(120 * time.Millisecond)

	log.Print(subscriber.isHandled) // true
	log.Print(subscriber.buffer)    // "BA" (because of async dispatch)
}

func syncExample() {
	// sync dispatch (events dispatched in main goroutine)
	hub.DispatchSync(&TestEvent{
		payload: "A",
		delay:   100 * time.Millisecond,
	})
	hub.DispatchSync(&TestEvent{
		payload: "B",
		delay:   0,
	})
	time.Sleep(120 * time.Millisecond)

	log.Print(subscriber.isHandled) // true
	log.Print(subscriber.buffer)    // "AB" (because of sync dispatch)
}

// set priority to set order of handling for subscribers subscribed to the same event

type TestPriorityEvent struct{}

func (e *TestPriorityEvent) GetKey() events.EventKey {
	return "test-priority"
}

func (e *TestPriorityEvent) GetPayload() events.EventPayload {
	return ""
}

type TestPrioritySubscriber struct {
	buffer   *string
	priority int
}

func (s *TestPrioritySubscriber) Handle(event events.Event[events.EventPayload]) {
	*s.buffer += fmt.Sprintf("%d", s.priority)
}

func (s *TestPrioritySubscriber) GetKey() events.EventKey {
	return "test-priority"
}

func (s *TestPrioritySubscriber) GetPriority() int {
	return s.priority
}

func priorityExample() {
	hub := events.GetEventHub()
	buffer := ""
	subscriber1 := &TestPrioritySubscriber{buffer: &buffer, priority: 1}
	subscriber2 := &TestPrioritySubscriber{buffer: &buffer, priority: 2}
	subscriber_1 := &TestPrioritySubscriber{buffer: &buffer, priority: -1}
	subscriber0 := &TestPrioritySubscriber{buffer: &buffer}
	hub.Subscribe(subscriber1)
	hub.Subscribe(subscriber2)
	hub.Subscribe(subscriber_1)
	hub.Subscribe(subscriber0)
	hub.DispatchSync(&TestPriorityEvent{})
	hub.DispatchSync(&TestPriorityEvent{})
	time.Sleep(20 * time.Millisecond)
	
	log.Print(buffer) // "-1012-1012" (according to priority)
}
```

License
-------
This package is under the MIT license. See the complete license in the root directory of the bundle.
