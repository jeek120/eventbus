package eventbus

// EventHandlerFunc is a function that can be used as a event handler.
type EventHandlerFunc func(Event) error

type EventHandler interface {
	HandleEvent(ev Event) error
}

// HandleCommand implements the HandleCommand method of the CommandHandler.
func (h EventHandlerFunc) HandleCommand(ev Event) error {
	return h(ev)
}

// EventHandlerMiddleware is a function that middlewares can implement to be
// able to chain.
type EventHandlerMiddleware func(EventHandlerFunc) EventHandlerFunc

// UseEventHandlerMiddleware wraps a CommandHandler in one or more middleware.
func UseEventHandlerMiddleware(h EventHandlerFunc, middleware ...EventHandlerMiddleware) EventHandlerFunc {
	// Apply in reverse order.
	for i := len(middleware) - 1; i >= 0; i-- {
		m := middleware[i]
		h = m(h)
	}
	return h
}
