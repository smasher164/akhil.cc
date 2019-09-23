package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.SetPrefix("www.staging.akhil.cc: ")
	log.SetFlags(0)
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Your Host: %s", r.Host)
	}))
}
