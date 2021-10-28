package hg

import (
	"fmt"
	"io"
	"strings"
)

const DebugHandler internalDebugHandler = internalDebugHandler(0)

type internalDebugHandler int

var _ Handler = internalDebugHandler(0)

func (internalDebugHandler) ServeMercury(w ResponseWriter, r Request) {
	var storage strings.Builder

	storage.WriteString("```")
	storage.WriteString(
`
  ███╗   ███╗███████╗██████╗  ██████╗██╗   ██╗██████╗ ██╗   ██╗     
  ████╗ ████║██╔════╝██╔══██╗██╔════╝██║   ██║██╔══██╗╚██╗ ██╔╝     
  ██╔████╔██║█████╗  ██████╔╝██║     ██║   ██║██████╔╝ ╚████╔╝      
  ██║╚██╔╝██║██╔══╝  ██╔══██╗██║     ██║   ██║██╔══██╗  ╚██╔╝       
  ██║ ╚═╝ ██║███████╗██║  ██║╚██████╗╚██████╔╝██║  ██║   ██║        
  ╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═╝   ╚═╝        
                                                                    
██████╗ ██████╗  ██████╗ ████████╗ ██████╗  ██████╗ ██████╗ ██╗     
██╔══██╗██╔══██╗██╔═══██╗╚══██╔══╝██╔═══██╗██╔════╝██╔═══██╗██║     
██████╔╝██████╔╝██║   ██║   ██║   ██║   ██║██║     ██║   ██║██║     
██╔═══╝ ██╔══██╗██║   ██║   ██║   ██║   ██║██║     ██║   ██║██║     
██║     ██║  ██║╚██████╔╝   ██║   ╚██████╔╝╚██████╗╚██████╔╝███████╗
╚═╝     ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚═════╝  ╚═════╝ ╚═════╝ ╚══════╝
`,
	)
	storage.WriteString("```\r\n")
	storage.WriteString("\r\n")
	storage.WriteString("# Mercury Request\r\n")
	storage.WriteString("\r\n")
	storage.WriteString("The following was the Mercury Protocol request received:\r\n")
	storage.WriteString("```")
	fmt.Fprintf(&storage, "%q", r)
	storage.WriteString("```\r\n")
	storage.WriteString("☿\r\n")

	w.WriteHeader(StatusSuccess, "text/gemini")
	io.WriteString(w, storage.String())
}
