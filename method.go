//Package method provides helpers for dispatching http requests
//based on the method of the request.
//
//All of the code is trivial.
//The meat of the package is a little more than a dozen SLOC.
//It is just as easy to inline this into your implementation
//or to write custom helpers more suited to your needs.
//
//Motivation
//
//This functionality is often built into http routers,
//but this is a troublesome conflation of ideas.
//
//An http router should recognize a route.
//
//What happens after that is up to the handler for that route.
//
//Mixing route dispatching with method dispatching complicates the router.
//The router needs to store multiple handlers per route for the different
//methods it handles.
//
//While most routers let you set a 404 Not Found handler,
//few let you set a 405 Invalid Method handler, and fewer still per route.
//Thus, if you need to handle 405's specially, you end up having to
//write a handler that matches all methods that contains your custom
//405 logic and your own method dispatching logic.
//
//By separating route dispatch from method dispatch, you end up with two
//simple, composable pieces that provide more flexibility.
//
//The router need only dispatch based on the url
//and the method dispatcher need only dispatch based on the request method.
//
//Of course, there is certainly a convenience to be gained
//by specifying something like
//	router.Get("path", handler)
//versus
//	router.Add("path", method.Get(handler))
//but this is easily handled by a trivial helper
//	router := NewRouter()
//	get := func(path string, handler http.Handler) {
//		router.Add(path, method.Get(handler))
//	}
//	get("path", handler)
//without sacrificing any simplicity or flexibility.
//
//For that matter, if all your routes are GET only,
//you can just write
//	http.ListenAndServe(addr, method.Get(router))
package method

import "net/http"

//Invalid is the default invalid method handler.
//It simply 405's.
//
//It is always invoked by the Get and Post helpers and by a Switch
//without the special * entry.
var Invalid http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Invalid method", 405)
})

//If returns a Handler that invokes Then if the Request Method == method
//and otherwise invokes Else.
//
//Note that method is case sensitive.
func If(method string, Then, Else http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			Then.ServeHTTP(w, r)
		} else {
			Else.ServeHTTP(w, r)
		}
	})
}

//Get is a partial application of If, namely If("GET", h, Invalid).
func Get(h http.Handler) http.Handler {
	return If("GET", h, Invalid)
}

//Post is a partial application of If, namely If("POST", h, Invalid).
func Post(h http.Handler) http.Handler {
	return If("POST", h, Invalid)
}

//Switch maps Request.Method to the appropriate Handler.
//
//If no Handler is set, it first tries to find the special "*" Handler
//and, if that is also unset, defaults to the global Invalid handler.
//
//Note that the keys are case sensitive.
type Switch map[string]http.Handler

func (s Switch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, ok := s[r.Method]
	if !ok {
		h, ok = s["*"]
		if !ok {
			h = Invalid
		}
	}
	h.ServeHTTP(w, r)
}
