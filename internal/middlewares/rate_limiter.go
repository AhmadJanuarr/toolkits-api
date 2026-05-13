package middlewares

import (
	"net/http"
	"sync"
	"toolkits/internal/utils"

	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips   map[string]*rate.Limiter
	mutex *sync.RWMutex
	r     rate.Limit //rate
	b     int        //burst
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips:   make(map[string]*rate.Limiter),
		mutex: &sync.RWMutex{},
		r:     r,
		b:     b,
	}
}

// method
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mutex.RLock()
	limiter, exists := i.ips[ip]
	i.mutex.RUnlock()

	if !exists {
		i.mutex.Lock()
		defer i.mutex.Unlock()

		if limiter, exists := i.ips[ip]; exists {
			return limiter
		}

		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
		return limiter
	}
	return limiter
}

func RateLimitMiddleware(rps float64, burst int) func(http.Handler) http.Handler {
	limiter := NewIPRateLimiter(rate.Limit(rps), burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			limiterForIP := limiter.GetLimiter(ip)
			if !limiterForIP.Allow() {
				utils.JSONResponse(w, http.StatusTooManyRequests, map[string]any{
					"status":  http.StatusTooManyRequests,
					"message": "Terlalu banyak permintaan, silahkan tunggu sebentar",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
