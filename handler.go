package hg

// Handler represents someething that responds to an Mercury Protocol request.
type Handler interface {
	ServeMercury(w ResponseWriter, r Request)
}
