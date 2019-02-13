package main

import (
	"fmt"
	"github.com/fixermark/structeditor"
	"log"
	"net/http"
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

func View(e structeditor.Editor) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := e.Render()
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "%s", result)
		}
	}
}

func main() {
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

	editor := structeditor.NewEditor(demoData, "/mutate")
	http.HandleFunc("/", View(editor))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
