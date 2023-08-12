package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

const SetKey = "dataID_set:"

type DataChecker struct {
	RedisClient *redis.Client
	Logger      *zap.Logger
}

func (ch DataChecker) Checker(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		dataID := c.Request().Header.Get("DataID")
		userID := c.Request().Header.Get("UserID")

		dataIDHash := sha256.Sum256([]byte(dataID))
		dataIDHashString := hex.EncodeToString(dataIDHash[:])
		member := dataIDHashString
		set := SetKey + userID

		isMember := ch.RedisClient.SIsMember(ctx, set, member).Val()
		if isMember {
			ch.Logger.Info("data already exist")
			return c.String(http.StatusOK, "dataID already exists in Redis set")

		} else {
			err := ch.RedisClient.SAdd(ctx, set, member).Err()
			if err != nil {
				ch.Logger.Error("Failed to add dataID to Redis", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process request")
			}
		}

		return next(c)
	}
}

func (ch DataChecker) Register(g *echo.Group) {
}
