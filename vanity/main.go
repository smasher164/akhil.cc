package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	log.SetPrefix("vanity: ")
	log.SetFlags(0)
	tmpl := template.Must(template.New("").Parse(
		`<html><head>
			<meta name="go-import" content="{{.Prefix}} {{.VCS}} {{.Root}}">
			<meta name="go-source" content="{{.Prefix}} {{.Home}} {{.Dir}} {{.File}}">
		</head></html>`,
	))
	imports := []struct {
		Path,
		Prefix,
		VCS,
		Root,
		Home,
		Dir,
		File string
	}{
		{
			Path:   "/mexdown",
			Prefix: "akhil.cc/mexdown",
			VCS:    "git", Root: "https://github.com/smasher164/mexdown",
			Home: "https://github.com/smasher164/mexdown/",
			Dir:  "https://github.com/smasher164/mexdown/blob/master{/dir}",
			File: "https://github.com/smasher164/mexdown/blob/master{/dir}/{file}#L{line}",
		},
	}
	m := http.NewServeMux()
	for _, im := range imports {
		m.HandleFunc(im.Path+"/", func(w http.ResponseWriter, r *http.Request) {
			host, _, err := net.SplitHostPort(r.Host)
			if err != nil {
				host = r.Host
			}
			if strings.HasPrefix(host, "www.") {
				host = host[4:]
			}
			if r.FormValue("go-get") != "1" {
				http.Redirect(w, r, "http://godoc.org/"+host+r.URL.Path, http.StatusFound)
				return
			}
			if err := tmpl.Execute(w, im); err != nil {
				log.Println(err)
			}
		})
	}
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "www.akhil.cc" {
			http.NotFound(w, r)
			return
		}
		m.ServeHTTP(w, r)
	}))
}
