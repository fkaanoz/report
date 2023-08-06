package logger

import (
	"fmt"
	"os"
)

type EventType int

const (
	EventDelete EventType = iota + 1
	EventPut
)

func (e EventType) String() string {
	switch e {
	case EventDelete:
		return "DELETE"
	case EventPut:
		return "PUT"
	default:
		return ""
	}
}

type Event struct {
	eventType EventType
	key       string
	value     string
}

type TransactionLogger struct {
	events chan Event
	errors chan error
	file   *os.File
}

func NewTransactionLogger(path string) (*TransactionLogger, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return nil, err
	}

	return &TransactionLogger{
		file: file,
	}, nil
}

func (t *TransactionLogger) WritePut(key, value string) {
	t.events <- Event{
		eventType: EventPut,
		key:       key,
		value:     value,
	}
}

func (t *TransactionLogger) WriteDelete(key string) {
	t.events <- Event{
		eventType: EventDelete,
		key:       key,
	}
}

func (t *TransactionLogger) Run() {
	eventCh := make(chan Event, 16)
	errorCh := make(chan error, 1)

	t.events = eventCh
	t.errors = errorCh

	go func() {
		for ev := range eventCh {
			_, err := fmt.Fprintf(t.file, "%v\t%v\t%v\n", ev.eventType, ev.key, ev.value)
			if err != nil {
				errorCh <- err
			}
		}
	}()
}
