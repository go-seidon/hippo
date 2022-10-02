package grpc_app_test

import (
	"context"
	"fmt"
	"io"
	"time"

	grpc_v1 "github.com/go-seidon/local/generated/proto/api/grpc/v1"
	mock_grpcv1 "github.com/go-seidon/local/generated/proto/api/grpc/v1/mock"
	mock_context "github.com/go-seidon/local/internal/context/mock"
	"github.com/go-seidon/local/internal/file"
	mock_file "github.com/go-seidon/local/internal/file/mock"
	grpc_app "github.com/go-seidon/local/internal/grpc-app"
	"github.com/go-seidon/local/internal/healthcheck"
	mock_healthcheck "github.com/go-seidon/local/internal/healthcheck/mock"
	mock_io "github.com/go-seidon/local/internal/io/mock"
	"github.com/go-seidon/local/internal/status"
	"github.com/go-seidon/local/internal/validation"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/metadata"
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
			config := &grpc_app.GrpcAppConfig{}
			handler = grpc_app.NewFileHandler(fileService, config)
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

		When("there are invalid param", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Eq(delParam)).
					Return(nil, validation.Error("invalid data")).
					Times(1)

				res, err := handler.DeleteFile(ctx, p)

				expectRes := &grpc_v1.DeleteFileResult{
					Code:    1002,
					Message: "invalid data",
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

	Context("RetrieveFile function", Label("unit"), func() {
		var (
			handler     grpc_v1.FileServiceServer
			fileService *mock_file.MockFile
			ctx         *mock_context.MockContext
			p           *grpc_v1.RetrieveFileParam
			pendingRes  *grpc_v1.RetrieveFileResult
			canceledRes *grpc_v1.RetrieveFileResult
			failedRes   *grpc_v1.RetrieveFileResult
			successRes  *grpc_v1.RetrieveFileResult
			stream      *mock_grpcv1.MockFileService_RetrieveFileServer
			retParam    file.RetrieveFileParam
			retRes      *file.RetrieveFileResult
			rc          *mock_io.MockReadCloser
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileService = mock_file.NewMockFile(ctrl)
			config := &grpc_app.GrpcAppConfig{}
			handler = grpc_app.NewFileHandler(fileService, config)
			ctx = mock_context.NewMockContext(ctrl)
			p = &grpc_v1.RetrieveFileParam{
				FileId: "file-id",
			}
			pendingRes = &grpc_v1.RetrieveFileResult{
				Code:    1005,
				Message: "retrieving file",
			}
			canceledRes = &grpc_v1.RetrieveFileResult{
				Code:    1001,
				Message: context.Canceled.Error(),
			}
			failedRes = &grpc_v1.RetrieveFileResult{
				Code:    1001,
				Message: "i/o error",
			}
			successRes = &grpc_v1.RetrieveFileResult{
				Code:    1000,
				Message: "success retrieve file",
			}
			stream = mock_grpcv1.NewMockFileService_RetrieveFileServer(ctrl)
			retParam = file.RetrieveFileParam{
				FileId: "file-id",
			}
			rc = mock_io.NewMockReadCloser(ctrl)
			retRes = &file.RetrieveFileResult{
				UniqueId:  "file-id",
				Name:      "file-name",
				MimeType:  "image/jpeg",
				Extension: "jpg",
				Data:      rc,
			}
			stream.
				EXPECT().
				Context().
				Return(ctx).
				Times(1)
		})

		When("failed send stream during file is not found", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(nil, file.ErrorNotFound).
					Times(1)

				res := &grpc_v1.RetrieveFileResult{
					Code:    1004,
					Message: file.ErrorNotFound.Error(),
				}
				stream.
					EXPECT().
					Send(gomock.Eq(res)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("file is not found", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(nil, file.ErrorNotFound).
					Times(1)

				res := &grpc_v1.RetrieveFileResult{
					Code:    1004,
					Message: file.ErrorNotFound.Error(),
				}
				stream.
					EXPECT().
					Send(gomock.Eq(res)).
					Return(nil).
					Times(1)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(BeNil())
			})
		})

		When("failed send stream during invalid data", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(nil, validation.Error("invalid data")).
					Times(1)

				res := &grpc_v1.RetrieveFileResult{
					Code:    1002,
					Message: "invalid data",
				}
				stream.
					EXPECT().
					Send(gomock.Eq(res)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("there are invalid data", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(nil, validation.Error("invalid data")).
					Times(1)

				res := &grpc_v1.RetrieveFileResult{
					Code:    1002,
					Message: "invalid data",
				}
				stream.
					EXPECT().
					Send(gomock.Eq(res)).
					Return(nil).
					Times(1)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(BeNil())
			})
		})

		When("failed send stream during failed retrieve file", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(nil, fmt.Errorf("db error")).
					Times(1)

				res := &grpc_v1.RetrieveFileResult{
					Code:    1001,
					Message: "db error",
				}
				stream.
					EXPECT().
					Send(gomock.Eq(res)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("failed retrieve file", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(nil, fmt.Errorf("db error")).
					Times(1)

				res := &grpc_v1.RetrieveFileResult{
					Code:    1001,
					Message: "db error",
				}
				stream.
					EXPECT().
					Send(gomock.Eq(res)).
					Return(nil).
					Times(1)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(BeNil())
			})
		})

		When("failed send header during success retrieve file", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(retRes, nil).
					Times(1)

				md := metadata.New(map[string]string{
					"file_name":      retRes.Name,
					"file_mimetype":  retRes.MimeType,
					"file_extension": retRes.Extension,
					"file_size":      fmt.Sprintf("%d", retRes.Size),
				})
				stream.
					EXPECT().
					SendHeader(gomock.Eq(md)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("failed send pending state during success retrieve file", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(retRes, nil).
					Times(1)

				md := metadata.New(map[string]string{
					"file_name":      retRes.Name,
					"file_mimetype":  retRes.MimeType,
					"file_extension": retRes.Extension,
					"file_size":      fmt.Sprintf("%d", retRes.Size),
				})
				stream.
					EXPECT().
					SendHeader(gomock.Eq(md)).
					Return(nil).
					Times(1)

				stream.
					EXPECT().
					Send(gomock.Eq(pendingRes)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("failed send stream during action is cancelled by client", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(retRes, nil).
					Times(1)

				stream.
					EXPECT().
					Send(gomock.Eq(pendingRes)).
					Return(nil).
					Times(1)

				md := metadata.New(map[string]string{
					"file_name":      retRes.Name,
					"file_mimetype":  retRes.MimeType,
					"file_extension": retRes.Extension,
					"file_size":      fmt.Sprintf("%d", retRes.Size),
				})
				stream.
					EXPECT().
					SendHeader(gomock.Eq(md)).
					Return(nil).
					Times(1)

				rc.
					EXPECT().
					Close().
					Times(1)

				ctx.
					EXPECT().
					Err().
					Return(context.Canceled).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				stream.
					EXPECT().
					Send(gomock.Eq(canceledRes)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("action is cancelled by client", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(retRes, nil).
					Times(1)

				stream.
					EXPECT().
					Send(gomock.Eq(pendingRes)).
					Return(nil).
					Times(1)

				md := metadata.New(map[string]string{
					"file_name":      retRes.Name,
					"file_mimetype":  retRes.MimeType,
					"file_extension": retRes.Extension,
					"file_size":      fmt.Sprintf("%d", retRes.Size),
				})
				stream.
					EXPECT().
					SendHeader(gomock.Eq(md)).
					Return(nil).
					Times(1)

				rc.
					EXPECT().
					Close().
					Times(1)

				ctx.
					EXPECT().
					Err().
					Return(context.Canceled).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				stream.
					EXPECT().
					Send(gomock.Eq(canceledRes)).
					Return(nil).
					Times(1)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(BeNil())
			})
		})

		When("failed send action fail during read data", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(retRes, nil).
					Times(1)

				stream.
					EXPECT().
					Send(gomock.Eq(pendingRes)).
					Return(nil).
					Times(1)

				md := metadata.New(map[string]string{
					"file_name":      retRes.Name,
					"file_mimetype":  retRes.MimeType,
					"file_extension": retRes.Extension,
					"file_size":      fmt.Sprintf("%d", retRes.Size),
				})
				stream.
					EXPECT().
					SendHeader(gomock.Eq(md)).
					Return(nil).
					Times(1)

				rc.
					EXPECT().
					Close().
					Times(1)

				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				rc.
					EXPECT().
					Read(gomock.Any()).
					Return(0, fmt.Errorf("i/o error")).
					Times(1)

				stream.
					EXPECT().
					Send(gomock.Eq(failedRes)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("success send action fail during read data", func() {
			It("should return result", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(retRes, nil).
					Times(1)

				stream.
					EXPECT().
					Send(gomock.Eq(pendingRes)).
					Return(nil).
					Times(1)

				md := metadata.New(map[string]string{
					"file_name":      retRes.Name,
					"file_mimetype":  retRes.MimeType,
					"file_extension": retRes.Extension,
					"file_size":      fmt.Sprintf("%d", retRes.Size),
				})
				stream.
					EXPECT().
					SendHeader(gomock.Eq(md)).
					Return(nil).
					Times(1)

				rc.
					EXPECT().
					Close().
					Times(1)

				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				rc.
					EXPECT().
					Read(gomock.Any()).
					Return(0, fmt.Errorf("i/o error")).
					Times(1)

				stream.
					EXPECT().
					Send(gomock.Eq(failedRes)).
					Return(nil).
					Times(1)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(BeNil())
			})
		})

		When("failed send stream during success read data", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(retRes, nil).
					Times(1)

				stream.
					EXPECT().
					Send(gomock.Eq(pendingRes)).
					Return(nil).
					Times(1)

				md := metadata.New(map[string]string{
					"file_name":      retRes.Name,
					"file_mimetype":  retRes.MimeType,
					"file_extension": retRes.Extension,
					"file_size":      fmt.Sprintf("%d", retRes.Size),
				})
				stream.
					EXPECT().
					SendHeader(gomock.Eq(md)).
					Return(nil).
					Times(1)

				rc.
					EXPECT().
					Close().
					Times(1)

				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				rc.
					EXPECT().
					Read(gomock.Any()).
					Return(1, nil).
					Times(1)

				stream.
					EXPECT().
					Send(gomock.Any()).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("failed send stream during end of data", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(retRes, nil).
					Times(1)

				stream.
					EXPECT().
					Send(gomock.Eq(pendingRes)).
					Return(nil).
					Times(1)

				md := metadata.New(map[string]string{
					"file_name":      retRes.Name,
					"file_mimetype":  retRes.MimeType,
					"file_extension": retRes.Extension,
					"file_size":      fmt.Sprintf("%d", retRes.Size),
				})
				stream.
					EXPECT().
					SendHeader(gomock.Eq(md)).
					Return(nil).
					Times(1)

				rc.
					EXPECT().
					Close().
					Times(1)

				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(2)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(2)

				firstRead := rc.
					EXPECT().
					Read(gomock.Any()).
					Return(1, nil).
					Times(1)

				lastRead := rc.
					EXPECT().
					Read(gomock.Any()).
					Return(0, io.EOF).
					Times(1)

				gomock.InOrder(firstRead, lastRead)

				firstStream := stream.
					EXPECT().
					Send(gomock.Any()).
					Return(nil).
					Times(1)

				lastStream := stream.
					EXPECT().
					Send(gomock.Eq(successRes)).
					Return(fmt.Errorf("network error")).
					Times(1)

				gomock.InOrder(firstStream, lastStream)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("success send stream during end of data", func() {
			It("should return result", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(retRes, nil).
					Times(1)

				stream.
					EXPECT().
					Send(gomock.Eq(pendingRes)).
					Return(nil).
					Times(1)

				md := metadata.New(map[string]string{
					"file_name":      retRes.Name,
					"file_mimetype":  retRes.MimeType,
					"file_extension": retRes.Extension,
					"file_size":      fmt.Sprintf("%d", retRes.Size),
				})
				stream.
					EXPECT().
					SendHeader(gomock.Eq(md)).
					Return(nil).
					Times(1)

				rc.
					EXPECT().
					Close().
					Times(1)

				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(2)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(2)

				firstRead := rc.
					EXPECT().
					Read(gomock.Any()).
					Return(1, nil).
					Times(1)

				lastRead := rc.
					EXPECT().
					Read(gomock.Any()).
					Return(0, io.EOF).
					Times(1)

				gomock.InOrder(firstRead, lastRead)

				firstStream := stream.
					EXPECT().
					Send(gomock.Any()).
					Return(nil).
					Times(1)

				lastStream := stream.
					EXPECT().
					Send(gomock.Eq(successRes)).
					Return(nil).
					Times(1)

				gomock.InOrder(firstStream, lastStream)

				err := handler.RetrieveFile(p, stream)

				Expect(err).To(BeNil())
			})
		})
	})

	Context("UploadFile function", Label("unit"), func() {
		var (
			handler     grpc_v1.FileServiceServer
			fileService *mock_file.MockFile
			ctx         *mock_context.MockContext
			currentTs   time.Time
			stream      *mock_grpcv1.MockFileService_UploadFileServer
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileService = mock_file.NewMockFile(ctrl)
			config := &grpc_app.GrpcAppConfig{
				UploadFormSize: 1073741824, //1GB
			}
			handler = grpc_app.NewFileHandler(fileService, config)
			ctx = mock_context.NewMockContext(ctrl)
			currentTs = time.Now()
			stream = mock_grpcv1.NewMockFileService_UploadFileServer(ctrl)
		})

		When("failed send stream during action cancelled by client", func() {
			It("should return error", func() {
				ctx.
					EXPECT().
					Err().
					Return(context.Canceled).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				failedRes := &grpc_v1.UploadFileResult{
					Code:    status.ACTION_FAILED,
					Message: context.Canceled.Error(),
				}
				stream.
					EXPECT().
					SendAndClose(gomock.Eq(failedRes)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.UploadFile(stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("action cancelled by client", func() {
			It("should return error", func() {
				ctx.
					EXPECT().
					Err().
					Return(context.Canceled).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				failedRes := &grpc_v1.UploadFileResult{
					Code:    status.ACTION_FAILED,
					Message: context.Canceled.Error(),
				}
				stream.
					EXPECT().
					SendAndClose(gomock.Eq(failedRes)).
					Return(nil).
					Times(1)

				err := handler.UploadFile(stream)

				Expect(err).To(BeNil())
			})
		})

		When("failed send stream during failed receive stream", func() {
			It("should return error", func() {
				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				stream.
					EXPECT().
					Recv().
					Return(nil, fmt.Errorf("client cancelled")).
					Times(1)

				failedRes := &grpc_v1.UploadFileResult{
					Code:    status.ACTION_FAILED,
					Message: "client cancelled",
				}
				stream.
					EXPECT().
					SendAndClose(gomock.Eq(failedRes)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.UploadFile(stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("failed receive stream", func() {
			It("should return error", func() {
				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				stream.
					EXPECT().
					Recv().
					Return(nil, fmt.Errorf("client cancelled")).
					Times(1)

				failedRes := &grpc_v1.UploadFileResult{
					Code:    status.ACTION_FAILED,
					Message: "client cancelled",
				}
				stream.
					EXPECT().
					SendAndClose(gomock.Eq(failedRes)).
					Return(nil).
					Times(1)

				err := handler.UploadFile(stream)

				Expect(err).To(BeNil())
			})
		})

		When("failed send stream during max file size reached", func() {
			It("should return error", func() {
				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				config := &grpc_app.GrpcAppConfig{
					UploadFormSize: 1,
				}
				handler := grpc_app.NewFileHandler(fileService, config)

				param := &grpc_v1.UploadFileParam{
					Data: &grpc_v1.UploadFileParam_Chunks{
						Chunks: []byte{1, 2, 3},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(param, nil).
					Times(1)

				failedRes := &grpc_v1.UploadFileResult{
					Code:    status.ACTION_FAILED,
					Message: "file is too large",
				}
				stream.
					EXPECT().
					SendAndClose(gomock.Eq(failedRes)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.UploadFile(stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("max file size reached", func() {
			It("should return error", func() {
				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				config := &grpc_app.GrpcAppConfig{
					UploadFormSize: 1,
				}
				handler := grpc_app.NewFileHandler(fileService, config)

				param := &grpc_v1.UploadFileParam{
					Data: &grpc_v1.UploadFileParam_Chunks{
						Chunks: []byte{1, 2, 3},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(param, nil).
					Times(1)

				failedRes := &grpc_v1.UploadFileResult{
					Code:    status.ACTION_FAILED,
					Message: "file is too large",
				}
				stream.
					EXPECT().
					SendAndClose(gomock.Eq(failedRes)).
					Return(nil).
					Times(1)

				err := handler.UploadFile(stream)

				Expect(err).To(BeNil())
			})
		})

		When("failed send stream during failed upload file", func() {
			It("should return error", func() {
				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(3)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(3)

				infoParam := &grpc_v1.UploadFileParam{
					Data: &grpc_v1.UploadFileParam_Info{
						Info: &grpc_v1.UploadFileInfo{
							Name:      "file-name",
							Mimetype:  "file-mimetype",
							Extension: "file-extension",
						},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(infoParam, nil).
					Times(1)

				chunkParam := &grpc_v1.UploadFileParam{
					Data: &grpc_v1.UploadFileParam_Chunks{
						Chunks: []byte{1, 2, 3},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(chunkParam, nil).
					Times(1)

				stream.
					EXPECT().
					Recv().
					Return(nil, io.EOF).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				fileService.
					EXPECT().
					UploadFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, fmt.Errorf("db error")).
					Times(1)

				failedRes := &grpc_v1.UploadFileResult{
					Code:    status.ACTION_FAILED,
					Message: "db error",
				}
				stream.
					EXPECT().
					SendAndClose(gomock.Eq(failedRes)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.UploadFile(stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("failed upload file", func() {
			It("should return error", func() {
				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(3)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(3)

				infoParam := &grpc_v1.UploadFileParam{
					Data: &grpc_v1.UploadFileParam_Info{
						Info: &grpc_v1.UploadFileInfo{
							Name:      "file-name",
							Mimetype:  "file-mimetype",
							Extension: "file-extension",
						},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(infoParam, nil).
					Times(1)

				chunkParam := &grpc_v1.UploadFileParam{
					Data: &grpc_v1.UploadFileParam_Chunks{
						Chunks: []byte{1, 2, 3},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(chunkParam, nil).
					Times(1)

				stream.
					EXPECT().
					Recv().
					Return(nil, io.EOF).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				fileService.
					EXPECT().
					UploadFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, fmt.Errorf("db error")).
					Times(1)

				failedRes := &grpc_v1.UploadFileResult{
					Code:    status.ACTION_FAILED,
					Message: "db error",
				}
				stream.
					EXPECT().
					SendAndClose(gomock.Eq(failedRes)).
					Return(nil).
					Times(1)

				err := handler.UploadFile(stream)

				Expect(err).To(BeNil())
			})
		})

		When("failed send stream during invalid data", func() {
			It("should return error", func() {
				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(3)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(3)

				infoParam := &grpc_v1.UploadFileParam{
					Data: &grpc_v1.UploadFileParam_Info{
						Info: &grpc_v1.UploadFileInfo{
							Name:      "file-name",
							Mimetype:  "file-mimetype",
							Extension: "file-extension",
						},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(infoParam, nil).
					Times(1)

				chunkParam := &grpc_v1.UploadFileParam{
					Data: &grpc_v1.UploadFileParam_Chunks{
						Chunks: []byte{1, 2, 3},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(chunkParam, nil).
					Times(1)

				stream.
					EXPECT().
					Recv().
					Return(nil, io.EOF).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				fileService.
					EXPECT().
					UploadFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, validation.Error("invalid data")).
					Times(1)

				failedRes := &grpc_v1.UploadFileResult{
					Code:    status.INVALID_PARAM,
					Message: "invalid data",
				}
				stream.
					EXPECT().
					SendAndClose(gomock.Eq(failedRes)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.UploadFile(stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("there are invalid data", func() {
			It("should return error", func() {
				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(3)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(3)

				infoParam := &grpc_v1.UploadFileParam{
					Data: &grpc_v1.UploadFileParam_Info{
						Info: &grpc_v1.UploadFileInfo{
							Name:      "file-name",
							Mimetype:  "file-mimetype",
							Extension: "file-extension",
						},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(infoParam, nil).
					Times(1)

				chunkParam := &grpc_v1.UploadFileParam{
					Data: &grpc_v1.UploadFileParam_Chunks{
						Chunks: []byte{1, 2, 3},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(chunkParam, nil).
					Times(1)

				stream.
					EXPECT().
					Recv().
					Return(nil, io.EOF).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				fileService.
					EXPECT().
					UploadFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, validation.Error("invalid data")).
					Times(1)

				failedRes := &grpc_v1.UploadFileResult{
					Code:    status.INVALID_PARAM,
					Message: "invalid data",
				}
				stream.
					EXPECT().
					SendAndClose(gomock.Eq(failedRes)).
					Return(nil).
					Times(1)

				err := handler.UploadFile(stream)

				Expect(err).To(BeNil())
			})
		})

		When("failed send stream during success upload file", func() {
			It("should return error", func() {
				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(3)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(3)

				infoParam := &grpc_v1.UploadFileParam{
					Data: &grpc_v1.UploadFileParam_Info{
						Info: &grpc_v1.UploadFileInfo{
							Name:      "file-name",
							Mimetype:  "file-mimetype",
							Extension: "file-extension",
						},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(infoParam, nil).
					Times(1)

				chunkParam := &grpc_v1.UploadFileParam{
					Data: &grpc_v1.UploadFileParam_Chunks{
						Chunks: []byte{1, 2, 3},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(chunkParam, nil).
					Times(1)

				stream.
					EXPECT().
					Recv().
					Return(nil, io.EOF).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				uploadRes := &file.UploadFileResult{
					UniqueId:   "file-id",
					Name:       "file-name",
					Path:       "file/path",
					Mimetype:   "file-mime-type",
					Extension:  "jpeg",
					Size:       100,
					UploadedAt: currentTs,
				}
				fileService.
					EXPECT().
					UploadFile(gomock.Eq(ctx), gomock.Any()).
					Return(uploadRes, nil).
					Times(1)

				failedRes := &grpc_v1.UploadFileResult{
					Code:    status.ACTION_SUCCESS,
					Message: "success upload file",
					Data: &grpc_v1.UploadFileData{
						Id:         uploadRes.UniqueId,
						Name:       uploadRes.Name,
						Path:       uploadRes.Path,
						Mimetype:   uploadRes.Mimetype,
						Extension:  uploadRes.Extension,
						Size:       uploadRes.Size,
						UploadedAt: uploadRes.UploadedAt.UnixMilli(),
					},
				}
				stream.
					EXPECT().
					SendAndClose(gomock.Eq(failedRes)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.UploadFile(stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("success upload file", func() {
			It("should return result", func() {
				ctx.
					EXPECT().
					Err().
					Return(nil).
					Times(3)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(3)

				infoParam := &grpc_v1.UploadFileParam{
					Data: &grpc_v1.UploadFileParam_Info{
						Info: &grpc_v1.UploadFileInfo{
							Name:      "file-name",
							Mimetype:  "file-mimetype",
							Extension: "file-extension",
						},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(infoParam, nil).
					Times(1)

				chunkParam := &grpc_v1.UploadFileParam{
					Data: &grpc_v1.UploadFileParam_Chunks{
						Chunks: []byte{1, 2, 3},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(chunkParam, nil).
					Times(1)

				stream.
					EXPECT().
					Recv().
					Return(nil, io.EOF).
					Times(1)

				stream.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				uploadRes := &file.UploadFileResult{
					UniqueId:   "file-id",
					Name:       "file-name",
					Path:       "file/path",
					Mimetype:   "file-mime-type",
					Extension:  "jpeg",
					Size:       100,
					UploadedAt: currentTs,
				}
				fileService.
					EXPECT().
					UploadFile(gomock.Eq(ctx), gomock.Any()).
					Return(uploadRes, nil).
					Times(1)

				failedRes := &grpc_v1.UploadFileResult{
					Code:    status.ACTION_SUCCESS,
					Message: "success upload file",
					Data: &grpc_v1.UploadFileData{
						Id:         uploadRes.UniqueId,
						Name:       uploadRes.Name,
						Path:       uploadRes.Path,
						Mimetype:   uploadRes.Mimetype,
						Extension:  uploadRes.Extension,
						Size:       uploadRes.Size,
						UploadedAt: uploadRes.UploadedAt.UnixMilli(),
					},
				}
				stream.
					EXPECT().
					SendAndClose(gomock.Eq(failedRes)).
					Return(nil).
					Times(1)

				err := handler.UploadFile(stream)

				Expect(err).To(BeNil())
			})
		})
	})
})
