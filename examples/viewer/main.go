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

package main

import (
	"log"
	"net/http"

	"github.com/google/structeditor"
)

type customer struct {
	Name    string
	Balance float32
}
type example struct {
	Company       string
	Id            int
	BillingActive bool
	Customers     []customer
}

func main() {
	demoData := example{
		Company:       "ExampleCo",
		Id:            123,
		BillingActive: true,
		Customers: []customer{
			customer{
				Name:    "Bob",
				Balance: 15.50,
			}, customer{
				Name:    "Shepherd's Plumbing And Fences",
				Balance: -5,
			},
		},
	}

	structeditor.ServeEditor(demoData, "/", http.DefaultServeMux)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
