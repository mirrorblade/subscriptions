package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mirrorblade/subscriptions/internal/config"
	"github.com/mirrorblade/subscriptions/internal/handler"
	"github.com/mirrorblade/subscriptions/internal/repository"
	"github.com/mirrorblade/subscriptions/internal/repository/postgresql"
	"github.com/mirrorblade/subscriptions/internal/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("/var/log/backend/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	var logger *zap.Logger

	if config.App.Production {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

		logger = zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(file),
			zapcore.InfoLevel,
		), zap.AddCaller())
	} else {
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

		logger = zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(file),
			zapcore.DebugLevel,
		), zap.AddCaller())
	}

	defer logger.Sync()

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := pool.Ping(ctx); err != nil {
		cancel()
		panic(err)
	}
	cancel()

	subscriptionsRepository := postgresql.NewSubscriptions(pool, "subscriptions")
	repository := repository.New(subscriptionsRepository)

	subscriptionsService := service.NewSubscriptionsService(repository.Subscriptions)
	service := service.New(subscriptionsService)

	handler := handler.New(service, logger, &config.Server)
	handler.Init()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		if err := handler.Start(); err != nil && err != http.ErrServerClosed {
			stop()
			panic("shutting down the server")
		}
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := handler.Shutdown(ctx); err != nil {
		panic(err)
	}
}
