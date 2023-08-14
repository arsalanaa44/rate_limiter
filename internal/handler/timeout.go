package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"time"
)

func TimeoutMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Create a new context with a timeout
		ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
		defer cancel()

		// Replace the request with a new one using the context with a timeout
		c.SetRequest(c.Request().WithContext(ctx))

		// Call the next handler in the chain
		return next(c)
	}
}
