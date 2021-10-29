package hg

// Caller represents someething that responds to an Mercury Protocol response.
type Caller interface {
	CallMercury(ResponseReader, Request)
}
