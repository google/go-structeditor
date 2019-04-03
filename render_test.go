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
	"fmt"
	"testing"
)

type exampleStruct struct {
	myString string
	myNumber int
	myBool   bool
}

func inputString(value string, index int) string {
	return fmt.Sprintf("<input type='text' id='input-%d' value='%s'>", index, value)

}

func primitiveEditString(value string, path string, index int) string {
	return fmt.Sprintf("<input type='text' id='input-%d' value='%s'><button onclick=\"update('%s', 'input-%d')\">change</button>", index, value, path, index)
}

func TestRenderElement(t *testing.T) {
	addressableValue := 5

	data := []struct {
		input  interface{}
		result string
	}{
		{3, inputString("3", 0)},
		{int32(5), inputString("5", 0)},
		{uint64(10), inputString("10", 0)},
		{3.0, inputString("3.000000", 0)},
		{false, inputString("false", 0)},
		{"hi", inputString("hi", 0)},
		{[3]int{1, 2, 3},
			"<div>[3]int {<ul><li>" +
				inputString("1", 0) +
				",</li><li>" +
				inputString("2", 1) +
				",</li><li>" +
				inputString("3", 2) +
				",</li>}</ul></div>"},
		{[]int{1, 2, 3},
			"<div>[]int {<ul><li>" +
				inputString("1", 0) +
				",</li><li>" +
				inputString("2", 1) +
				",</li><li>" +
				inputString("3", 2) +
				",</li>}</ul></div>"},
		{&addressableValue, "&" + primitiveEditString("5", "", 0)},
	}

	for _, step := range data {

		e := &editor{state: step.input}
		result, err := e.unwrappedRender()

		if err != nil {
			t.Error("Rendering error:", err)
		}

		if result != step.result {
			t.Error("Expected", step.result, "saw", result)
		}
	}
}

func TestRenderStruct(t *testing.T) {
	testCase := exampleStruct{
		myString: "hello",
		myNumber: 5,
		myBool:   true,
	}

	e := editor{state: testCase}
	result, err := e.unwrappedRender()

	if err != nil {
		t.Error("Rendering error:", err)
	}
	expected := "<div>exampleStruct {<ul><li>myString: " + inputString("hello", 0) +
		",</li><li>myNumber: " + inputString("5", 1) +
		",</li><li>myBool: " + inputString("true", 2) +
		",</li>}</ul></div>"

	if result != expected {
		t.Error("Expected", expected, "saw", result)
	}

}
