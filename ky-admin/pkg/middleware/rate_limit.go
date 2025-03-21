package middleware

import (
	"sync"
	"time"

	"github.com/SupenBysz/ky-admin/pkg/api"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter IP速率限制器
type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     *sync.RWMutex
	r      rate.Limit
	b      int
	expiry time.Duration
	last   map[string]time.Time
}

// NewIPRateLimiter 创建一个新的IP速率限制器
// r 表示每秒允许的请求数（令牌产生速率）
// b 表示令牌桶的容量
// expiry 表示IP限制器的过期时间
func NewIPRateLimiter(r rate.Limit, b int, expiry time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		mu:     &sync.RWMutex{},
		r:      r,
		b:      b,
		expiry: expiry,
		last:   make(map[string]time.Time),
	}
}

// GetLimiter 获取IP对应的限制器，如果不存在则创建
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	last, _ := i.last[ip]
	i.mu.RUnlock()

	now := time.Now()
	if !exists || now.Sub(last) > i.expiry {
		i.mu.Lock()
		// 双重检查
		if limiter, exists = i.ips[ip]; !exists || now.Sub(i.last[ip]) > i.expiry {
			limiter = rate.NewLimiter(i.r, i.b)
			i.ips[ip] = limiter
			i.last[ip] = now
		}
		i.mu.Unlock()
	} else {
		// 更新最后访问时间
		i.mu.Lock()
		i.last[ip] = now
		i.mu.Unlock()
	}

	return limiter
}

// RateLimit 速率限制中间件
func RateLimit(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		ip := c.ClientIP()
		// 获取对应的限制器
		ipLimiter := limiter.GetLimiter(ip)

		// 尝试获取令牌
		if !ipLimiter.Allow() {
			api.Fail(c, 429, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// DefaultRateLimit 默认速率限制中间件（每秒5个请求，最多10个令牌，30分钟过期）
func DefaultRateLimit() gin.HandlerFunc {
	limiter := NewIPRateLimiter(5, 10, 30*time.Minute)
	return RateLimit(limiter)
}
