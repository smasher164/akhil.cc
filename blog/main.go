package main

import (
	"bytes"
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/feeds"
)

var site string
var file string
var dir string
var static string

func init() {
	flag.StringVar(&site, "site", "www.blog.akhil.cc", "primary host name for site")
	flag.StringVar(&file, "conf", "posts.toml", "path to posts file")
	flag.StringVar(&dir, "posts", "/posts", "path to posts directory")
	flag.StringVar(&static, "static", "/static", "path to static file resources")
}

type blog struct {
	cache map[string][]byte
	posts []struct {
		Route       string
		File        string
		Title       string
		Description string
		Created     time.Time
	}
	feed *feeds.Feed
}

const stub = `<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8"><meta name="viewport" content="width=device-width, minimum-scale=1.0">
        <meta name="Description" content="{{.Description}}">
        <title>{{.Title}}</title>
        <link href="/static/blog.css" rel="stylesheet">
    </head>
    <body>
        {{printf "%s" .Body}}
    </body>
</html>
`

var tmpl = template.Must(template.New("").Parse(stub))

func (b *blog) init(posts map[string]string) {
	// load each post from file and into cache
	for _, p := range b.posts {
		p := p
		fpath := posts[p.Route]
		file, err := os.Open(fpath)
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		out, err := ioutil.ReadAll(file)
		if err != nil {
			log.Print(err)
			return
		}
		buf := new(bytes.Buffer)
		desc := fmt.Sprintf("Author: Akhil Indurti, Created: %v", p.Created)
		if len(p.Description) != 0 {
			desc += ", Description: " + p.Description
		}
		res := struct {
			Title       string
			Description string
			Body        []byte
		}{Title: p.Title, Description: desc, Body: out}
		tmpl.Execute(buf, res)
		b.cache[p.Route] = buf.Bytes()
	}
	// create atom feed from posts
	now := time.Now()
	b.feed = &feeds.Feed{
		Title:   "blog.akhil.cc",
		Link:    &feeds.Link{Href: "https://www.blog.akhil.cc"},
		Author:  &feeds.Author{Name: "Akhil Indurti", Email: "aindurti@gmail.com"},
		Created: now,
	}
	for _, post := range b.posts {
		b.feed.Items = append(b.feed.Items, &feeds.Item{
			Title:       post.Title,
			Link:        &feeds.Link{Href: site + post.Route},
			Description: post.Description,
			Created:     post.Created,
		})
	}
}

func (b *blog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, ok := b.cache[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Write(body)
}

func main() {
	flag.Parse()
	log.SetPrefix(site + ": ")
	log.SetFlags(0)
	var conf struct {
		Posts []struct {
			Route       string
			File        string
			Title       string
			Description string
			Created     time.Time
		}
	}
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		log.Fatalln(err)
	}
	posts := make(map[string]string)
	for _, p := range conf.Posts {
		posts[p.Route] = dir + "/" + p.File
	}
	blog := &blog{cache: make(map[string][]byte), posts: conf.Posts}
	blog.init(posts)
	m := http.NewServeMux()
	m.Handle("/", blog)
	h := http.StripPrefix("/static/", http.FileServer(http.Dir(static)))
	m.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		fpath := path.Join(static, strings.TrimPrefix(r.URL.Path, "/static/"))
		info, err := os.Stat(fpath)
		if err != nil {
			log.Fatal(err)
		}
		etag := fmt.Sprintf("%x", sha256.Sum256([]byte(info.ModTime().String())))
		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, etag) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		w.Header().Set("Etag", etag)
		w.Header().Set("Cache-Control", "no-cache")
		if strings.HasPrefix(path.Ext(r.URL.Path), "woff2") {
			w.Header().Set("Content-Type", "application/font-woff2")
		}
		h.ServeHTTP(w, r)
	})
	m.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(nil)
	})
	m.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Sitemap: https://www.blog.akhil.cc/feed.atom"))
	})
	m.HandleFunc("/feed.atom", func(w http.ResponseWriter, r *http.Request) {
		blog.feed.WriteAtom(w)
	})
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != site {
			http.NotFound(w, r)
			return
		}
		m.ServeHTTP(w, r)
	}))
}
