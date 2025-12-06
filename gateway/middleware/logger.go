package middleware

import (
	"bytes"
	"go-api-gateway/utils"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceID := uuid.New().String()
		c.Set("traceID", traceID)

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		w := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		latency := time.Since(start).Milliseconds()

		entry := utils.LogEntry{
			Timestamp:    start.UTC(),
			TraceID:      traceID,
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			StatusCode:   c.Writer.Status(),
			ClientIP:     c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			RequestBody:  string(requestBody),
			ResponseBody: w.body.String(),
			LatencyMs:    latency,
		}

		if len(entry.RequestBody) > 1024 {
			entry.RequestBody = entry.RequestBody[:1024] + "...(truncated)"
		}
		if len(entry.ResponseBody) > 1024 {
			entry.ResponseBody = entry.ResponseBody[:1024] + "...(truncated)"
		}

		utils.SendLogToES(entry)
	}
}
