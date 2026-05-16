package middlewares

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func RequestTracer() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			req := c.Request()

			traceID := req.Header.Get("X-Trace-ID")
			if traceID == "" {
				traceID = uuid.New().String()
			}

			req.Header.Set("X-Trace-ID", traceID)

			return next(c)
		}
	}
}
