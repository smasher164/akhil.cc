package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func main() {
	m := http.NewServeMux()
	m.Handle("/", http.FileServer(http.Dir("/root")))
	m.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		var m struct {
			Name  string
			Bytes string
		}
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		b, err := base64.StdEncoding.DecodeString(m.Bytes)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if err := ioutil.WriteFile(filepath.Join("/root/", m.Name), b, 0666); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "www.static.akhil.cc" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		m.ServeHTTP(w, r)
	}))
}
