package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"sync"
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

type serveMux struct {
	mu       sync.RWMutex
	handlers []struct {
		r *regexp.Regexp
		h http.Handler
	}
}

func (mux *serveMux) Handle(pattern string, handler http.Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	r := regexp.MustCompile(pattern)
	mux.handlers = append(mux.handlers, struct {
		r *regexp.Regexp
		h http.Handler
	}{r, handler})
}

func (mux *serveMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	for _, m := range mux.handlers {
		if m.r.MatchString(r.Host+r.URL.Path) || m.r.MatchString(r.Host) {
			m.h.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func makeProxy(hosts []host) http.Handler {
	m := new(serveMux)
	for _, h := range hosts {
		utarget, err := url.Parse(h.Target)
		if err != nil {
			log.Fatalln(err)
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

func wildcardWhitelist(hostPatterns ...string) autocert.HostPolicy {
	pattern := strings.Join(hostPatterns, "|")
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatalln(err)
	}
	return func(_ context.Context, host string) error {
		if !re.MatchString(host) {
			return fmt.Errorf("acme/autocert: host %q not configured in HostWhitelist", host)
		}
		return nil
	}
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
		Cache:  autocert.DirCache(conf.CacheDir),
		Prompt: autocert.AcceptTOS,
		// HostPolicy: autocert.HostWhitelist(conf.Valid...),
		HostPolicy: wildcardWhitelist(conf.Valid...),
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
