package main

import (
	"net/http"
)

func main() {
	m := http.NewServeMux()
	m.Handle("/", http.FileServer(http.Dir("/root")))
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "www.static.akhil.cc" {
			http.NotFound(w, r)
			return
		}
		m.ServeHTTP(w, r)
	}))
}
