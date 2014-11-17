package method

import "net/http"

var Mux *http.ServeMux

var HandleGet, HandlePost, HandlePatch, CustomInvalid http.Handler

func ExampleIf() {
	Mux.Handle("/patch", If("PATCH", HandlePatch, CustomInvalid))

	//Which is equivalent to:
	Mux.HandleFunc("/patch-inlined", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			HandlePatch.ServeHTTP(w, r)
		} else {
			CustomInvalid.ServeHTTP(w, r)
		}
	})
}

func ExampleGet() {
	Mux.Handle("/get", Get(HandleGet))

	//Which is equivalent to:
	Mux.Handle("/get-long", If("GET", HandleGet, Invalid))

	//Which is equivalent to:
	Mux.HandleFunc("/get-inlined", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			HandleGet.ServeHTTP(w, r)
		} else {
			Invalid.ServeHTTP(w, r)
		}
	})
}

func ExamplePost() {
	Mux.Handle("/post", Post(HandlePost))

	//Which is equivalent to:
	Mux.Handle("/post-long", If("POST", HandlePost, Invalid))

	//Which is equivalent to:
	Mux.HandleFunc("/post-inlined", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			HandlePost.ServeHTTP(w, r)
		} else {
			Invalid.ServeHTTP(w, r)
		}
	})
}

func ExampleSwitch() {
	Mux.Handle("/multiple-methods", Switch{
		"GET":   HandleGet,
		"POST":  HandlePost,
		"PATCH": HandlePatch,

		// * is dispatched when r.Method is not one of "GET", "POST", or "PATCH".
		// If no * were defined, the package level Invalid handler would be called.
		"*": CustomInvalid,
	})

	//Which is equivalent to:
	Mux.HandleFunc("/multiple-methods-inlined", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			HandleGet.ServeHTTP(w, r)
		case "POST":
			HandlePost.ServeHTTP(w, r)
		case "PATCH":
			HandlePatch.ServeHTTP(w, r)
		default:
			CustomInvalid.ServeHTTP(w, r)
		}
	})
}
