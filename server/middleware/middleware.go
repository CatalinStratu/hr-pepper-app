package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger logs request information
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[%s] %s %d %v", c.Request.Method, path, status, latency)
	}
}

// ErrorHandler handles panics and returns proper error responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.AbortWithStatusJSON(500, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}
