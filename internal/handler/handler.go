package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mirrorblade/subscriptions/internal/handler/rest"
	"github.com/mirrorblade/subscriptions/internal/service"
)

type Handler struct {
	router  *echo.Echo
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Init() {
	h.router = echo.New()

	h.checkHealth()

	h.initRest()
}

func (h *Handler) initRest() {
	handler := rest.New(h.service, bluemonday.UGCPolicy())
	group := h.router.Group("/rest")
	handler.Init(group)
}

func (h *Handler) checkHealth() {
	h.router.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "ok",
		})
	})
}
