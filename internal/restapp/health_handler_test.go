package restapp_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	api "github.com/go-seidon/hippo/api/restapp"
	"github.com/go-seidon/hippo/internal/healthcheck"
	mock_healthcheck "github.com/go-seidon/hippo/internal/healthcheck/mock"
	"github.com/go-seidon/hippo/internal/restapp"
	mock_restapp "github.com/go-seidon/hippo/internal/restapp/mock"
	mock_logging "github.com/go-seidon/provider/logging/mock"
	mock_serialization "github.com/go-seidon/provider/serialization/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Health Handler", func() {

	Context("CheckHealth Handler", Label("unit"), func() {
		var (
			ctx           context.Context
			handler       http.HandlerFunc
			r             *http.Request
			w             *mock_restapp.MockResponseWriter
			log           *mock_logging.MockLogger
			serializer    *mock_serialization.MockSerializer
			healthService *mock_healthcheck.MockHealthCheck
		)

		BeforeEach(func() {
			ctx = context.Background()
			t := GinkgoT()
			r = &http.Request{}
			ctrl := gomock.NewController(t)
			w = mock_restapp.NewMockResponseWriter(ctrl)
			log = mock_logging.NewMockLogger(ctrl)
			serializer = mock_serialization.NewMockSerializer(ctrl)
			healthService = mock_healthcheck.NewMockHealthCheck(ctrl)
			healthHandler := restapp.NewHealthHandler(restapp.HealthHandlerParam{
				Logger:        log,
				Serializer:    serializer,
				HealthService: healthService,
			})
			handler = healthHandler.CheckHealth
		})

		When("failed check service health", func() {
			It("should write response", func() {

				err := fmt.Errorf("failed check health")

				b := restapp.ResponseBody{
					Code:    1001,
					Message: err.Error(),
				}

				healthService.
					EXPECT().
					Check(gomock.Eq(ctx)).
					Return(nil, err).
					Times(1)

				serializer.
					EXPECT().
					Marshal(gomock.Eq(b)).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(500)).
					Times(1)

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})

		When("success check service health", func() {
			It("should write response", func() {

				currentTimestamp := time.Now()
				res := &healthcheck.CheckResult{
					Status: "WARNING",
					Items: map[string]healthcheck.CheckResultItem{
						"app-disk": {
							Name:      "app-disk",
							Status:    "FAILED",
							Error:     "Critical: disk usage too high 96.71 percent",
							CheckedAt: currentTimestamp,
						},
						"internet-connection": {
							Name:      "internet-connection",
							Status:    "OK",
							Error:     "",
							CheckedAt: currentTimestamp,
						},
					},
				}

				healthService.
					EXPECT().
					Check(gomock.Eq(ctx)).
					Return(res, nil).
					Times(1)

				details := map[string]api.CheckHealthDetail{
					"app-disk": {
						Name:      "app-disk",
						Status:    "FAILED",
						Error:     "Critical: disk usage too high 96.71 percent",
						CheckedAt: currentTimestamp.UnixMilli(),
					},
					"internet-connection": {
						Name:      "internet-connection",
						Status:    "OK",
						Error:     "",
						CheckedAt: currentTimestamp.UnixMilli(),
					},
				}
				b := restapp.ResponseBody{
					Code:    1000,
					Message: "success check service health",
					Data: &api.CheckHealthData{
						Status: "WARNING",
						Details: api.CheckHealthData_Details{
							AdditionalProperties: details,
						},
					},
				}
				serializer.
					EXPECT().
					Marshal(gomock.Eq(b)).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(200))

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})
	})

})
