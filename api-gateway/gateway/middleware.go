package gateway

import (
	"crypto/subtle"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now
		rec := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}

func requireApiKey(key string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}
		got := r.Header.Get("X-API-Key")
		if subtle.ConstantTimeCompare([]byte(got), []byte(key)) != 1 {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type client struct {
	count       int
	windowStart time.Time
}

// rateLimit is a middleware that limits the number of requests per client IP address.
func rateLimit(limit int, window time.Duration, next http.Handler) http.Handler {
	var mu sync.Mutex
	clients := make(map[string]*client)

	//periodically drop clients whose window has expired to prevent memory leak
	go func() {
		for range time.Tick(window) {
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.indowStart) > window {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr // fallback to the whole address if splitting fails
		}
		mu.Lock()
		c, ok := clients[ip]
		if !ok || time.Since(c.windowStart) > window {
			c = &client{windowStart: time.Now()}
			clients[ip] = c
		}
		c.count++
		allowed := c.count <= limit
		mu.Unlock()

		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(int(window.Seconds())))
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})

}
