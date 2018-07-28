package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/NYTimes/gziphandler"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

var site string
var file string

func init() {
	flag.StringVar(&site, "site", "www.akhil.cc", "primary host name for site")
	flag.StringVar(&file, "conf", "conf.toml", "path to configuration file")
}

func hsts(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		h.ServeHTTP(w, r)
	})
}

type host struct {
	Route  string
	Target string
}

func makeProxy(hosts []host) http.Handler {
	m := http.NewServeMux()
	for _, h := range hosts {
		utarget, err := url.Parse(h.Target)
		if err != nil {
			log.Fatalln(err)
		}
		if strings.HasPrefix(h.Route, "/") {
			h.Route = site + h.Route
		}
		m.Handle(h.Route, httputil.NewSingleHostReverseProxy(utarget))
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Host, "www.") {
			r.Host = "www." + r.Host
			redURL := "https://" + r.Host + r.URL.String()
			http.Redirect(w, r, redURL, http.StatusFound)
		} else {
			m.ServeHTTP(w, r)
		}
	})
}

func main() {
	flag.Parse()
	log.SetPrefix(site + ": ")
	log.SetFlags(0)
	var conf struct {
		CacheDir string
		Valid    []string
		Email    string
		Hosts    []host
	}
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		log.Fatalln(err)
	}
	cert := &autocert.Manager{
		Cache:      autocert.DirCache(conf.CacheDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(conf.Valid...),
		Email:      conf.Email,
	}
	var redirMux http.ServeMux
	redirMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		redURL := "https://" + r.Host + r.URL.String()
		http.Redirect(w, r, redURL, http.StatusFound)
	})
	redir := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      cert.HTTPHandler(&redirMux),
		Addr:         ":80",
	}
	proxy := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      gziphandler.GzipHandler(hsts(makeProxy(conf.Hosts))),
		Addr:         ":443",
		TLSConfig: &tls.Config{
			GetCertificate:           cert.GetCertificate,
			PreferServerCipherSuites: true,
			NextProtos: []string{
				"h2", "http/1.1", // enable HTTP/2
				acme.ALPNProto, // enable tls-alpn ACME challenges
			},
		},
	}
	go func() {
		if err := redir.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()
	if err := proxy.ListenAndServeTLS("", ""); err != nil {
		log.Fatalln(err)
	}
}
