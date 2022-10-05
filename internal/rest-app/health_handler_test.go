package rest_app_test

import (
	"fmt"
	"net/http"
	"time"

	rest_v1 "github.com/go-seidon/local/generated/rest-v1"
	"github.com/go-seidon/local/internal/healthcheck"
	mock_healthcheck "github.com/go-seidon/local/internal/healthcheck/mock"
	mock_logging "github.com/go-seidon/local/internal/logging/mock"
	rest_app "github.com/go-seidon/local/internal/rest-app"
	mock_restapp "github.com/go-seidon/local/internal/rest-app/mock"
	mock_serialization "github.com/go-seidon/local/internal/serialization/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Health Handler", func() {

	Context("CheckHealth Handler", Label("unit"), func() {
		var (
			handler       http.HandlerFunc
			r             *http.Request
			w             *mock_restapp.MockResponseWriter
			log           *mock_logging.MockLogger
			serializer    *mock_serialization.MockSerializer
			healthService *mock_healthcheck.MockHealthCheck
		)

		BeforeEach(func() {
			t := GinkgoT()
			r = &http.Request{}
			ctrl := gomock.NewController(t)
			w = mock_restapp.NewMockResponseWriter(ctrl)
			log = mock_logging.NewMockLogger(ctrl)
			serializer = mock_serialization.NewMockSerializer(ctrl)
			healthService = mock_healthcheck.NewMockHealthCheck(ctrl)
			healthHandler := rest_app.NewHealthHandler(rest_app.HealthHandlerParam{
				Logger:        log,
				Serializer:    serializer,
				HealthService: healthService,
			})
			handler = healthHandler.CheckHealth
		})

		When("failed check service health", func() {
			It("should write response", func() {

				err := fmt.Errorf("failed check health")

				b := rest_app.ResponseBody{
					Code:    1001,
					Message: err.Error(),
				}

				healthService.
					EXPECT().
					Check().
					Return(nil, err).
					Times(1)

				serializer.
					EXPECT().
					Marshal(b).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(400)).
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

				details := map[string]rest_v1.CheckHealthDetail{
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
				b := rest_app.ResponseBody{
					Code:    1000,
					Message: "success check service health",
					Data: &rest_v1.CheckHealthData{
						Status:  "WARNING",
						Details: details,
					},
				}

				healthService.
					EXPECT().
					Check().
					Return(res, nil).
					Times(1)

				serializer.
					EXPECT().
					Marshal(b).
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
