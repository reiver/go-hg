package hg

// A Mercury Protocol server runs over TCP.
// TCP has communications happening over TCP-ports.
// A client-server protocol (including the Mercury Protocol) typically defines a default-TCP-port for servers.
// For the Mercury Protocol, this default-TCP-port is: 1961.
//
// This constant — ‘DefaultTCPPort’ — can be used when once wants to use the default-TCP--port for a Mercury Protocol server.
//
// For example:
//
//	var domain string = "example.com"
//	
//	var address string = fmt.Sprint("%s:%d", domain, hg.DefaultTCPPort)
//	
//	err := hg.hg.ListenAndServe(address, handler)
const DefaultTCPPort = 1961
