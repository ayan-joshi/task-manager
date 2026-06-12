package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"taskmanager/internal/utils"
)

// visitor tracks a per-client rate limiter and the last time it was seen, so
// idle entries can be reclaimed.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter is an in-memory, per-IP token-bucket limiter. For a single-node
// deployment this is sufficient; a multi-node deployment would back this with
// Redis. It self-cleans idle visitors to bound memory.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rps      rate.Limit
	burst    int
}

// NewRateLimiter constructs a RateLimiter and starts its cleanup loop.
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rps:      rate.Limit(rps),
		burst:    burst,
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, ok := rl.visitors[ip]
	if !ok {
		lim := rate.NewLimiter(rl.rps, rl.burst)
		rl.visitors[ip] = &visitor{limiter: lim, lastSeen: time.Now()}
		return lim
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Middleware returns the gin handler enforcing the limit.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.getLimiter(c.ClientIP()).Allow() {
			c.Header("Retry-After", "1")
			c.JSON(http.StatusTooManyRequests, utils.Envelope{
				Success: false,
				Error: &utils.ErrorBody{
					Code:    "RATE_LIMITED",
					Message: "too many requests, please slow down",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
