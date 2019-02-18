package structeditor

import (
	"net/url"
)

type Editor interface {
	// Render the HTML for the editor UI
	Render() (string, error)
	// Run the specified operator on the data
	// referenced by the path.
	Mutate(path string, operator Operator) error
	// Create an operator described by the query params
	// in a URL
	OperatorFor(values url.Values) (Operator, error)
}

type editor struct {
	state     interface{}
	mutateUrl string
}

func NewEditor(state interface{}, mutateUrl string) Editor {
	return &editor{
		state:     state,
		mutateUrl: mutateUrl,
	}
}
