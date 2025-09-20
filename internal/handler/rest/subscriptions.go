package rest

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/mirrorblade/subscriptions/internal/domain"
	"github.com/mirrorblade/subscriptions/internal/repository"
)

type subscriptionJSON struct {
	ServiceName string `json:"service_name"`
	Price       int64  `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

func (h *Handler) initSubscriptions(g *echo.Group) {
	group := g.Group("/subscriptions")
	group.GET("/:id", h.getSubscription)
	group.GET("/", h.getSubscriptions)
	group.GET("/price", h.getSubscriptionsSum)
	group.POST("/", h.createSubscription)
	group.PATCH("/:id", h.updateSubscription)
	group.DELETE("/:id", h.deleteSubscription)
}

func (h *Handler) getSubscription(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad request",
		})
	}

	subscription, err := h.service.Subscriptions.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrSubscriptionNotFound) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad request",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, subscription)
}

func (h *Handler) getSubscriptions(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad request",
		})
	}

	subscriptions, err := h.service.Subscriptions.GetListByUserID(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad request",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, subscriptions)
}

func (h *Handler) getSubscriptionsSum(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad request",
		})
	}

	var serviceName *string
	dirtyServiceName := c.Param("service_name")
	if dirtyServiceName != "" {
		sanitizedServiceName := h.sanitizer.Sanitize(dirtyServiceName)
		serviceName = &sanitizedServiceName
	}

	var fromDate *time.Time
	dirtyFromDate := c.Param("from_date")
	if dirtyFromDate != "" {
		date, err := time.Parse("01-2006", dirtyFromDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad request",
			})
		}
		fromDate = &date
	}

	var toDate *time.Time
	dirtyToDate := c.Param("to_date")
	if dirtyToDate != "" {
		date, err := time.Parse("01-2006", dirtyToDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad request",
			})
		}
		toDate = &date
	}

	parameters := repository.GetSumParameters{
		ServiceName: serviceName,
		FromDate:    fromDate,
		ToDate:      toDate,
	}

	price, err := h.service.Subscriptions.GetPriceSumByUserID(c.Request().Context(), userID, parameters)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad request",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, price)
}

func (h *Handler) createSubscription(c echo.Context) error {
	body := new(subscriptionJSON)
	if err := c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad request",
		})
	}

	userID, err := uuid.Parse(body.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad request",
		})
	}

	startDate, err := time.Parse("01-2006", body.StartDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad request",
		})
	}

	var endDate *time.Time

	if body.EndDate != "" {
		date, err := time.Parse("01-2006", body.EndDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad request",
			})
		}

		endDate = &date
	}

	subscription := domain.Subscription{
		ServiceName: h.sanitizer.Sanitize(body.ServiceName),
		Price:       body.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := h.service.Subscriptions.Create(c.Request().Context(), subscription); err != nil {
		if errors.Is(err, domain.ErrInvalidPrice) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad request",
			})
		}

		if errors.Is(err, domain.ErrInvalidDate) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad request",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "internal server error",
		})
	}
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "subscription was succesfully created",
	})
}

func (h *Handler) updateSubscription(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad request",
		})
	}

	var price *int64
	dirtyPrice := c.Param("price")
	if dirtyPrice != "" {
		clearPrice, err := strconv.ParseInt(dirtyPrice, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad request",
			})
		}
		price = &clearPrice
	}

	var endDate *time.Time
	dirtyEndDate := c.Param("end_date")
	if dirtyEndDate != "" {
		date, err := time.Parse("01-2006", dirtyEndDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad request",
			})
		}
		endDate = &date
	}

	parameters := repository.UpdateParameters{
		Price:   price,
		EndDate: endDate,
	}

	if err := h.service.Subscriptions.UpdateByID(c.Request().Context(), id, parameters); err != nil {
		if errors.Is(err, domain.ErrSubscriptionNotFound) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad request",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "internal server error",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) deleteSubscription(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad request",
		})
	}

	if err := h.service.Subscriptions.DeleteByID(c.Request().Context(), id); err != nil {
		if errors.Is(err, domain.ErrSubscriptionNotFound) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad request",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "internal server error",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
