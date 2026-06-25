package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = newRequestID()
		}

		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Set("request_id", requestID)

		start := time.Now()
		c.Next()

		path := c.Request.URL.Path
		if fullPath := c.FullPath(); fullPath != "" {
			path = fullPath
		}

		slog.Info("request",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"duration", time.Since(start).String(),
		)
	}
}

func newRequestID() string {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return "unknown"
	}

	return hex.EncodeToString(value[:])
}
