package grpchandler_test

import (
	"context"
	"time"

	api "github.com/go-seidon/hippo/api/grpcapp"
	"github.com/go-seidon/hippo/internal/grpchandler"
	"github.com/go-seidon/hippo/internal/healthcheck"
	mock_healthcheck "github.com/go-seidon/hippo/internal/healthcheck/mock"
	"github.com/go-seidon/provider/system"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"

	grpc_status "google.golang.org/grpc/status"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Health Package", func() {

	Context("CheckHealth function", Label("unit"), func() {
		var (
			handler       api.HealthServiceServer
			healthService *mock_healthcheck.MockHealthCheck
			ctx           context.Context
			p             *api.CheckHealthParam
			r             *api.CheckHealthResult
			checkRes      *healthcheck.CheckResult
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			healthService = mock_healthcheck.NewMockHealthCheck(ctrl)
			handler = grpchandler.NewHealth(grpchandler.HealthParam{
				HealthClient: healthService,
			})
			ctx = context.Background()
			currentTs := time.Now().UTC()
			p = &api.CheckHealthParam{}
			r = &api.CheckHealthResult{
				Code:    1000,
				Message: "success check health",
				Data: &api.CheckHealthData{
					Status: "WARNING",
					Details: map[string]*api.CheckHealthDetail{
						"inet-conn": {
							Name:      "inet-conn",
							Status:    "OK",
							Error:     "",
							CheckedAt: currentTs.UnixMilli(),
						},
						"disk-check": {
							Name:      "disk-check",
							Status:    "FAILED",
							Error:     "Critical: disk usage too high 61.93 percent",
							CheckedAt: currentTs.UnixMilli(),
						},
					},
				},
			}
			checkRes = &healthcheck.CheckResult{
				Success: system.Success{
					Code:    1000,
					Message: "success check health",
				},
				Status: "WARNING",
				Items: map[string]healthcheck.CheckResultItem{
					"inet-conn": {
						Name:      "inet-conn",
						Status:    "OK",
						Error:     "",
						Fatal:     false,
						CheckedAt: currentTs,
					},
					"disk-check": {
						Name:      "disk-check",
						Status:    "FAILED",
						Error:     "Critical: disk usage too high 61.93 percent",
						Fatal:     false,
						CheckedAt: currentTs,
					},
				},
			}
		})

		When("failed check service health", func() {
			It("should return error", func() {
				healthService.
					EXPECT().
					Check(gomock.Eq(ctx)).
					Return(nil, &system.Error{
						Code:    1001,
						Message: "routine error",
					}).
					Times(1)

				res, err := handler.CheckHealth(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(
					grpc_status.Error(codes.Unknown, "routine error"),
				))
			})
		})

		When("there is no states available", func() {
			It("should return result", func() {
				checkRes := &healthcheck.CheckResult{
					Success: system.Success{
						Code:    1000,
						Message: "success check health",
					},
					Status: "OK",
					Items:  map[string]healthcheck.CheckResultItem{},
				}
				healthService.
					EXPECT().
					Check(gomock.Eq(ctx)).
					Return(checkRes, nil).
					Times(1)

				res, err := handler.CheckHealth(ctx, p)

				r := &api.CheckHealthResult{
					Code:    1000,
					Message: "success check health",
					Data: &api.CheckHealthData{
						Status:  "OK",
						Details: map[string]*api.CheckHealthDetail{},
					},
				}
				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})

		When("there are states available", func() {
			It("should return result", func() {
				healthService.
					EXPECT().
					Check(gomock.Eq(ctx)).
					Return(checkRes, nil).
					Times(1)

				res, err := handler.CheckHealth(ctx, p)

				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})
	})
})
