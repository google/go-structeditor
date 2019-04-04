// Copyright 2019 Google LLC
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
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

type Operator interface {
	Do(v reflect.Value) error
	// Return true if this Operator modifies a pointer
	ModifiesPointer() bool
}

func (e *editor) Mutate(path string, operator Operator) error {
	p, err := StringToPath(path)
	if err != nil {
		return err
	}

	v, err := e.findValueToChange(p, reflect.ValueOf(e.state), operator.ModifiesPointer())
	if err != nil {
		return err
	}

	return operator.Do(v)
}

func (e *editor) findValueToChange(p *Path, v reflect.Value, modifiesPtr bool) (reflect.Value, error) {
	if p == nil {
		return v, nil
	}
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if modifiesPtr {
			return v, nil
		} else {
			dereferenced := v.Elem()
			if !dereferenced.IsValid() {
				return reflect.Value{}, errors.New("Attempted to dereference nil pointer.")
			}
			return e.findValueToChange(p, dereferenced, modifiesPtr)
		}
	case reflect.Struct:
		if p.Name == "" {
			return reflect.Value{}, errors.New("Attempted numeric indexing on a struct or interface.")
		}
		el := v.FieldByName(p.Name)
		if !el.IsValid() {
			return reflect.Value{}, errors.New("No field by name '" + p.Name + "'")
		}
		return e.findValueToChange(p.Next, el, modifiesPtr)
	case reflect.Array, reflect.Slice:
		if p.Name != "" {
			return reflect.Value{}, errors.New("Attempted to index into array or slice using a name string '" + p.Name + "'.")
		}
		if v.Len() <= p.Index {
			return reflect.Value{}, fmt.Errorf("Attempted to fetch element %d, but array or slice is length %d", p.Index, v.Len())
		}
		el := v.Index(p.Index)
		return e.findValueToChange(p.Next, el, modifiesPtr)
	}
	return reflect.Value{}, errors.New("Could not follow path through element with type '" + v.Kind().String() + "'")

}

/// Operators

// Set a value to a new value (indicated by a string)
type operatorSet struct {
	newValue string
}

func OperatorSet(newValue string) Operator {
	return &operatorSet{newValue}
}

func (o *operatorSet) ModifiesPointer() bool {
	return false
}

func (o *operatorSet) Do(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		newValue, err := strconv.ParseInt(o.newValue, 0, 64)
		if err != nil {
			return err
		}
		v.SetInt(newValue)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		newValue, err := strconv.ParseUint(o.newValue, 0, 64)
		if err != nil {
			return err

		}
		v.SetUint(newValue)
		return nil
	case reflect.Float32, reflect.Float64:
		newValue, err := strconv.ParseFloat(o.newValue, 64)
		if err != nil {
			return err
		}
		v.SetFloat(newValue)
	case reflect.String:
		v.SetString(o.newValue)
		return nil
	case reflect.Bool:
		newValue, err := strconv.ParseBool(o.newValue)
		if err != nil {
			return err
		}
		v.SetBool(newValue)
		return nil
	default:
		return fmt.Errorf("Unable to set value on type %v", v.Kind())
	}
	return nil
}

// Grow a slice by one element (the element takes on the zero value for th slice)
type operatorGrow struct {
}

func OperatorGrow() Operator {
	return &operatorGrow{}
}

func (o *operatorGrow) ModifiesPointer() bool {
	return false
}

func (o *operatorGrow) Do(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Slice:
		if v.Len() >= v.Cap() {
			v.Set(doubleCapacity(v))
		}
		v.SetLen(v.Len() + 1)
	default:
		return fmt.Errorf("Unable to grow type %v", v.Kind())
	}
	return nil
}

// Shrink a slice by one element. Does nothing if the slice's size is already
// zero.
type operatorShrink struct{}

func OperatorShrink() Operator {
	return &operatorShrink{}
}

func (o *operatorShrink) ModifiesPointer() bool {
	return false
}

func (o *operatorShrink) Do(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Slice:
		if v.Len() > 0 {
			v.SetLen(v.Len() - 1)
		}
	default:
		return fmt.Errorf("Unable to shrink type %v", v.Kind())
	}
	return nil
}

// doubleCapacity takes a slice (as a Value) and returns a copy of the slice
// with the capacity doubled
func doubleCapacity(sliceValue reflect.Value) reflect.Value {
	newSlice := reflect.MakeSlice(sliceValue.Type(), sliceValue.Len(), 2*sliceValue.Cap())
	reflect.Copy(newSlice, sliceValue)
	return newSlice
}

/// end operators

func (e *editor) OperatorFor(values url.Values) (Operator, error) {
	operatorName := values.Get("operator")
	switch operatorName {
	case "set":
		return OperatorSet(values.Get("value")), nil
	case "grow":
		return OperatorGrow(), nil
	case "shrink":
		return OperatorShrink(), nil
	}
	return nil, errors.New("Unable to build Operator named '" + operatorName + "'")
}
