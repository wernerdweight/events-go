Events for Go
====================================

A simple event system (dispatcher-subscriber) for Go.

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
)

// Create event
type MyEvent struct{}

func (e *MyEvent) GetKey() events.EventKey {
    return "my-event"
}

func (e *MyEvent) GetPayload() interface{} {
    return "payload"
}

// Create subscriber
type MySubscriber struct{}

func (s *MySubscriber) Handle(event events.Event[events.EventPayload]) {
    fmt.Println(event.GetPayload())
}

func (s *MySubscriber) GetKey() events.EventKey {
    return "my-event"
}

func (s *MySubscriber) GetPriority() int {
    return 0
}

func main() {
    // Create hub
    hub := events.GetEventHub()
	defer hub.Close()

    subscriber := &MySubscriber{}

    // Subscribe to event
    hub.Subscribe(subscriber)

    // Dispatch event
    hub.Dispatch(&MyEvent{})
}
```

License
-------
This package is under the MIT license. See the complete license in the root directory of the bundle.
