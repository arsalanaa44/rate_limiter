package handler

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"strconv"
)

type User struct {
	ID             string
	MonthSizeLimit int
	SizeConsumed   int
}

type MonthlyQuotaChecker struct {
	RedisClient *redis.Client
	Logger      *zap.Logger
}

func (m MonthlyQuotaChecker) LimitConsumption(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		userID := c.Request().Header.Get(userID)
		dataSize, _ := strconv.Atoi(c.Request().Header.Get(dataSize))

		userData, err := m.RedisClient.HGetAll(ctx, "users:"+userID).Result()
		if err != nil {
			m.Logger.Error("Failed to get user record from Redis", zap.Error(err))
			return echo.ErrForbidden
		}
		// TODO handle errors
		var (
			consumed int
			limit    int
		)

		if val, ok := userData[monthlyLimitKey]; ok {
			limit, _ = strconv.Atoi(val)
		}

		if val, ok := userData[consumptionKey]; ok {
			consumed, _ = strconv.Atoi(val)
		}

		consumed += dataSize

		if consumed > limit {
			m.Logger.Error("size limitation reached", zap.Any("consumed", consumed))
			return echo.NewHTTPError(405, "you are using more than your quota")
		}

		// Update the user record in Redis
		err = m.RedisClient.HSet(ctx, "users:"+userID,
			consumptionKey, consumed).Err()
		if err != nil {
			m.Logger.Error("Failed to update user record in Redis", zap.Error(err))
			return echo.ErrInternalServerError
		}
		return next(c)
	}
}
