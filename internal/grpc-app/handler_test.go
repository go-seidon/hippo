package grpc_app_test

import (
	"context"
	"fmt"
	"time"

	grpc_v1 "github.com/go-seidon/local/generated/proto/api/grpc/v1"
	"github.com/go-seidon/local/internal/file"
	mock_file "github.com/go-seidon/local/internal/file/mock"
	grpc_app "github.com/go-seidon/local/internal/grpc-app"
	"github.com/go-seidon/local/internal/healthcheck"
	mock_healthcheck "github.com/go-seidon/local/internal/healthcheck/mock"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler Package", func() {

	Context("CheckHealth function", Label("unit"), func() {
		var (
			handler       grpc_v1.HealthServiceServer
			healthService *mock_healthcheck.MockHealthCheck
			ctx           context.Context
			p             *grpc_v1.CheckHealthParam
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			healthService = mock_healthcheck.NewMockHealthCheck(ctrl)
			handler = grpc_app.NewHealthHandler(healthService)
			ctx = context.Background()
			p = &grpc_v1.CheckHealthParam{}
		})

		When("failed check service health", func() {
			It("should return error", func() {
				expectedErr := fmt.Errorf("routine error")

				healthService.
					EXPECT().
					Check().
					Return(nil, expectedErr).
					Times(1)

				res, err := handler.CheckHealth(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(
					grpc_status.Error(codes.Unknown, expectedErr.Error()),
				))
			})
		})

		When("no health check states available", func() {
			It("should return result", func() {
				checkRes := &healthcheck.CheckResult{
					Status: "OK",
					Items:  map[string]healthcheck.CheckResultItem{},
				}

				healthService.
					EXPECT().
					Check().
					Return(checkRes, nil).
					Times(1)

				res, err := handler.CheckHealth(ctx, p)

				expectedRes := &grpc_v1.CheckHealthResult{
					Code:    1000,
					Message: "success check service health",
					Data: &grpc_v1.CheckHealthData{
						Status:  "OK",
						Details: map[string]*grpc_v1.CheckHealthDetail{},
					},
				}

				Expect(res).To(Equal(expectedRes))
				Expect(err).To(BeNil())
			})
		})

		When("there are health check states", func() {
			It("should return result", func() {
				currentTimestamp := time.Now()

				checkRes := &healthcheck.CheckResult{
					Status: "WARNING",
					Items: map[string]healthcheck.CheckResultItem{
						"inet-conn": {
							Name:      "inet-conn",
							Status:    "OK",
							Error:     "",
							Fatal:     false,
							CheckedAt: currentTimestamp,
						},
						"disk-check": {
							Name:      "disk-check",
							Status:    "FAILED",
							Error:     "Critical: disk usage too high 61.93 percent",
							Fatal:     false,
							CheckedAt: currentTimestamp,
						},
					},
				}

				healthService.
					EXPECT().
					Check().
					Return(checkRes, nil).
					Times(1)

				res, err := handler.CheckHealth(ctx, p)

				expectedRes := &grpc_v1.CheckHealthResult{
					Code:    1000,
					Message: "success check service health",
					Data: &grpc_v1.CheckHealthData{
						Status: "WARNING",
						Details: map[string]*grpc_v1.CheckHealthDetail{
							"inet-conn": {
								Name:      "inet-conn",
								Status:    "OK",
								Error:     "",
								CheckedAt: currentTimestamp.UnixMilli(),
							},
							"disk-check": {
								Name:      "disk-check",
								Status:    "FAILED",
								Error:     "Critical: disk usage too high 61.93 percent",
								CheckedAt: currentTimestamp.UnixMilli(),
							},
						},
					},
				}

				Expect(res).To(Equal(expectedRes))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("DeleteFile function", Label("unit"), func() {
		var (
			handler     grpc_v1.FileServiceServer
			fileService *mock_file.MockFile
			ctx         context.Context
			currentTs   time.Time
			p           *grpc_v1.DeleteFileParam
			delParam    file.DeleteFileParam
			delRes      *file.DeleteFileResult
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileService = mock_file.NewMockFile(ctrl)
			handler = grpc_app.NewFileHandler(fileService)
			ctx = context.Background()
			currentTs = time.Now()
			p = &grpc_v1.DeleteFileParam{
				FileId: "file-id",
			}
			delParam = file.DeleteFileParam{
				FileId: "file-id",
			}
			delRes = &file.DeleteFileResult{
				DeletedAt: currentTs,
			}
		})

		When("success delete file", func() {
			It("should return result", func() {
				fileService.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Eq(delParam)).
					Return(delRes, nil).
					Times(1)

				res, err := handler.DeleteFile(ctx, p)

				expectRes := &grpc_v1.DeleteFileResult{
					Code:    1000,
					Message: "success delete file",
					Data: &grpc_v1.DeleteFileData{
						DeletedAt: currentTs.UnixMilli(),
					},
				}
				Expect(res).To(Equal(expectRes))
				Expect(err).To(BeNil())
			})
		})

		When("file is not found", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Eq(delParam)).
					Return(nil, file.ErrorNotFound).
					Times(1)

				res, err := handler.DeleteFile(ctx, p)

				expectRes := &grpc_v1.DeleteFileResult{
					Code:    1004,
					Message: "not found",
				}
				Expect(res).To(Equal(expectRes))
				Expect(err).To(BeNil())
			})
		})

		When("failed delete file", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Eq(delParam)).
					Return(nil, fmt.Errorf("db error")).
					Times(1)

				res, err := handler.DeleteFile(ctx, p)

				expectRes := &grpc_v1.DeleteFileResult{
					Code:    1001,
					Message: "db error",
				}
				Expect(res).To(Equal(expectRes))
				Expect(err).To(BeNil())
			})
		})
	})
})
