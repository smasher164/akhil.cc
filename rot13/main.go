package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func rot13(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			// Rotate lowercase letters 13 places.
			if r >= 'm' {
				return r - 13
			} else {
				return r + 13
			}
		} else if r >= 'A' && r <= 'Z' {
			// Rotate uppercase letters 13 places.
			if r >= 'M' {
				return r - 13
			} else {
				return r + 13
			}
		}
		// Do nothing.
		return r
	}, s)
}

func main() {
	log.SetPrefix("rot13: ")
	log.SetFlags(0)
	m := mux.NewRouter().SkipClean(true)
	m.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == rot13("/") {
			fmt.Fprintf(w, "pass a rot 13'd URL")
			return
		}
		// log.Printf("URL=%v", r.URL)
		// log.Printf("RequestURI=%v", r.RequestURI)
		uri := strings.TrimPrefix(r.RequestURI, "/")
		if !strings.HasPrefix(uri, rot13("http://")) && !strings.HasPrefix(uri, rot13("https://")) {
			uri = rot13("https://") + uri
		}
		// log.Printf("Trimmed URI=%v", uri)
		http.Redirect(w, r, rot13(uri), http.StatusFound)
	})
	m.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(nil)
	})
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "www.rot13.akhil.cc" {
			http.NotFound(w, r)
			return
		}
		m.ServeHTTP(w, r)
	}))
}
