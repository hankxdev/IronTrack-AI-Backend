package router

import (
	"bytes"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// DevelopmentLogger logs detailed request and response information for errors in debug mode
func DevelopmentLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only run in debug mode
		if gin.Mode() != gin.DebugMode {
			c.Next()
			return
		}

		// Read and restore request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap response writer to capture response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		// For errors (status code >= 400), log extra information
		status := c.Writer.Status()
		if status >= 400 {
			log.Printf("\n--- [DEVELOPMENT ERROR LOG] ---")
			log.Printf("Path:   %s %s", c.Request.Method, c.Request.URL.Path)
			log.Printf("Status: %d", status)

			if len(requestBody) > 0 {
				log.Printf("Request Body: %s", string(requestBody))
			}

			log.Printf("Response Body: %s", blw.body.String())

			if len(c.Errors) > 0 {
				log.Printf("Gin Errors: %v", c.Errors.String())
			}
			log.Printf("-------------------------------\n")
		}
	}
}
