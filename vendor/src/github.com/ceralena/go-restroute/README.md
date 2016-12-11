# go-restroute

[![Build Status](https://travis-ci.org/ceralena/go-restroute.svg?branch=master)](https://travis-ci.org/ceralena/go-restroute) [![codecov](https://codecov.io/gh/ceralena/go-restroute/branch/master/graph/badge.svg)](https://codecov.io/gh/ceralena/go-restroute) [![Go Report Card](https://goreportcard.com/badge/github.com/ceralena/go-restroute)](https://goreportcard.com/report/github.com/ceralena/go-restroute) [![GoDoc](https://godoc.org/github.com/stretchr/testify?status.svg)](https://godoc.org/github.com/stretchr/testify)

go-restroute provides a [Go](https://golang.org) package - `restroute` - for
defining the routes of a RESTful HTTP service.

restroute has no opinions about the controller or anything below that, beyond
specifying that a route handler must have a particular function signature.

Routes are specified as regular expressions, where match groups are parameters.

See godoc for documentation: https://godoc.org/github.com/ceralena/go-restroute

## Design Rationale

* just a handful of types to understand
* prefer strongly typed function composition for middleware over dependency injection with `reflect`
* do not assume or prevent the use of any other library - directly compatible with `net/http`
* no ad-hoc syntax for routes and params - just regular expressions
* no special logic for subroutes or precedence: the caller is expected to
	provide unique routes
* low performance overhead - do as little as possible to find and call a handler

## Usage

Import `github.com/ceralena/go-restroute` in your package.

The types exposed by the package are:

* `restroute.Request` - contains the state for a request
* `restroute.Map`: a map of paths to MethodMap values
* `restroute.MethodMap`: a map of methods to route handlers
* `restroute.Handler` - the expected function signature of a
	route handler: `func(restroute.Request)`

The Request value passed into the handler has the state you'd get in an ordinary
Go HTTP handler, as well as the params corresponding to named groups in the
route regular expression:

	type Request struct {
			W      http.ResponseWriter
			R      *http.Request
			Params map[string]string // Named matches from the URL
	}

## Examples

See [godoc](https://godoc.org/github.com/ceralena/go-restroute) for a full list of examples.


## Implementing Middleware; Stateful Handlers

The obvious downside to the type signature of handlers being a function is
that they don't hold state. This is problematic when (for example) the HTTP
layer of an application should use a shared database connection pool.

Instead of using global state, you can get around this in one of the two following ways:

* function composition & middleware
* using methods or closures as your handlers

Both approaches - or a combination of the two - are appropriate in different circumstances.

## License

MIT license. See the LICENSE file.
