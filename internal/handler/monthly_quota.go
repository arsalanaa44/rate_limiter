package handler

import (
	"fmt"
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

func (m MonthlyQuotaChecker) Checker(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		userID := c.Request().Header.Get("UserID")
		dataSize, _ := strconv.Atoi(c.Request().Header.Get("cl"))

		userData, err := m.RedisClient.HGetAll(ctx, "users:"+userID).Result()
		if err != nil {
			m.Logger.Error("Failed to get user record from Redis", zap.Error(err))
			return echo.ErrInternalServerError
		}

		user := User{
			ID:             userID,
			MonthSizeLimit: 0,
			SizeConsumed:   0,
		}

		if val, ok := userData["MonthSizeLimit"]; ok {
			fmt.Sscanf(val, "%d", &user.MonthSizeLimit)
		}

		if val, ok := userData["SizeConsumed"]; ok {
			fmt.Sscanf(val, "%d", &user.SizeConsumed)
		}

		user.SizeConsumed += dataSize

		if user.SizeConsumed > user.MonthSizeLimit {
			m.Logger.Error("size limitation reached", zap.Any("consumed", user.SizeConsumed))
			return echo.ErrNotAcceptable
		}

		// Update the user record in Redis
		err = m.RedisClient.HSet(ctx, "users:"+user.ID,
			"SizeConsumed", user.SizeConsumed).Err()
		if err != nil {
			m.Logger.Error("Failed to update user record in Redis", zap.Error(err))
			return echo.ErrInternalServerError
		}

		return next(c)
	}
}

func (m MonthlyQuotaChecker) Register(g *echo.Group) {
}
