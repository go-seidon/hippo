package resthandler

import (
	"net/http"

	"github.com/go-seidon/hippo/api/restapp"
	"github.com/go-seidon/hippo/internal/healthcheck"
	"github.com/go-seidon/provider/status"
	"github.com/labstack/echo/v4"
)

type healthHandler struct {
	healthClient healthcheck.HealthCheck
}

func (h *healthHandler) CheckHealth(ctx echo.Context) error {
	health, err := h.healthClient.Check(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, &restapp.ResponseBodyInfo{
			Code:    status.ACTION_FAILED,
			Message: err.Error(),
		})
	}

	details := restapp.CheckHealthData_Details{
		AdditionalProperties: map[string]restapp.CheckHealthDetail{},
	}
	for _, item := range health.Items {
		details.AdditionalProperties[item.Name] = restapp.CheckHealthDetail{
			Name:      item.Name,
			Status:    item.Status,
			Error:     item.Error,
			CheckedAt: item.CheckedAt.UnixMilli(),
		}
	}

	return ctx.JSON(http.StatusOK, &restapp.CheckHealthResponse{
		Code:    status.ACTION_SUCCESS,
		Message: "success check health",
		Data: restapp.CheckHealthData{
			Details: details,
			Status:  health.Status,
		},
	})
}

type HealthParam struct {
	HealthClient healthcheck.HealthCheck
}

func NewHealth(p HealthParam) *healthHandler {
	return &healthHandler{
		healthClient: p.HealthClient,
	}
}
