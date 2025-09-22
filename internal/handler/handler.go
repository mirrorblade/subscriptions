package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mirrorblade/subscriptions/internal/config"
	"github.com/mirrorblade/subscriptions/internal/handler/rest"
	"github.com/mirrorblade/subscriptions/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	router  *echo.Echo
	service *service.Service

	logger *zap.Logger

	config *config.Server
}

func New(service *service.Service, logger *zap.Logger, config *config.Server) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
		config:  config,
	}
}

func (h *Handler) Init() {
	h.router = echo.New()

	h.router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogMethod:   true,
		LogStatus:   true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fields := []zap.Field{
				zap.String("uri", v.URI),
				zap.String("method", v.Method),
				zap.Int("status", v.Status),
				zap.String("ip", v.RemoteIP),
			}

			if v.Status < 400 {
				h.logger.Info("request", fields...)
			} else if v.Status < 500 {
				errorMessage, ok := c.Get("error").(error)
				if ok {
					fields = append(fields, zap.Error(errorMessage))
				}

				h.logger.Warn("request", fields...)

			} else {
				errorMessage, ok := c.Get("error").(error)
				if ok {
					fields = append(fields, zap.Error(errorMessage))
				}

				h.logger.Error("request", fields...)
			}

			return nil
		},
	}))

	corsConfig := middleware.DefaultCORSConfig

	corsConfig.AllowOrigins = h.config.CORS.AllowOrigins
	corsConfig.AllowMethods = h.config.CORS.AllowMethods
	corsConfig.AllowHeaders = h.config.CORS.AllowHeaders
	corsConfig.AllowCredentials = h.config.CORS.AllowCredentials
	corsConfig.MaxAge = int(h.config.CORS.MaxAge.Seconds())

	h.router.Use(middleware.CORSWithConfig(corsConfig))

	h.router.Use(middleware.AddTrailingSlash())

	h.checkHealth()

	h.initRest()
}

func (h *Handler) Start() error {
	return h.router.Start(h.config.Host + ":" + h.config.Port)
}

func (h *Handler) Shutdown(context context.Context) error {
	return h.router.Shutdown(context)
}

func (h *Handler) initRest() {
	group := h.router.Group("/rest")

	handler := rest.New(h.service, bluemonday.UGCPolicy())
	handler.Init(group)
}

func (h *Handler) checkHealth() {
	h.router.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "ok",
		})
	})
}
