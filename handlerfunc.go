package hg

// HandlerFunc is an adapter that allows one to turn a function into a Handler if the function has the same signature as ServeMercury.
type HandlerFunc func(ResponseWriter, Request)

func (fn HandlerFunc) ServeMercury(w ResponseWriter, r Request) {
	fn(w,r)
}
