package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

func main() {
	log.SetPrefix("www.akhil.cc: ")
	log.SetFlags(0)
	cache := make(map[string][]byte)
	for _, name := range []string{"home.html", "plate.html", "stayhome.html"} {
		b, err := ioutil.ReadFile(name)
		if err != nil {
			log.Fatal(err)
			return
		}
		cache[name] = b
	}
	const stub = `<html>
<head>
	<meta charset="UTF-8"><meta name="viewport" content="width=device-width, maximum-scale=1.0">
<style>
html {
	position: relative;
	width: 100%;
	height: 100%;
}
</style>
</head>
<body>{{printf "%s" .}}</body></html>`
	tmpl := template.Must(template.New("").Parse(stub))
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, cache["home.html"])
	})
	m.HandleFunc("/plate", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, cache["plate.html"])
	})
	m.HandleFunc("/stayhome", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, cache["stayhome.html"])
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
