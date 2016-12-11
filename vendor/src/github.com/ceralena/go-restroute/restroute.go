// Package restroute implements a simple RESTful HTTP routing layer
package restroute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

/*
// TRACE and CONNECT unsupported
var http_methods = []string{
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"DELETE",
}
*/

// Request holds the state for a request, including any named matches parsed
// from the route regular expression.. It's passed to Handlers when their route
// matches.
type Request struct {
	W      http.ResponseWriter
	R      *http.Request
	Params map[string]string // Named matches from the URL
}

// Map is a map of route regular expressions to a MethodMap for each.
type Map map[string]MethodMap

// Compile will prepare this map to be used as a HTTP handler.
//
// Fails if any of the MethodMap regular expressions are invalid.
func (m Map) Compile() (http.Handler, error) {
	return newRouter(m)
}

// MustCompile is the same as Compile(), but panics with any error.
func (m Map) MustCompile() http.Handler {
	h, err := m.Compile()
	if err != nil {
		panic(err)
	}
	return h
}

// Merge merges a set of Maps into a single Map.
//
// If more than one map specifies the same route, the last one in the list
// wins.
func Merge(maps ...Map) Map {
	merged := make(Map)

	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}

	return merged
}

// MethodMap is a map of HTTP request methods to handlers.
type MethodMap map[string]Handler

// Handler is a request handler function that takes a Request. It is called
// when the route and method match.
type Handler func(req Request)

type router struct {
	routes []compiledRoute
}

// Given a method map, create a new router by compiling the regexp for each
// path.
//
// router implements http.Handler
func newRouter(m Map) (*router, error) {
	routes := make([]compiledRoute, 0, len(m))

	for path, handlers := range m {
		route, err := newCompiledRoute(path, handlers)
		if err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}

	return &router{routes}, nil
}

func newCompiledRoute(path string, handlers MethodMap) (compiledRoute, error) {
	r, err := regexp.Compile(path)
	if err != nil {
		return compiledRoute{}, fmt.Errorf("bad path: could not compile :%q: %s", path, err)
	}

	return compiledRoute{r, handlers}, nil

}

func (rtr *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Look for a matching route
	route, params, ok := rtr.getRouteFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		// TODO(Cera) - Clean this up; handle error
		enc := json.NewEncoder(w)
		enc.Encode(map[string]string{"error": "404 not found", "message": fmt.Sprintf("No handler for '%s %s'", r.Method, r.URL.Path)})
		return
	}

	// Look for a matching handler for this HTTP method
	handler, ok := route.handlers[r.Method]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		// TODO(Cera) - Clean this up; handle error
		enc := json.NewEncoder(w)
		enc.Encode(map[string]string{"error": "405 method not allowed", "message": fmt.Sprintf("No handler for '%s %s'", r.Method, r.URL.Path)})
		return
	}

	// Call the handler
	rtr.callHandler(w, r, params, handler)
}

func (rtr *router) getRouteFromRequest(r *http.Request) (*compiledRoute, map[string]string, bool) {
	for _, route := range rtr.routes {
		params, ok := route.match(r.URL)
		if ok {
			return &route, params, true
		}
	}
	return nil, nil, false
}

func (rtr *router) callHandler(w http.ResponseWriter, r *http.Request, params map[string]string, handler Handler) {
	// Build up the request state
	req := Request{w, r, params}

	// Call the handler
	handler(req)
}

type compiledRoute struct {
	r        *regexp.Regexp
	handlers MethodMap
}

// Try to match this route on this URL.
func (cr *compiledRoute) match(u *url.URL) (map[string]string, bool) {
	names := cr.r.SubexpNames()

	sms := cr.r.FindStringSubmatch(u.Path)
	if len(sms) == 0 {
		return nil, false
	}

	params := make(map[string]string, len(names))

	for i, sm := range sms {
		if names[i] != "" {
			params[names[i]] = sm
		}
	}

	return params, true
}
