package restapp

import (
	"net/http"

	"github.com/go-seidon/hippo/api/restapp"
	"github.com/go-seidon/hippo/internal/healthcheck"
	"github.com/go-seidon/provider/logging"
	"github.com/go-seidon/provider/serialization"
	"github.com/go-seidon/provider/status"
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

	details := restapp.CheckHealthData_Details{
		AdditionalProperties: map[string]restapp.CheckHealthDetail{},
	}
	for checkName, item := range r.Items {
		details.AdditionalProperties[checkName] = restapp.CheckHealthDetail{
			Name:      item.Name,
			Status:    item.Status,
			Error:     item.Error,
			CheckedAt: item.CheckedAt.UnixMilli(),
		}
	}
	d := &restapp.CheckHealthData{
		Status: r.Status,

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
