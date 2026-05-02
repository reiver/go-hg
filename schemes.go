package hg

const (
	Scheme    = "mercury"
	SchemeTLS = "gemini"
)

func hasValidScheme(scheme string) bool {
	return  Scheme == scheme || SchemeTLS == scheme
}
