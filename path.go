// Copyright 2019 Google, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package structeditor

import (
	"fmt"
	"strconv"
	"strings"
)

// Path to a specific variable
type Path struct {
	// Only one is true:
	// name is not ""
	// index has meaning

	// Name of struct field in path
	Name string
	// Current variable is array or slice and should be indexed
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

// Appends the specified newElement to the path and
// returns the new root of the path.
func (p *Path) Append(newElement *Path) *Path {
	var prevEl *Path
	var curEl *Path
	for prevEl, curEl = nil, p; curEl != nil; prevEl, curEl = curEl, curEl.Next {
	}
	if prevEl != nil {
		prevEl.Next = newElement
		return p
	} else {
		return newElement
	}
}

// Removes the last element from the path and returns
// the new root of the path
func (p *Path) RemoveLast() *Path {
	if p == nil || p.Next == nil {
		return nil
	}
	var prevEl *Path
	var curEl *Path
	for prevEl, curEl = p, p.Next; curEl.Next != nil; prevEl, curEl = curEl, curEl.Next {
	}
	prevEl.Next = nil
	return p
}

type VisitingFunc func(updatedPath *Path)

// Attaches the specified element to the path, runs the specified function, and
// then detaches the specified element. Convenience function for processes that,
// for example, need to visit every field in a struct.
func (p *Path) Visiting(element *Path, doing VisitingFunc) {
	p = p.Append(element)
	doing(p)
	p = p.RemoveLast()
}

func (p *Path) String() string {
	if p == nil {
		return ""
	}
	subpath := p.Next.String()
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
