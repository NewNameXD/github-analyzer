package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	limit    time.Duration
}

type visitor struct {
	lastSeen time.Time
}

func NewRateLimiter(limit time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
	}

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for ip, v := range rl.visitors {
		if time.Since(v.lastSeen) > rl.limit*2 {
			delete(rl.visitors, ip)
		}
	}
}

func (rl *rateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{}
		rl.visitors[ip] = v
	}

	return v
}

func (rl *rateLimiter) isAllowed(ip string) (bool, int) {
	v := rl.getVisitor(ip)

	now := time.Now()
	timeSince := now.Sub(v.lastSeen)

	if timeSince < rl.limit {
		remainingSeconds := int(rl.limit.Seconds() - timeSince.Seconds())
		return false, remainingSeconds
	}

	return true, 0
}

func (rl *rateLimiter) updateLastSeen(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if v, exists := rl.visitors[ip]; exists {
		v.lastSeen = time.Now()
	}
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		if ip := net.ParseIP(forwarded); ip != nil {
			return forwarded
		}
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		if ip := net.ParseIP(realIP); ip != nil {
			return realIP
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

func RateLimitMiddleware(limiter *rateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/evaluate" {
				next.ServeHTTP(w, r)
				return
			}

			ip := getIP(r)
			allowed, remainingSeconds := limiter.isAllowed(ip)

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"success":          false,
					"error":            "Rate limit exceeded. Please wait before making another request.",
					"remainingSeconds": remainingSeconds,
				})
				return
			}

			limiter.updateLastSeen(ip)

			next.ServeHTTP(w, r)
		})
	}
}
