package middlewares

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
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

func RateLimitMiddleware(rps float64, burst int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiterForIP := limiter.GetLimiter(ip)

		if !limiterForIP.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  http.StatusTooManyRequests,
				"message": "Terlalu banyak permintaan, silahkan tunggu sebentar",
			})
			c.Abort()
			return
		}
		c.Next()

	}
}
