package hg

const (
	// A Mercury Protocol server runs over TCP.
	// TCP has communications happening over TCP-ports.
	// A client-server protocol (including the Mercury Protocol) typically defines a default-TCP-port for servers.
	// For the Mercury Protocol, this default-TCP-port is: 1961.
	//
	// This constant — ‘DefaultTCPPort’ — can be used when once wants to use the default-TCP-port for a Mercury Protocol server.
	//
	// For example:
	//
	//	var domain string = "example.com"
	//	
	//	var address string = fmt.Sprintf("%s:%d", domain, hg.DefaultTCPPort)
	//	
	//	err := hg.ListenAndServe(address, handler)
	DefaultTCPPort       =  1961
	DefaultTCPPortString = "1961"

	// A Gemini Protocol server runs over TLS over TCP.
	// TCP has communications happening over TCP-ports.
	// A client-server protocol (including the Gemini Protocol) typically defines a default-TCP-port for servers.
	// For the Gemini Protocol, this default-TCP-port is: 1965.
	//
	// This constant — ‘DefaultTCPPortTLS’ — can be used when once wants to use the default-TCP-port for a Gemini Protocol server.
	//
	// For example:
	//
	//	var domain string = "example.com"
	//	
	//	var address string = fmt.Sprintf("%s:%d", domain, hg.DefaultTCPPortTLS)
	//	
	//	err := hg.ListenAndServe(address, handler)
	DefaultTCPPortTLS       =  1965
	DefaultTCPPortTLSString = "1965"
)
