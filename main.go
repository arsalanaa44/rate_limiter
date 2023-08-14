package main

import (
	"fmt"
	"github.com/arsalanaa44/rate_limiter/internal/config"
	"github.com/arsalanaa44/rate_limiter/internal/db/redis"
	"github.com/arsalanaa44/rate_limiter/internal/handler"
	"github.com/arsalanaa44/rate_limiter/pkg/strategy"
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
		if err != nil {
			panic(err)
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	}

	db, err := redis.NewRedisClient(cfg.Database)
	if err != nil {
		logger.Fatal("unable to connect to redis", zap.Error(err))
	}

	hs := handler.SignUp{
		RedisClient: db,
		Logger:      logger.Named("signup"),
	}
	hs.Register(app.Group(""))

	hd := handler.Cache{
		RedisClient: db,
		Logger:      logger.Named("cache"),
	}
	hm := handler.MonthlyQuotaChecker{
		RedisClient: db,
		Logger:      logger.Named("monthly limit"),
	}

	hr := handler.RateLimiter{
		RedisClient: db,
		Logger:      logger.Named("rate limiter"),
		Strategy:    strategy.NewSortedSetCounterStrategy(db, time.Now),
	}

	app.GET("/hello", handler.Hello,
		hr.RateLimit,
		hm.LimitConsumption,
		hd.CheckDataCache,
	)

	app.Debug = cfg.Debug

	if err := app.Start(":" + fmt.Sprint(cfg.Port)); err != nil {
		logger.Error("cannot start the http server", zap.Error(err))
	}
}
