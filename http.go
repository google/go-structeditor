package structeditor

import (
	"fmt"
	"net/http"
)

// Handlers for serving the view interface via HTTP and handling mutation requests

// ViewHandler is an HTTP request handler that returns the structeditor user
// interface.
func (e *editor) ViewHandler(w http.ResponseWriter, r *http.Request) {
	result, err := e.Render()
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "%s", result)
	}
}

// MutateHandler is an HTTP request handler that modifies the editable state in
// response to mutation operations (usually generated by the UI built by
// ViewHandler).
func (e *editor) MutateHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	path := values.Get("path")
	operator, err := e.OperatorFor(values)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = e.Mutate(path, operator)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Error(w, "", 200)
}