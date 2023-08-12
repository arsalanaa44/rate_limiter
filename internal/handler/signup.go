package handler

import (
	"github.com/arsalanaa44/rate_limiter/internal/model"
	"github.com/arsalanaa44/rate_limiter/internal/request"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type SignUp struct {
	RedisClient *redis.Client
	Logger      *zap.Logger
}

func (s SignUp) RegisterUser(c echo.Context) error {

	ctx := c.Request().Context()

	var req request.SignUp

	if err := c.Bind(&req); err != nil {
		s.Logger.Error("cannot bind request to user quota",
			zap.Error(err),
		)

		return echo.ErrBadRequest
	}

	if err := req.Validate(); err != nil {
		s.Logger.Error("request validation failed",
			zap.Error(err),
			zap.Any("request", req),
		)

		return echo.ErrBadRequest
	}

	m := model.User{
		ID:    uuid.New(),
		Quota: req.Quota,
	}

	if err := s.RedisClient.HSet(ctx, "users:"+m.ID.String(),
		"MonthSizeLimit", m.MonthSizeLimit,
		"SizeConsumed", 0,
		"MinuteRateLimit", m.MinuteRateLimit).
		Err(); err != nil {
		s.Logger.Error("hashset addition failed", zap.Error(err))
	}

	s.Logger.Info("user creation success")

	return c.JSON(http.StatusCreated, m)

}

func (s SignUp) Register(g *echo.Group) {
	g.POST("/signup", s.RegisterUser)
}
