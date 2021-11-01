package hg

// HandlerFunc is an adapter that allows one to turn a function into a Handler if the function has the same signature as ServeMercury.
//
// For example:
//
//	func doIt(w hg.ResponseWriter, r hg.Request) {
//		// ...
//	}
//
//	// ...
//
//	var handler hg.Handler = hg.HandlerFunc(doIt)
type HandlerFunc func(ResponseWriter, Request)

func (fn HandlerFunc) ServeMercury(w ResponseWriter, r Request) {
	fn(w,r)
}
