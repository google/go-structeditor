package structeditor

import (
	"reflect"
	"testing"
)

type equality struct {
	t         *testing.T
	compareTo string
}

func (e *equality) Do(v reflect.Value) error {
	if e.compareTo != v.String() {
		e.t.Error("Expected", e.compareTo, "\nsaw", v)
	}
	return nil
}

func (e *equality) ModifiesPointer() bool {
	return false
}

func equalityOf(t *testing.T, v string) *equality {
	return &equality{
		t, v,
	}
}

type testEmployee struct {
	Name string
	Id   string
}

func TestFindValue(t *testing.T) {
	scanning := struct {
		Foo       string
		Bar       string
		Employees []testEmployee
		Boss      *testEmployee
	}{
		"5",
		"hello",
		[]testEmployee{
			{
				Name: "Bob",
				Id:   "A",
			}, {
				Name: "Sue",
				Id:   "B",
			},
		},
		&testEmployee{
			"Snake",
			"0",
		},
	}

	data := []struct {
		path        string
		shouldEqual Operator
	}{
		{
			"Foo",
			equalityOf(t, "5"),
		}, {
			"Bar",
			equalityOf(t, "hello"),
		}, {
			"Employees.0.Name",
			equalityOf(t, "Bob"),
		}, {
			"Employees.1.Id",
			equalityOf(t, "B"),
		}, {
			"Boss.Name",
			equalityOf(t, "Snake"),
		},
	}

	e := NewEditor(scanning, "")
	for _, step := range data {
		err := e.Mutate(step.path, step.shouldEqual)
		if err != nil {
			t.Error(step.path, "- saw error:", err)
		}
	}
}

type modify struct {
	Foo int
	Bar string
	Baz bool
}

func TestModifyValue(t *testing.T) {
	data := modify{
		Foo: 5,
		Bar: "hello",
		Baz: false,
	}

	target := modify{
		Foo: 7,
		Bar: "hi",
		Baz: true,
	}

	mutations := []struct {
		path     string
		newValue string
	}{
		{"Foo", "7"},
		{"Bar", "hi"},
		{"Baz", "true"},
	}

	e := NewEditor(&data, "")

	for _, mutation := range mutations {
		err := e.Mutate(
			mutation.path,
			OperatorSet(mutation.newValue))
		if err != nil {
			t.Error(mutation.path, "-", err)
		}
	}

	if !reflect.DeepEqual(data, target) {
		t.Error("Expected", target, "saw", data)
	}
}
