# go-hg ☿

Package **hg** provides ☿ **Mercury Protocol** client and server implementations, for the Go programming language.

The **hg** package provides an API in a style similar to the `"net/http"` library that is part of the Go standard library, including support for "middleware".

## Documention

Online documentation, which includes examples, can be found at: http://godoc.org/github.com/reiver/go-hg

[![GoDoc](https://godoc.org/github.com/reiver/go-hg?status.svg)](https://godoc.org/github.com/reiver/go-hg)

## Mercury Protocol

The ☿ **Mercury Protocol** is a simple client-server protocol.

The ☿ **Mercury Protocol** is derived from the _Gemini Protocol_ — basically the _Mercury Protocol_ is the _Gemini Protocol_ without the TLS encryption.
In a sense, the ☿ _Mercury Protocol_ is a “naked” form of the _Gemini Protocol_.

The ☿ **Mercury Protocol**, through the _Gemini Protocol_, was inpired by the _Gopher Protocol_.

## Hypermedia

The ☿ **Mercury Protocol** (and the _Gemini Protocol_) are designed to work with a (specific) **hypermedia** file data format known as **gemtext**.

(The name “gemtext” is short for “gemini text”.)

**Gemtext** is a **formatted text** file data format similar to _markdown_, and inspired by the line typing convention in Gopher.

Some of the built-in handlers in this package will output **gemtext**.

## Mercury Protocol + TLS = Gemini Protocol

If one wants to have a _Gemini Protocol_ server, but handle the TLS encryption at another level from server, then (using this package and) setting up a _Mercury Protocol_ server can enable that.

## Example Mercury Protocol Server

A very very simple ☿ **Mercury Protocol** server is shown in the following code.

This particular ☿ **Mercury Protocol** server just responds to the client with the URI that was in the request, plus the remote address.

```go
package main

import (
	"github.com/reiver/go-hg"
)

func main() {

	var handler hg.Handler = hg.DebugHandler

	err := hg.ListenAndServe(":1961", handler)
	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}
}
```

## Tilde Example Mercury Protocol Server

Another example  ☿ **Mercury Protocol** server is shown in the following code:

```go
package main

import (
	"github.com/reiver/go-hg"
)

func main() {

	var handler hg.Handler = hg.UserDirHandler

	err := hg.ListenAndServe(":1961", handler)
	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}
}
```

Here the handler — `hg.UserDirHandler` — operates similar to Apache's HTTP Server Project's `mod_userdir` — in that it enables user-specific directories such as `/home/username/mercury_public/` to be accessed over the **Mercury Protocol** using the tilde path `mercury://example.com/~username/`

## Example Mercury Protocol Servers With Custom Handler

And finally, here is a custom handler being used in a  ☿ **Mercury Protocol** server:
```go
package main

import (
	"github.com/reiver/go-hg"
)

func main() {

	var handler hg.Handler = myCustomHandler{}

	err := hg.ListenAndServe(":1961", handler)
	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}
}

type myCustomHandler {}

func (receiver myCustomHandler) ServeMercury(w hg.ResponseWriter, r hg.Request) {
	io.WriteString(w, "Hello world!")
}
```

Alternatively, this could be made a bit simple it `hg.HandlerFunc()` is used:
```go
package main

import (
	"github.com/reiver/go-hg"
)

func main() {

	var handler hg.Handler = hg.HandlerFunc(helloworld)

	err := hg.ListenAndServe(":1961", handler)
	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}
}

func helloworld(w hg.ResponseWriter, r *hg.Request) {
	io.WriteString(w, "Hello world!")
}
```

## Mercury Protocol Response Helpers

This package provides a number of helper-functions that make responding to a ☿ Mercury Protocol request easier.
The helper functions are:

| Mercury Protocol Response  | Basic Usage                    | Intermediate Usage                |
| -------------------------- | ------------------------------ | --------------------------------- |
| `10 INPUT`                 | `hg.Input(w, prompt)`          |                                   |
| `11 SENTITIVE INPUT`       | `hg.SensitiveInput(w, prompt)` |                                   |
| `20 SUCCESS`               |                                |                                   |
| `30 REDIRECT - TEMPORARY`  | `hg.RedirectTemporary(w, url)` |                                   |
| `31 REDIRECT - PERMANENT`  | `hg.RedirectPermanent(w, url)` |                                   |
| `40 TEMPORARY FAILURE`     | `hg.TemporaryFailure(w)`       | `hg.TemporaryFailure(w, info)`    |
| `41 SERVER UNAVAILABLE`    | `hg.ServerUnavailable(w)`      | `hg.ServerUnavailable(w, info)`   |
| `42 CGI ERROR`             | `hg.CGIError(w)`               | `hg.CGIError(w, info)`            |
| `43 PROXY ERROR`           | `hg.ProxyError(w)`             | `hg.ProxyError(w, info)`          |
| `44 SLOW DOWN`             | `hg.SlowDown(w, retryAfter)`   |                                   |
| `50 PERMANENT FAILURE`     | `hg.PermanentFailure(w)`       | `hg.PermanentFail ure(w, info)`   |
| `51 NOT FOUND`             | `hg.NotFound(w)`               | `hg.NotFound(w, info)`            |
| `52 GONE`                  | `hg.Gone(w)`                   | `hg.Gone(w, info)`                |
| `53 PROXY REQUEST REFUSED` | `hg.ProxyRequestRefused(w)`    | `hg.ProxyRequestRefused(w, info)` |
| `59 BAD REQUEST`           | `hg.BadRequest(w)`             | `hg.BadRequest(w, info)`          |

## Package Name

The package name of this Go package is **hg** rather than **mercury** because **Hg** is often used as a shorthard for **mercury**.

Nowadays the word **mercury** is used to refer to multiple things —
a Roman god named “Mercury”,
a chemical element named “mercury”,
a planet named “mercury”,
a space-mission named “Project Mercury”, and
now also a network protocol named the “Mercury Protocol”.

The relationship between these different things named “mercury” is as follows —

The **Mercury Protocol** was named after the **Project Mercury** space-mission.

The **Project Mercury** space-mission was named after the **Roman god** named **Mercury**.
The **Project Mercury** space-mission also used a modified version **astrological-symbol** for the **planet mercury** (☿) for its logo.

The chemical-element _mercury_ was also named after **Roman god** named **Mercury**.

An older name for the chemicalelement _mercury_ is **hydrargyrum**.

“Hydrargyrum” is a romanized version of the ancient Greek word “ὑδράργυρος” (hydrargyros).
The ancient Greek word “ὑδράργυρος” (hydrargyros) is a compound word: “ὑδρ” + “άργυρος”.
The first part **ὑδρ-** (hydr-) comes from the root ὕδωρ **water** (although in this context it might be more accurate to interpret it as **liquid** rather than _water_).
The second part **ἄργυρος** (argyros) means **silver** (although in this context it might be more accurate to interpret it as **shiny** rather than _silver_).
So ὑδράργυρος” (hydrargyros) is **water-silver**, although perhaps more accurately interpretted as **liquid-shiny**

“Hg” is the chemical-symbol for the chemical-element _mercury_ because “Hg” is short for “hydrargyrum”.

And thus this, a package that implements the **Mercury Protocol**, is named  `hg`.

## See Also
* [The Mercury protocol (gemini)](gemini://gemini.circumlunar.space/users/solderpunk/gemlog/the-mercury-protocol.gmi)
* [The Mercury protocol (http proxy)](https://portal.mozz.us/gemini/gemini.circumlunar.space/users/solderpunk/gemlog/the-mercury-protocol.gmi)
* [Mailing List thread: “Mercury”](https://lists.orbitalfox.eu/archives/gemini/2020/thread.html#1842)
