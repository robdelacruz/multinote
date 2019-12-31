package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := "8000"

	http.HandleFunc("/", rootHandler())
	fmt.Printf("Listening on %s...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	log.Fatal(err)
}

func rootHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		fmt.Fprintf(w, "hello.\n")
	}
}
