package middlewares

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

// RateLimit limits requests per IP.
// Default: 5 requests per 10 seconds (safe for auth endpoints)
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := clientIP(r)
		now := time.Now()

		mu.Lock()
		v, exists := visitors[ip]

		if !exists {
			visitors[ip] = &visitor{
				lastSeen: now,
				count:    1,
			}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// Reset window after 10 seconds
		if now.Sub(v.lastSeen) > 10*time.Second {
			v.lastSeen = now
			v.count = 1
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// Allow max 5 requests per window
		if v.count >= 5 {
			mu.Unlock()
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}

		v.count++
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// Cleanup stale IPs (call once at startup)
func StartRateLimitCleanup() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 5*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

// Extract real client IP (supports proxies)
func clientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		return xff
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
