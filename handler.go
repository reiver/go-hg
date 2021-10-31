package hg

// Handler represents something that responds to an Mercury Protocol request.
//
// Typically, someone who wants to create a custom Mercury Protocol server would create a type that fits this Handler interface,
// and then (directly or indirectly) pass it to the Serve to ListenAndServe functions.
type Handler interface {
	ServeMercury(w ResponseWriter, r Request)
}
