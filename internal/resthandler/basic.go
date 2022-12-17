package resthandler

import (
	"net/http"

	"github.com/go-seidon/hippo/api/restapp"
	"github.com/go-seidon/provider/status"
	"github.com/labstack/echo/v4"
)

type basicHandler struct {
	config *BasicConfig
}

func (h *basicHandler) GetAppInfo(ctx echo.Context) error {
	res := &restapp.GetAppInfoResponse{
		Code:    status.ACTION_SUCCESS,
		Message: "success get app info",
		Data: restapp.GetAppInfoData{
			AppName:    h.config.AppName,
			AppVersion: h.config.AppVersion,
		},
	}
	return ctx.JSON(http.StatusOK, res)
}

type BasicConfig struct {
	AppName    string
	AppVersion string
}

type BasicParam struct {
	Config *BasicConfig
}

func NewBasic(p BasicParam) *basicHandler {
	return &basicHandler{
		config: p.Config,
	}
}
