package rest_app

import (
	"net/http"

	rest_v1 "github.com/go-seidon/hippo/generated/rest-v1"
	"github.com/go-seidon/hippo/internal/healthcheck"
	"github.com/go-seidon/hippo/internal/logging"
	"github.com/go-seidon/hippo/internal/serialization"
	"github.com/go-seidon/hippo/internal/status"
)

type healthHandler struct {
	logger        logging.Logger
	serializer    serialization.Serializer
	healthService healthcheck.HealthCheck
}

func (h *healthHandler) CheckHealth(w http.ResponseWriter, req *http.Request) {
	r, err := h.healthService.Check(req.Context())
	if err != nil {

		Response(
			WithWriterSerializer(w, h.serializer),
			WithCode(status.ACTION_FAILED),
			WithMessage(err.Error()),
			WithHttpCode(http.StatusInternalServerError),
		)
		return
	}

	details := map[string]rest_v1.CheckHealthDetail{}
	for checkName, item := range r.Items {
		details[checkName] = rest_v1.CheckHealthDetail{
			Name:      item.Name,
			Status:    item.Status,
			Error:     item.Error,
			CheckedAt: item.CheckedAt.UnixMilli(),
		}
	}
	d := &rest_v1.CheckHealthData{
		Status:  r.Status,
		Details: details,
	}

	Response(
		WithWriterSerializer(w, h.serializer),
		WithData(d),
		WithMessage("success check service health"),
	)
}

type HealthHandlerParam struct {
	Logger        logging.Logger
	Serializer    serialization.Serializer
	HealthService healthcheck.HealthCheck
}

func NewHealthHandler(p HealthHandlerParam) *healthHandler {
	return &healthHandler{
		logger:        p.Logger,
		serializer:    p.Serializer,
		healthService: p.HealthService,
	}
}
