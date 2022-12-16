package rest_app

import (
	"net/http"

	"github.com/go-seidon/hippo/api/restapp"
	"github.com/go-seidon/provider/logging"
	"github.com/go-seidon/provider/serialization"
	"github.com/go-seidon/provider/status"
)

type basicHandler struct {
	logger     logging.Logger
	serializer serialization.Serializer
	config     *RestAppConfig
}

func (h *basicHandler) GetAppInfo(w http.ResponseWriter, req *http.Request) {
	d := &restapp.GetAppInfoData{
		AppName:    h.config.AppName,
		AppVersion: h.config.AppVersion,
	}

	Response(
		WithWriterSerializer(w, h.serializer),
		WithData(d),
	)
}

func (h *basicHandler) NotFound(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	Response(
		WithWriterSerializer(w, h.serializer),
		WithHttpCode(http.StatusNotFound),
		WithCode(status.RESOURCE_NOTFOUND),
		WithMessage("resource not found"),
	)
}

func (h *basicHandler) MethodNotAllowed(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	Response(
		WithWriterSerializer(w, h.serializer),
		WithHttpCode(http.StatusMethodNotAllowed),
		WithCode(status.ACTION_FAILED),
		WithMessage("method is not allowed"),
	)
}

type BasicHandlerParam struct {
	Logger     logging.Logger
	Serializer serialization.Serializer
	Config     *RestAppConfig
}

func NewBasicHandler(p BasicHandlerParam) *basicHandler {
	return &basicHandler{
		logger:     p.Logger,
		serializer: p.Serializer,
		config:     p.Config,
	}
}
