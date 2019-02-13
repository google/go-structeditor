package structeditor

type Editor interface {
	// Render the HTML for the editor UI
	Render() (string, error)
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
