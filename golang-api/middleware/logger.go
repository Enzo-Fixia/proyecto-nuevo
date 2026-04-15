package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		go func(status int, latency time.Duration) {
			log.Printf("[%s] %s %d (%v)", method, path, status, latency)
		}(c.Writer.Status(), time.Since(start))
	}
}
