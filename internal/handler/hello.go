package handler

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Hello(c echo.Context) error {

	doneCh := make(chan error)

	go func() {
		//	time.Sleep(20 * time.Second) busy work (queue e.g)

		doneCh <- nil
	}()

	select {
	case err := <-doneCh:
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, "HI")
	case <-c.Request().Context().Done():
		if errors.Is(c.Request().Context().Err(), context.DeadlineExceeded) {
			return c.String(http.StatusServiceUnavailable, "timeout")
		}
		return c.Request().Context().Err()
	}
}
