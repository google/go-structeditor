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

type growable struct {
	Foo int
	Bar []int
}

func TestOperatorGrow(t *testing.T) {
	data := growable{
		Foo: 1,
		Bar: []int{
			2, 3, 4,
		},
	}

	target := growable{
		Foo: 1,
		Bar: []int{
			2, 3, 4, 0,
		},
	}

	e := NewEditor(&data, "")
	err := e.Mutate("Bar", OperatorGrow())

	if err != nil {
		t.Error("Could not mutate Bar -", err)
	}

	if !reflect.DeepEqual(data, target) {
		t.Error("Expected", target, "saw", data)
	}
}

func TestOperatorShrink(t *testing.T) {
	data := growable{
		Foo: 1,
		Bar: []int{
			2, 3, 4,
		},
	}

	target := growable{
		Foo: 1,
		Bar: []int{
			2, 3,
		},
	}

	e := NewEditor(&data, "")
	err := e.Mutate("Bar", OperatorShrink())

	if err != nil {
		t.Error("Could not mutate Bar -", err)
	}

	if !reflect.DeepEqual(data, target) {
		t.Error("Expected", target, "saw", data)
	}
}
