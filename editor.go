package structeditor

type Editor interface {
	// Render the HTML for the editor UI
	Render() (string, error)
}

type editor struct {
	state interface{}
}

func NewEditor(state interface{}) Editor {
	return &editor{
		state: state,
	}
}
