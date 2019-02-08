package structeditor

import (
	"testing"
)

func TestEmptyStringToPath(t *testing.T) {
	p, err := StringToPath("")
	if err != nil {
		t.Error(err)
	}
	if p != nil {
		t.Error("p was not empty, was", p)
	}
}

func TestStringToPath(t *testing.T) {
	data := []struct {
		input    string
		expected *Path
	}{
		{
			"base", &Path{
				Name: "base",
			},
		},
		{
			"1", &Path{
				Index: 1,
			},
		},
		{
			"base.1", &Path{
				Name: "base",
				Next: &Path{
					Index: 1,
				},
			},
		},
		{
			"1.customer.name", &Path{
				Index: 1,
				Next: &Path{
					Name: "customer",
					Next: &Path{
						Name: "name",
					},
				},
			},
		},
	}

	for _, step := range data {
		topResult, err := StringToPath(step.input)
		if err != nil {
			t.Error(step.input, err)
			continue
		}
		if topResult == nil && step.expected != nil {
			t.Error(step.input, 0, ": result was nil but was expected to be non-nil")
		}
		for result, expected, depth := topResult, step.expected, 0; result != nil; result, expected, depth = result.Next, expected.Next, depth+1 {
			if result == nil && expected != nil {
				t.Error(step.input, depth, ": Expected non-nil result, was nil")
			}
			if result != nil && expected == nil {
				t.Error(step.input, depth, ": Expected nil, was not nil")
			}
			if result.Name != expected.Name {
				t.Error(step.input, depth, ": Expected Name to be", expected.Name, ", was", result.Name)
			}
			if result.Index != expected.Index {
				t.Error(step.input, depth, ": Expected Index to be", expected.Index, ", was", result.Index)
			}
		}
	}
}

func TestPathToString(t *testing.T) {
	data := []struct {
		input    *Path
		expected string
	}{
		{
			nil, "",
		},
		{
			&Path{
				Name: "base",
			}, "base",
		},
		{
			&Path{
				Index: 3,
			}, "3",
		},
		{
			&Path{
				Name: "customers",
				Next: &Path{
					Index: 3,
					Next: &Path{
						Name: "balance",
					},
				},
			}, "customers.3.balance",
		},
	}

	for _, step := range data {
		result := step.input.ToString()
		if result != step.expected {
			t.Error("Expected", step.expected, "saw", result)
		}
	}
}
