package handler

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

var errRepeatedData = errors.New("data already exist")

const (
	setKey = "data:IDs:"
	dataID = "Data-ID"
	UserID = "User-ID"
)

type Cache struct {
	RedisClient *redis.Client
	Logger      *zap.Logger
}

func (ch Cache) IsDataCached(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		dataID := c.Request().Header.Get(dataID)
		userID := c.Request().Header.Get(UserID)

		set := setKey + userID
		if ch.RedisClient.SIsMember(ctx, set, dataID).Val() {
			ch.Logger.Info(errRepeatedData.Error())
			return c.String(http.StatusOK, errRepeatedData.Error())

		} else {
			err := ch.RedisClient.SAdd(ctx, set, dataID).Err()
			if err != nil {
				ch.Logger.Error("Failed to add dataID to Redis", zap.Error(err))
			}
		}

		return next(c)
	}
}
