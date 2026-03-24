package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func AuditLogger() gin.HandlerFunc {
	audit := log.New(os.Stdout, "[AUDIT] ", log.LstdFlags)

	return func(c *gin.Context) {
		if !shouldAuditMethod(c.Request.Method) {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		latency := time.Since(start)
		audit.Printf(
			"method=%s path=%s status=%d client_ip=%s latency_ms=%d user_agent=%q",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			c.ClientIP(),
			latency.Milliseconds(),
			c.Request.UserAgent(),
		)
	}
}

func shouldAuditMethod(method string) bool {
	switch method {
	case "POST", "PUT", "PATCH", "DELETE":
		return true
	default:
		return false
	}
}
