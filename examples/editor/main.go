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
	// Note difference from examples/viewer/main.go: demoData is a pointer
	// to the structure, so the structure will be editable.
	demoData := &example{
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
