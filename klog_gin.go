package klog

import (
	"time"

	"github.com/gin-gonic/gin"
)

func setupLogging(duration time.Duration) {
	go func() {
		for range time.Tick(duration) {
			Flush()
		}
	}()
}

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// ErrorLogger returns an ErrorLoggerT with parameter gin.ErrorTypeAny
func GinErrorLogger() gin.HandlerFunc {
	return GinErrorLoggerT(gin.ErrorTypeAny)
}

// GinErrorLoggerT returns an GinErrorLoggerT middleware with the given
// type gin.ErrorType.
func GinErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !c.Writer.Written() {
			json := c.Errors.ByType(typ).JSON()
			if json != nil {
				c.JSON(-1, json)
			}
		}
	}
}

// Logger prints a logline for each request and measures the time to
// process for a call. It formats the log entries similar to
// http://godoc.org/github.com/gin-gonic/gin#Logger does.
//
// Example:
//        router := gin.New()
//        router.Use(gin.Recovery())
//        router.Use(klog.Logger(3 * time.Second))
func GinLogger(duration time.Duration) gin.HandlerFunc {
	setupLogging(duration)
	return func(c *gin.Context) {
		t := time.Now()

		// process request
		c.Next()

		latency := time.Since(t)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)
		path := c.Request.URL.Path

		switch {
		case statusCode >= 400 && statusCode <= 499:
			{
				Warningf("[GIN] |%s %3d %s| %12v | %s |%s  %s %-7s %s\n%s",
					statusColor, statusCode, reset,
					latency,
					clientIP,
					methodColor, method, reset,
					path,
					c.Errors.String(),
				)
			}
		case statusCode >= 500:
			{
				Errorf("[GIN] |%s %3d %s| %12v | %s |%s  %s %-7s %s\n%s",
					statusColor, statusCode, reset,
					latency,
					clientIP,
					methodColor, method, reset,
					path,
					c.Errors.String(),
				)
			}
		default:
			V(0).Infof("[GIN] |%s %3d %s| %12v | %s |%s  %s %-7s %s\n%s",
				statusColor, statusCode, reset,
				latency,
				clientIP,
				methodColor, method, reset,
				path,
				c.Errors.String(),
			)
		}

	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return reset
	}
}
