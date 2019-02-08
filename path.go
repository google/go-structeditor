package structeditor

import (
	"fmt"
	"strconv"
	"strings"
)

// Path to a piece of a value
type Path struct {
	// Only one is filled: index or name
	Name  string
	Index int
	// if nil, this Path refers to the top-level element
	Next *Path
}

// Converts a string to a Path pointer. If the pointer is nil,
// the Path refers to the top-level ("current") element.
func StringToPath(s string) (*Path, error) {
	if s == "" {
		return nil, nil
	}
	sliced := strings.Split(s, ".")
	return sliceToPath(sliced)
}

// Converts a string slice to a Path pointer. If the pointer is nil,
// the Path refers to the top-level ("current") element.
func sliceToPath(slice []string) (*Path, error) {
	var first *Path
	var cur *Path
	for _, elt := range slice {
		pathEl, err := encodePath(elt)
		if err != nil {
			return nil, err
		}
		if first == nil {
			first = pathEl
		} else {
			cur.Next = pathEl
		}
		cur = pathEl

	}
	return first, nil
}

// Encodes path part
func encodePath(s string) (*Path, error) {
	if strings.IndexAny(s, "0123456789") == 0 {
		result, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		return &Path{
			Index: result,
		}, nil
	}

	return &Path{
		Name: s,
	}, nil
}

func (p *Path) ToString() string {
	if p == nil {
		return ""
	}
	subpath := p.Next.ToString()
	elt := p.Name
	if elt == "" {
		elt = fmt.Sprintf("%d", p.Index)
	}
	if subpath == "" {
		return elt
	} else {
		return elt + "." + subpath
	}
}
