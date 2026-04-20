package hg

import (
	"context"
	"time"
)

// Handler represents something that responds to a Mercury Protocol request.
//
// Typically, someone who wants to create a custom Mercury Protocol server would create a type that fits this Handler interface,
// and then (directly or indirectly) pass it to the Serve or ListenAndServe functions.
type Handler interface {
	ServeMercury(ctx context.Context, w ResponseWriter, r Request)
}

// TimeoutHandler is an optional interface that a Handler may implement to declare
// how long it should be allowed to run.
//
// If Timeout returns a positive duration, the server will apply that as a deadline
// on the context passed to ServeMercury.
//
// If Timeout returns zero or a negative duration, no timeout is applied (the handler
// is considered long-lived).
//
// Handlers that do not implement this interface fall back to Server.HandlerTimeout.
type TimeoutHandler interface {
	Handler
	Timeout() time.Duration
}
