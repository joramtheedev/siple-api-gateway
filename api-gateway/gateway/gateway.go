package gateway

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"simple-gateway/config"
)

// New builds an http.Handler that reverse-proxies each route prefix
// to its backend service.
func New(routes []config.Route, apiKey string) (http.Handler, error) {
	mux := http.NewServeMux()

	for _, rt := range routes {
		target, err := url.Parse(rt.Target)
		if err != nil {
			return nil, fmt.Errorf("bad target %q: %w", rt.Target, err)
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("proxy error for %s: %v", r.URL.Path, err)
			http.Error(w, "bad gateway", http.StatusBadGateway)
		}

		// /users/42 -> backend receives /42
		mux.Handle(rt.Prefix+"/", http.StripPrefix(rt.Prefix, proxy))
		log.Printf("route %s -> %s", rt.Prefix, rt.Target)
	}

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return logging(requireAPIKey(apiKey, mux)), nil
}
