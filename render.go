package structeditor

import (
	"fmt"
	"reflect"
)

// Contains state used as a render is being evaluated
type renderer struct {
	nextId   int
	editable bool
}

// Render the state into HTML for serving
func (e *editor) Render() (string, error) {
	r := renderer{}
	v := reflect.ValueOf(e.state)
	if v.Kind() == reflect.Ptr {
		r.editable = true
	}
	return r.renderElement(reflect.ValueOf(e.state))
}

// Render an unknown element
func (r *renderer) renderElement(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.Uint8:
		return r.renderEditField(fmt.Sprintf("%d", v.Uint()))
	case reflect.Uint16:
		return r.renderEditField(fmt.Sprintf("%d", v.Uint()))
	case reflect.Uint32:
		return r.renderEditField(fmt.Sprintf("%d", v.Uint()))
	case reflect.Uint64:
		return r.renderEditField(fmt.Sprintf("%d", v.Uint()))
	case reflect.Uint:
		return r.renderEditField(fmt.Sprintf("%d", v.Uint()))
	case reflect.Int8:
		return r.renderEditField(fmt.Sprintf("%d", v.Int()))
	case reflect.Int16:
		return r.renderEditField(fmt.Sprintf("%d", v.Int()))
	case reflect.Int32:
		return r.renderEditField(fmt.Sprintf("%d", v.Int()))
	case reflect.Int64:
		return r.renderEditField(fmt.Sprintf("%d", v.Int()))
	case reflect.Int:
		return r.renderEditField(fmt.Sprintf("%d", v.Int()))
	case reflect.Float32:
		return r.renderEditField(fmt.Sprintf("%f", v.Float()))
	case reflect.Float64:
		return r.renderEditField(fmt.Sprintf("%f", v.Float()))
	case reflect.Bool:
		return r.renderEditField(fmt.Sprintf("%v", v.Bool()))
	case reflect.String:
		return r.renderEditField(v.String())
	default:
		return r.renderComposite(v)
		// Next: render complex type function, struct
		// render function

	}
}

// Render a composite element type (any type containing another type): struct,
// array, slice, map, &c
func (r *renderer) renderComposite(elt reflect.Value) (string, error) {
	switch elt.Kind() {
	case reflect.Struct:
		return r.renderStruct(elt)
	case reflect.Array:
		return r.renderArray(elt)
	case reflect.Slice:
		return r.renderSlice(elt)
	case reflect.Ptr:
		return r.renderPtr(elt)
	default:
		return "", fmt.Errorf("Unknown composite render type: %v", elt.Kind())
	}
}

// Render a struct
func (r *renderer) renderStruct(v reflect.Value) (string, error) {
	t := v.Type()

	result := fmt.Sprintf("<div>%s {<ul>", t.Name())
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		result += fmt.Sprintf("<li>%s: ", sf.Name)
		subvalue := v.Field(i)
		rendered, err := r.renderElement(subvalue)
		if err != nil {
			return "", err
		}
		result += fmt.Sprintf("%s,</li>", rendered)
	}
	result += "}</ul></div>"
	return result, nil
}

func (r *renderer) renderArray(v reflect.Value) (string, error) {
	t := v.Type()
	innerType := t.Elem()
	len := v.Len()
	result := fmt.Sprintf("<div>[%d]%s {<ul>", len, innerType)
	for i := 0; i < len; i++ {
		subelem := v.Index(i)
		subtext, err := r.renderElement(subelem)
		if err != nil {
			return "", err
		}
		result += fmt.Sprintf("<li>%s,</li>", subtext)
	}
	result += "}</ul></div>"
	return result, nil
}

func (r *renderer) renderSlice(v reflect.Value) (string, error) {
	t := v.Type()
	innerType := t.Elem()
	len := v.Len()
	result := fmt.Sprintf("<div>[]%s {<ul>", innerType)
	for i := 0; i < len; i++ {
		subelem := v.Index(i)
		subtext, err := r.renderElement(subelem)
		if err != nil {
			return "", err
		}
		result += fmt.Sprintf("<li>%s,</li>", subtext)
	}
	result += "}</ul></div>"
	return result, nil
}

func (r *renderer) renderPtr(v reflect.Value) (string, error) {
	if v.IsNil() {
		return "nil", nil
	}
	innerValue := v.Elem()
	innerText, err := r.renderElement(innerValue)
	return fmt.Sprintf("&%s", innerText), err
}

func (r *renderer) getNextId() string {
	id := r.nextId
	r.nextId += 1
	return fmt.Sprintf("input-%d", id)
}

func (r *renderer) renderEditField(value string) (string, error) {
	return fmt.Sprintf("<input type='text' id='%s' value='%s'>", r.getNextId(), value), nil
}
