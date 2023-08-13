package handler

import (
	"context"
	"fmt"
	"github.com/arsalanaa44/rate_limiter/pkg/strategy"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)


type RateLimiter struct {
	RedisClient *redis.Client
	Logger      *zap.Logger
	Strategy    strategy.Strategy
}

func (l RateLimiter) RateLimit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		//TODO add error handling here 
		userID := c.Request().Header.Get(userID)

		// get the request
		limit := l.findLimit(ctx, userID)
		req := strategy.Request{
			Key:      userID,
			Limit:    limit,
			Duration: time.Minute,
		}

		// run the rate limiter strategy
		result, err := l.Strategy.Run(ctx, &req)
		if err != nil {
			l.Logger.Error("strategy run failed", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		// check the result
		if result.State == strategy.Deny {
			l.Logger.Error("limit reached for ", zap.Any("user", userID))
			return echo.NewHTTPError(http.StatusTooManyRequests)
		}

		// continue with the request
		return next(c)
	}
}

func (l RateLimiter) findLimit(ctx context.Context, userID string) uint64 {

	userData, _ := l.RedisClient.HGetAll(ctx, "users:"+userID).Result()

	var MinuteRate uint64
	if val, ok := userData[minuteRateLimit]; ok {
	_:
		fmt.Sscanf(val, "%d", &MinuteRate)
	}
	return MinuteRate
}
