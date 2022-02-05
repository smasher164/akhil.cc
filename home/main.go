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
	font-family: sans-serif;
	white-space: pre-line;
	padding-left: 1rem;
	font-size: 1.1rem;
}
.bullet {
	padding-bottom: 0.5rem;
	line-height: 1.1rem;
}
</style>
</head>
<body>{{printf "%s" .}}</body></html>`
	tmpl := template.Must(template.New("").Parse(stub))
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, cache["home.html"])
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
