package hg

import (
	"context"
)

// HandlerFunc is an adapter that allows one to turn a function into a Handler if the function has the same signature as ServeMercury.
//
// For example:
//
//	func doIt(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		// ...
//	}
//
//	// ...
//
//	var handler hg.Handler = hg.HandlerFunc(doIt)
type HandlerFunc func(context.Context, ResponseWriter, Request)

func (fn HandlerFunc) ServeMercury(ctx context.Context, w ResponseWriter, r Request) {
	fn(ctx, w, r)
}
