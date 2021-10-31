package hg

// CallerFunc is an adapter that allows one to turn a function into a Caller if the function has the same signature as CallMercury.
type CallerFunc func(ResponseReader, Request)

func (fn CallerFunc) CallMercury(rr ResponseReader, r Request) {
	fn(rr,r)
}
