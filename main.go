package main

import (
	"fmt"
	"github.com/arsalanaa44/rate_limiter/internal/config"
	"github.com/arsalanaa44/rate_limiter/internal/db/redis"
	"github.com/arsalanaa44/rate_limiter/internal/handler"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// just check the method
func main() {

	app := echo.New()

	cfg := config.New()
	var (
		logger *zap.Logger
		err    error
	)

	if cfg.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	db, err := redis.NewRedisClient(cfg.Database)
	if err != nil {
		logger.Fatal("unable to connect to redis", zap.Error(err))
	}

	ha := handler.SignUp{db, logger.Named("signup")}

	ha.Register(app.Group(""))

	app.Debug = cfg.Debug

	if err := app.Start(":" + fmt.Sprint(cfg.Port)); err != nil {
		logger.Error("cannot start the http server", zap.Error(err))
	}
}