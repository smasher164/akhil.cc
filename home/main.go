package main

import (
	"log"
	"net/http"
	"os"

	"akhil.cc/mexdown/gen/html"
	"akhil.cc/mexdown/parser"
)

func main() {
	log.SetPrefix("www.akhil.cc: ")
	log.SetFlags(0)
	file, err := os.Open("home.xd")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	f, err := parser.Parse(file)
	if err != nil {
		log.Print(err)
		return
	}
	out, err := html.Gen(f).Output()
	if err != nil {
		log.Print(err)
		return
	}
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(out)
	})
	m.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(nil)
	})
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "www.akhil.cc" {
			http.NotFound(w, r)
			return
		}
		m.ServeHTTP(w, r)
	}))
}
