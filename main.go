package main

import (
	"fmt"
	"github.com/arsalanaa44/rate_limiter/internal/config"
	"github.com/arsalanaa44/rate_limiter/internal/db/redis"
	"github.com/arsalanaa44/rate_limiter/internal/handler"
	"github.com/arsalanaa44/rate_limiter/pkg/redis_rate_limiter"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
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

	hs := handler.SignUp{RedisClient: db, Logger: logger.Named("signup")}
	app.POST("/signup", hs.RegisterUser)

	hd := handler.Cache{RedisClient: db, Logger: logger.Named("data checker")}
	hm := handler.MonthlyQuotaChecker{RedisClient: db, Logger: logger.Named("monthly quota checker")}

	hr := handler.RateLimiter{
		RedisClient: db,
		Logger:      logger.Named("rate limiter"),
		Strategy:    redis_rate_limiter.NewSortedSetCounterStrategy(db, time.Now),
	}

	grp := app.Group("")
	grp.Use(hr.RateLimit)
	grp.Use(hm.Checker)
	grp.Use(hd.IsDataCached)
	grp.GET("/hello", handler.Hello)

	app.Debug = cfg.Debug

	if err := app.Start(":" + fmt.Sprint(cfg.Port)); err != nil {
		logger.Error("cannot start the http server", zap.Error(err))
	}
}
