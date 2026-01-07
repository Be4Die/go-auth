package httpv1

import (
	"time"

	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = "req_" + time.Now().UTC().Format("20060102T150405.000000000")
		}
		c.Set("request_id", rid)
		c.Writer.Header().Set("X-Request-ID", rid)
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
		if c.Request.Method == "OPTIONS" {
			c.Status(204)
			return
		}
		c.Next()
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none';")
		c.Next()
	}
}

func RateLimit(max int, window time.Duration) gin.HandlerFunc {
	type bucket struct {
		count int
		until time.Time
	}
	store := make(map[string]*bucket)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		b := store[ip]
		now := time.Now()
		if b == nil || now.After(b.until) {
			b = &bucket{count: 0, until: now.Add(window)}
			store[ip] = b
		}
		if b.count >= max {
			c.JSON(429, gin.H{"error": "Too many requests", "code": "RATE_LIMITED"})
			c.Abort()
			return
		}
		b.count++
		c.Next()
	}
}
