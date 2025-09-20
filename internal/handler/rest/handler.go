package rest

import (
	"github.com/labstack/echo"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mirrorblade/subscriptions/internal/service"
)

type Handler struct {
	service *service.Service

	sanitizer *bluemonday.Policy
}

func New(service *service.Service, sanitizer *bluemonday.Policy) *Handler {
	return &Handler{
		service:   service,
		sanitizer: sanitizer,
	}
}

func (h *Handler) Init(group *echo.Group) {
	h.initSubscriptions(group)
}
