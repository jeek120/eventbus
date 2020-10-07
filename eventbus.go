package eventbus

import (
	"errors"
)

// EventBus is a command handler that handles commands by routing to the
// registered CommandHandlers.
type EventBus struct {
	handlers map[EventType][]EventHandler
}

// NewEventBus creates a EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]EventHandler),
	}
}

// HandleEvent handles a event with a handler capable of handling it.
func (h *EventBus) HandleEvent(ev Event) error {
	if handlers, ok := h.handlers[ev.EventType()]; ok {
		for _, handler := range handlers {
			return handler.HandleEvent(ev)
		}
	}

	return ErrHandlerNotFound
}

// AddHandler adds a handler for a specific event.
func (h *EventBus) AddHandler(handler EventHandler, evTypes ...EventType) error {
	for _, evType := range evTypes {
		if _, ok := h.handlers[evType]; !ok {
			h.handlers[evType] = make([]EventHandler, 0)
		}
		/*
		   if !reflect.TypeOf(handler).In(1).Implements( reflect.TypeOf((*Event)(nil)).Elem() ) {
		       panic(ErrHandlerNotEvent)
		   }*/

		h.handlers[evType] = append(h.handlers[evType], handler)
	}
	return nil
}

/*
func (h *EventBus) AddHandler(handler EventHandler, evs ...Event) error {
    for _, ev := range evs {
        evType := ev.EventType()
        if _, ok := h.handlers[evType]; !ok {
            h.handlers[evType] = make([]EventHandler, 0)
        }
        if !reflect.TypeOf(handler).In(1).Implements( reflect.TypeOf((*Event)(nil)).Elem() ) {
            panic(ErrHandlerNotEvent)
        }

        h.handlers[evType] = append(h.handlers[evType], handler)
    }
    return nil
}*/

// ErrHandlerNotFound is when no handler can be found.
var ErrHandlerNotFound = errors.New("no handlers for event")

// ErrHandlerNotFound is when no handler can be found.
var ErrHandlerNotEvent = errors.New("the second param is event")

// ErrHandlerAlreadySet is when a handler is already registered for a command.
var ErrHandlerAlreadySet = errors.New("handler is already set")
