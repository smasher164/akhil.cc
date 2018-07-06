package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"text/template"

	"akhil.cc/mexdown/gen/html"
	"akhil.cc/mexdown/parser"
	"github.com/BurntSushi/toml"
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
	posts map[string]string
	cache map[string][]byte
}

const stub = `<html>
    <head>
        <meta charset="UTF-8"><meta name="viewport" content="width=device-width, maximum-scale=1.0">
        <link href="/static/blog.css" rel="stylesheet">
    </head>
    <body>
        {{printf "%s" .}}
    </body>
</html>
`

var tmpl = template.Must(template.New("").Parse(stub))

func (b *blog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fpath, ok := b.posts[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	if b, ok := b.cache[r.URL.Path]; ok {
		tmpl.Execute(w, b)
		return
	}
	file, err := os.Open(fpath)
	if err != nil {
		log.Println(err)
		http.Error(w, "couldn't load post", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	f, err := parser.Parse(file)
	if err != nil {
		log.Print(err)
		http.Error(w, "couldn't load post", http.StatusInternalServerError)
		return
	}
	out, err := html.Gen(f).Output()
	if err != nil {
		log.Print(err)
		http.Error(w, "couldn't load post", http.StatusInternalServerError)
		return
	}
	b.cache[r.URL.Path] = out
	tmpl.Execute(w, out)
}

type cache struct {
	m          sync.Mutex
	route2etag map[string]string
}

func (c *cache) control(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.m.Lock()
		defer c.m.Unlock()
		h.ServeHTTP(w, r)
		if strings.HasPrefix(path.Ext(r.URL.Path), "woff2") {
			w.Header().Add("Content-Type", "application/font-woff2")
		}
		// tag := c.route2etag[r.URL.Path]
		// if match := r.Header.Get("If-Modified-Since"); match != "" {
		// 	if match == tag {
		// 		w.WriteHeader(http.StatusNotModified)
		// 		return
		// 	}
		// }
		// if match := r.Header.Get("If-Range"); match != "" {
		// 	if match == tag {
		// 		w.WriteHeader(http.StatusNotModified)
		// 		return
		// 	}
		// }
		// // tag = time.Now().Format(time.RFC1123)
		// // data := r.URL.Path + time.Now().String()
		// // hash := xxhash.Sum64String(data)
		// // tag = fmt.Sprintf("%08x", hash)
		// h.ServeHTTP(w, r)
		// w.Header().Set("Cache-Control", "no-cache")
		// tag = w.Header().Get("Last-Modified")
		// if tag == "" {
		// 	tag = time.Now().Format(time.RFC1123)
		// }
		// c.route2etag[r.URL.Path] = tag
	})
}

func main() {
	flag.Parse()
	log.SetPrefix(site + ": ")
	log.SetFlags(0)
	var conf struct {
		Posts []struct {
			Route string
			File  string
		}
	}
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		log.Fatalln(err)
	}
	posts := make(map[string]string)
	for _, p := range conf.Posts {
		posts[p.Route] = dir + "/" + p.File
	}
	cache := &cache{route2etag: make(map[string]string)}
	blog := &blog{posts: posts, cache: make(map[string][]byte)}
	m := http.NewServeMux()
	m.Handle("/", cache.control(blog))
	m.Handle("/static/", cache.control(http.StripPrefix("/static/", http.FileServer(http.Dir(static)))))
	m.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(nil)
	})
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != site {
			http.NotFound(w, r)
			return
		}
		m.ServeHTTP(w, r)
	}))
}
