package grpcapp_test

import (
	"context"
	"fmt"
	"io"
	"time"

	api "github.com/go-seidon/hippo/api/grpcapp"
	mock_grpcapp "github.com/go-seidon/hippo/api/grpcapp/mock"
	"github.com/go-seidon/hippo/internal/file"
	mock_file "github.com/go-seidon/hippo/internal/file/mock"
	"github.com/go-seidon/hippo/internal/grpcapp"
	"github.com/go-seidon/hippo/internal/healthcheck"
	mock_healthcheck "github.com/go-seidon/hippo/internal/healthcheck/mock"
	mock_context "github.com/go-seidon/provider/context/mock"
	mock_io "github.com/go-seidon/provider/io/mock"
	"github.com/go-seidon/provider/system"
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
			handler = grpcapp.NewHealthHandler(healthService)
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

	Context("DeleteFileById function", Label("unit"), func() {
		var (
			handler     api.FileServiceServer
			fileService *mock_file.MockFile
			ctx         context.Context
			currentTs   time.Time
			p           *api.DeleteFileByIdParam
			r           *api.DeleteFileByIdResult
			delParam    file.DeleteFileParam
			delRes      *file.DeleteFileResult
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileService = mock_file.NewMockFile(ctrl)
			config := &grpcapp.GrpcAppConfig{}
			handler = grpcapp.NewFileHandler(fileService, config)
			ctx = context.Background()
			currentTs = time.Now()
			p = &api.DeleteFileByIdParam{
				FileId: "file-id",
			}
			r = &api.DeleteFileByIdResult{
				Code:    1000,
				Message: "success delete file",
				Data: &api.DeleteFileByIdData{
					DeletedAt: currentTs.UnixMilli(),
				},
			}
			delParam = file.DeleteFileParam{
				FileId: "file-id",
			}
			delRes = &file.DeleteFileResult{
				Success: system.Success{
					Code:    1000,
					Message: "success delete file",
				},
				DeletedAt: currentTs,
			}
		})

		When("file is not found", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Eq(delParam)).
					Return(nil, &system.Error{
						Code:    1004,
						Message: "file is not found",
					}).
					Times(1)

				res, err := handler.DeleteFileById(ctx, p)

				r := &api.DeleteFileByIdResult{
					Code:    1004,
					Message: "file is not found",
				}
				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})

		When("there is invalid param", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Eq(delParam)).
					Return(nil, &system.Error{
						Code:    1002,
						Message: "invalid data",
					}).
					Times(1)

				res, err := handler.DeleteFileById(ctx, p)

				r := &api.DeleteFileByIdResult{
					Code:    1002,
					Message: "invalid data",
				}
				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})

		When("failed delete file", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Eq(delParam)).
					Return(nil, &system.Error{
						Code:    1001,
						Message: "network error",
					}).
					Times(1)

				res, err := handler.DeleteFileById(ctx, p)

				r := &api.DeleteFileByIdResult{
					Code:    1001,
					Message: "network error",
				}
				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})

		When("success delete file", func() {
			It("should return result", func() {
				fileService.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Eq(delParam)).
					Return(delRes, nil).
					Times(1)

				res, err := handler.DeleteFileById(ctx, p)

				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("RetrieveFileById function", Label("unit"), func() {
		var (
			handler     api.FileServiceServer
			fileService *mock_file.MockFile
			ctx         *mock_context.MockContext
			p           *api.RetrieveFileByIdParam
			pendingRes  *api.RetrieveFileByIdResult
			canceledRes *api.RetrieveFileByIdResult
			failedRes   *api.RetrieveFileByIdResult
			successRes  *api.RetrieveFileByIdResult
			stream      *mock_grpcapp.MockFileService_RetrieveFileByIdServer
			retParam    file.RetrieveFileParam
			retRes      *file.RetrieveFileResult
			rc          *mock_io.MockReadCloser
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileService = mock_file.NewMockFile(ctrl)
			config := &grpcapp.GrpcAppConfig{}
			handler = grpcapp.NewFileHandler(fileService, config)
			ctx = mock_context.NewMockContext(ctrl)
			p = &api.RetrieveFileByIdParam{
				FileId: "file-id",
			}
			pendingRes = &api.RetrieveFileByIdResult{
				Code:    1005,
				Message: "retrieving file",
			}
			canceledRes = &api.RetrieveFileByIdResult{
				Code:    1001,
				Message: context.Canceled.Error(),
			}
			failedRes = &api.RetrieveFileByIdResult{
				Code:    1001,
				Message: "i/o error",
			}
			successRes = &api.RetrieveFileByIdResult{
				Code:    1000,
				Message: "success retrieve file",
			}
			stream = mock_grpcapp.NewMockFileService_RetrieveFileByIdServer(ctrl)
			retParam = file.RetrieveFileParam{
				FileId: "file-id",
			}
			rc = mock_io.NewMockReadCloser(ctrl)
			retRes = &file.RetrieveFileResult{
				Success: system.Success{
					Code:    1000,
					Message: "success retrieve file",
				},
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
					Return(nil, &system.Error{
						Code:    1004,
						Message: "file is not found",
					}).
					Times(1)

				res := &api.RetrieveFileByIdResult{
					Code:    1004,
					Message: "file is not found",
				}
				stream.
					EXPECT().
					Send(gomock.Eq(res)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.RetrieveFileById(p, stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("file is not found", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(nil, &system.Error{
						Code:    1004,
						Message: "file is not found",
					}).
					Times(1)

				res := &api.RetrieveFileByIdResult{
					Code:    1004,
					Message: "file is not found",
				}
				stream.
					EXPECT().
					Send(gomock.Eq(res)).
					Return(nil).
					Times(1)

				err := handler.RetrieveFileById(p, stream)

				Expect(err).To(BeNil())
			})
		})

		When("failed send stream during invalid data", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(nil, &system.Error{
						Code:    1002,
						Message: "invalid data",
					}).
					Times(1)

				res := &api.RetrieveFileByIdResult{
					Code:    1002,
					Message: "invalid data",
				}
				stream.
					EXPECT().
					Send(gomock.Eq(res)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.RetrieveFileById(p, stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("there are invalid data", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(nil, &system.Error{
						Code:    1002,
						Message: "invalid data",
					}).
					Times(1)

				res := &api.RetrieveFileByIdResult{
					Code:    1002,
					Message: "invalid data",
				}
				stream.
					EXPECT().
					Send(gomock.Eq(res)).
					Return(nil).
					Times(1)

				err := handler.RetrieveFileById(p, stream)

				Expect(err).To(BeNil())
			})
		})

		When("failed send stream during failed retrieve file", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(nil, &system.Error{
						Code:    1001,
						Message: "network error",
					}).
					Times(1)

				res := &api.RetrieveFileByIdResult{
					Code:    1001,
					Message: "network error",
				}
				stream.
					EXPECT().
					Send(gomock.Eq(res)).
					Return(fmt.Errorf("network error")).
					Times(1)

				err := handler.RetrieveFileById(p, stream)

				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("failed retrieve file", func() {
			It("should return error", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retParam)).
					Return(nil, &system.Error{
						Code:    1001,
						Message: "network error",
					}).
					Times(1)

				res := &api.RetrieveFileByIdResult{
					Code:    1001,
					Message: "network error",
				}
				stream.
					EXPECT().
					Send(gomock.Eq(res)).
					Return(nil).
					Times(1)

				err := handler.RetrieveFileById(p, stream)

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

				err := handler.RetrieveFileById(p, stream)

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

				err := handler.RetrieveFileById(p, stream)

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

				err := handler.RetrieveFileById(p, stream)

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

				err := handler.RetrieveFileById(p, stream)

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

				err := handler.RetrieveFileById(p, stream)

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

				err := handler.RetrieveFileById(p, stream)

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

				err := handler.RetrieveFileById(p, stream)

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

				err := handler.RetrieveFileById(p, stream)

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

				err := handler.RetrieveFileById(p, stream)

				Expect(err).To(BeNil())
			})
		})
	})

	Context("UploadFile function", Label("unit"), func() {
		var (
			handler     api.FileServiceServer
			fileService *mock_file.MockFile
			ctx         *mock_context.MockContext
			currentTs   time.Time
			stream      *mock_grpcapp.MockFileService_UploadFileServer
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileService = mock_file.NewMockFile(ctrl)
			config := &grpcapp.GrpcAppConfig{
				UploadFormSize: 1073741824, //1GB
			}
			handler = grpcapp.NewFileHandler(fileService, config)
			ctx = mock_context.NewMockContext(ctrl)
			currentTs = time.Now()
			stream = mock_grpcapp.NewMockFileService_UploadFileServer(ctrl)
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

				failedRes := &api.UploadFileResult{
					Code:    1001,
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

				failedRes := &api.UploadFileResult{
					Code:    1001,
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

				failedRes := &api.UploadFileResult{
					Code:    1001,
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

				failedRes := &api.UploadFileResult{
					Code:    1001,
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

				config := &grpcapp.GrpcAppConfig{
					UploadFormSize: 1,
				}
				handler := grpcapp.NewFileHandler(fileService, config)

				param := &api.UploadFileParam{
					Data: &api.UploadFileParam_Chunks{
						Chunks: []byte{1, 2, 3},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(param, nil).
					Times(1)

				failedRes := &api.UploadFileResult{
					Code:    1001,
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

				config := &grpcapp.GrpcAppConfig{
					UploadFormSize: 1,
				}
				handler := grpcapp.NewFileHandler(fileService, config)

				param := &api.UploadFileParam{
					Data: &api.UploadFileParam_Chunks{
						Chunks: []byte{1, 2, 3},
					},
				}
				stream.
					EXPECT().
					Recv().
					Return(param, nil).
					Times(1)

				failedRes := &api.UploadFileResult{
					Code:    1001,
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

				infoParam := &api.UploadFileParam{
					Data: &api.UploadFileParam_Info{
						Info: &api.UploadFileInfo{
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

				chunkParam := &api.UploadFileParam{
					Data: &api.UploadFileParam_Chunks{
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
					Return(nil, &system.Error{
						Code:    1001,
						Message: "network error",
					}).
					Times(1)

				failedRes := &api.UploadFileResult{
					Code:    1001,
					Message: "network error",
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

				infoParam := &api.UploadFileParam{
					Data: &api.UploadFileParam_Info{
						Info: &api.UploadFileInfo{
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

				chunkParam := &api.UploadFileParam{
					Data: &api.UploadFileParam_Chunks{
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
					Return(nil, &system.Error{
						Code:    1001,
						Message: "network error",
					}).
					Times(1)

				failedRes := &api.UploadFileResult{
					Code:    1001,
					Message: "network error",
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

				infoParam := &api.UploadFileParam{
					Data: &api.UploadFileParam_Info{
						Info: &api.UploadFileInfo{
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

				chunkParam := &api.UploadFileParam{
					Data: &api.UploadFileParam_Chunks{
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
					Return(nil, &system.Error{
						Code:    1002,
						Message: "invalid data",
					}).
					Times(1)

				failedRes := &api.UploadFileResult{
					Code:    1002,
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

		When("there is invalid data", func() {
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

				infoParam := &api.UploadFileParam{
					Data: &api.UploadFileParam_Info{
						Info: &api.UploadFileInfo{
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

				chunkParam := &api.UploadFileParam{
					Data: &api.UploadFileParam_Chunks{
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
					Return(nil, &system.Error{
						Code:    1002,
						Message: "invalid data",
					}).
					Times(1)

				failedRes := &api.UploadFileResult{
					Code:    1002,
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

				infoParam := &api.UploadFileParam{
					Data: &api.UploadFileParam_Info{
						Info: &api.UploadFileInfo{
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

				chunkParam := &api.UploadFileParam{
					Data: &api.UploadFileParam_Chunks{
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
					Success: system.Success{
						Code:    1000,
						Message: "success upload file",
					},
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

				failedRes := &api.UploadFileResult{
					Code:    1000,
					Message: "success upload file",
					Data: &api.UploadFileData{
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

				infoParam := &api.UploadFileParam{
					Data: &api.UploadFileParam_Info{
						Info: &api.UploadFileInfo{
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

				chunkParam := &api.UploadFileParam{
					Data: &api.UploadFileParam_Chunks{
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
					Success: system.Success{
						Code:    1000,
						Message: "success upload file",
					},
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

				failedRes := &api.UploadFileResult{
					Code:    1000,
					Message: "success upload file",
					Data: &api.UploadFileData{
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
