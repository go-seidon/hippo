package grpclog_test

import (
	"context"
	"net"
	"time"

	mock_grpc "github.com/go-seidon/hippo/internal/grpc/mock"
	"github.com/go-seidon/hippo/internal/grpclog"
	mock_context "github.com/go-seidon/provider/context/mock"
	mock_datetime "github.com/go-seidon/provider/datetime/mock"
	mock_logging "github.com/go-seidon/provider/logging/mock"
	"github.com/go-seidon/provider/testdata"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var _ = Describe("Interceptor Package", func() {

	Context("UnaryServerInterceptor function", Label("unit"), func() {
		var (
			currentTs   time.Time
			ctx         *mock_context.MockContext
			req         interface{}
			info        *grpc.UnaryServerInfo
			logger      *mock_logging.MockLogger
			interceptor grpc.UnaryServerInterceptor
			handler     func(ctx context.Context, req interface{}) (interface{}, error)
			clock       *mock_datetime.MockClock
		)

		BeforeEach(func() {
			t := GinkgoT()
			currentTs = time.Now()
			req = struct{}{}
			info = &grpc.UnaryServerInfo{
				FullMethod: "/pkg.Service/MethodName",
			}
			ctrl := gomock.NewController(t)
			ctx = mock_context.NewMockContext(ctrl)
			logger = mock_logging.NewMockLogger(ctrl)
			clock = mock_datetime.NewMockClock(ctrl)
			handler = func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, nil
			}
		})

		When("grpc method is ignored", func() {
			It("should return result", func() {
				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				ctx.
					EXPECT().
					Deadline().
					Return(currentTs, false).
					Times(1)

				interceptor = grpclog.UnaryServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
					grpclog.IgnoredMethod([]string{
						info.FullMethod,
					}),
				)
				res, err := interceptor(ctx, req, info, handler)

				Expect(res).To(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("grpc method is not ignored", func() {
			It("should return result", func() {
				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				ctx.
					EXPECT().
					Deadline().
					Return(currentTs, false).
					Times(1)

				ctx.
					EXPECT().
					Value(gomock.Any()).
					Return(nil).
					Times(2)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Logf(
						gomock.Eq("info"),
						gomock.Eq("request: %s@%s"),
						gomock.Eq("pkg.Service"),
						gomock.Eq("MethodName"),
					).
					Times(1)

				interceptor = grpclog.UnaryServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
				)
				res, err := interceptor(ctx, req, info, handler)

				Expect(res).To(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("deadline occured", func() {
			It("should return result", func() {
				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				ctx.
					EXPECT().
					Deadline().
					Return(currentTs, true).
					Times(1)

				ctx.
					EXPECT().
					Value(gomock.Any()).
					Return(nil).
					Times(2)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Logf(
						gomock.Eq("info"),
						gomock.Eq("request: %s@%s"),
						gomock.Eq("pkg.Service"),
						gomock.Eq("MethodName"),
					).
					Times(1)

				interceptor = grpclog.UnaryServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
				)
				res, err := interceptor(ctx, req, info, handler)

				Expect(res).To(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("fatal error occured", func() {
			It("should return result", func() {
				expectErr := status.New(codes.Unknown, "fatal error").Err()

				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				ctx.
					EXPECT().
					Deadline().
					Return(currentTs, false).
					Times(1)

				errLog := map[string]interface{}{
					"error": expectErr,
				}
				logger.
					EXPECT().
					WithFields(gomock.Eq(errLog)).
					Return(logger).
					Times(1)

				ctx.
					EXPECT().
					Value(gomock.Any()).
					Return(nil).
					Times(2)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Logf(
						gomock.Eq("error"),
						gomock.Eq("request: %s@%s"),
						gomock.Eq("pkg.Service"),
						gomock.Eq("MethodName"),
					).
					Times(1)

				interceptor = grpclog.UnaryServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
				)
				res, err := interceptor(ctx, req, info, func(ctx context.Context, req interface{}) (interface{}, error) {
					return nil, expectErr
				})

				Expect(res).To(BeNil())
				Expect(err).To(Equal(expectErr))
			})
		})

		When("warning error occured", func() {
			It("should return result", func() {
				expectErr := status.New(codes.PermissionDenied, "warning error").Err()

				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				ctx.
					EXPECT().
					Deadline().
					Return(currentTs, false).
					Times(1)

				errLog := map[string]interface{}{
					"error": expectErr,
				}
				logger.
					EXPECT().
					WithFields(gomock.Eq(errLog)).
					Return(logger).
					Times(1)

				ctx.
					EXPECT().
					Value(gomock.Any()).
					Return(nil).
					Times(2)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Logf(
						gomock.Eq("warn"),
						gomock.Eq("request: %s@%s"),
						gomock.Eq("pkg.Service"),
						gomock.Eq("MethodName"),
					).
					Times(1)

				interceptor = grpclog.UnaryServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
				)
				res, err := interceptor(ctx, req, info, func(ctx context.Context, req interface{}) (interface{}, error) {
					return nil, expectErr
				})

				Expect(res).To(BeNil())
				Expect(err).To(Equal(expectErr))
			})
		})

		When("network info available", func() {
			It("should return result", func() {
				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				nctx := context.Background()
				pctx := peer.NewContext(nctx, &peer.Peer{
					Addr:     &net.IPAddr{},
					AuthInfo: credentials.TLSInfo{},
				})

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Logf(
						gomock.Eq("info"),
						gomock.Eq("request: %s@%s"),
						gomock.Eq("pkg.Service"),
						gomock.Eq("MethodName"),
					).
					Times(1)

				interceptor = grpclog.UnaryServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
				)
				res, err := interceptor(pctx, req, info, handler)

				Expect(res).To(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("metadata info available", func() {
			It("should return result", func() {
				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				nctx := context.Background()
				mctx := metadata.NewIncomingContext(nctx, metadata.New(map[string]string{
					"X-Correlation-Id": "123",
				}))

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Logf(
						gomock.Eq("info"),
						gomock.Eq("request: %s@%s"),
						gomock.Eq("pkg.Service"),
						gomock.Eq("MethodName"),
					).
					Times(1)

				interceptor = grpclog.UnaryServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
					grpclog.AllowedMetadata([]string{
						"X-Correlation-Id",
					}),
				)
				res, err := interceptor(mctx, req, info, handler)

				Expect(res).To(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("request is available", func() {
			It("should return result", func() {
				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				ctx.
					EXPECT().
					Deadline().
					Return(currentTs, false).
					Times(1)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				ctx.
					EXPECT().
					Value(gomock.Any()).
					Return(nil).
					Times(2)

				logger.
					EXPECT().
					Logf(
						gomock.Eq("info"),
						gomock.Eq("request: %s@%s"),
						gomock.Eq("pkg.Service"),
						gomock.Eq("MethodName"),
					).
					Times(1)

				interceptor = grpclog.UnaryServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
				)
				req := &testdata.SimpleMessage{
					String_: "key",
					Bool:    true,
				}
				res, err := interceptor(ctx, req, info, handler)

				Expect(res).To(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("response is available", func() {
			It("should return result", func() {
				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				ctx.
					EXPECT().
					Deadline().
					Return(currentTs, false).
					Times(1)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				ctx.
					EXPECT().
					Value(gomock.Any()).
					Return(nil).
					Times(2)

				logger.
					EXPECT().
					Logf(
						gomock.Eq("info"),
						gomock.Eq("request: %s@%s"),
						gomock.Eq("pkg.Service"),
						gomock.Eq("MethodName"),
					).
					Times(1)

				interceptor = grpclog.UnaryServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
				)

				res, err := interceptor(ctx, req, info, func(ctx context.Context, req interface{}) (interface{}, error) {
					res := &testdata.SimpleMessage{
						String_: "key",
						Bool:    true,
					}
					return res, nil
				})

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("StreamServerInterceptor function", Label("unit"), func() {
		var (
			currentTs   time.Time
			ctx         *mock_context.MockContext
			srv         interface{}
			ss          *mock_grpc.MockServerStream
			info        *grpc.StreamServerInfo
			logger      *mock_logging.MockLogger
			interceptor grpc.StreamServerInterceptor
			handler     func(srv interface{}, stream grpc.ServerStream) error
			clock       *mock_datetime.MockClock
		)

		BeforeEach(func() {
			t := GinkgoT()
			currentTs = time.Now()
			srv = struct{}{}
			info = &grpc.StreamServerInfo{
				FullMethod: "/pkg.Service/MethodName",
			}
			ctrl := gomock.NewController(t)
			ss = mock_grpc.NewMockServerStream(ctrl)
			ctx = mock_context.NewMockContext(ctrl)
			logger = mock_logging.NewMockLogger(ctrl)
			clock = mock_datetime.NewMockClock(ctrl)
			handler = func(srv interface{}, stream grpc.ServerStream) error {
				return nil
			}
		})

		When("grpc method is ignored", func() {
			It("should return result", func() {
				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				ss.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				ctx.
					EXPECT().
					Deadline().
					Return(currentTs, false).
					Times(1)

				interceptor = grpclog.StreamServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
					grpclog.IgnoredMethod([]string{
						info.FullMethod,
					}),
				)
				err := interceptor(srv, ss, info, handler)

				Expect(err).To(BeNil())
			})
		})

		When("grpc method is not ignored", func() {
			It("should return result", func() {
				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				ss.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				ctx.
					EXPECT().
					Deadline().
					Return(currentTs, false).
					Times(1)

				ctx.
					EXPECT().
					Value(gomock.Any()).
					Return(nil).
					Times(2)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Logf(
						gomock.Eq("info"),
						gomock.Eq("request: %s@%s"),
						gomock.Eq("pkg.Service"),
						gomock.Eq("MethodName"),
					).
					Times(1)

				interceptor = grpclog.StreamServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
				)
				err := interceptor(srv, ss, info, handler)

				Expect(err).To(BeNil())
			})
		})

		When("deadline occured", func() {
			It("should return result", func() {
				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				ss.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				ctx.
					EXPECT().
					Deadline().
					Return(currentTs, true).
					Times(1)

				ctx.
					EXPECT().
					Value(gomock.Any()).
					Return(nil).
					Times(2)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Logf(
						gomock.Eq("info"),
						gomock.Eq("request: %s@%s"),
						gomock.Eq("pkg.Service"),
						gomock.Eq("MethodName"),
					).
					Times(1)

				interceptor = grpclog.StreamServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
				)
				err := interceptor(srv, ss, info, handler)

				Expect(err).To(BeNil())
			})
		})

		When("fatal error occured", func() {
			It("should return result", func() {
				expectErr := status.New(codes.Unknown, "fatal error").Err()

				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				ss.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				ctx.
					EXPECT().
					Deadline().
					Return(currentTs, false).
					Times(1)

				errLog := map[string]interface{}{
					"error": expectErr,
				}
				logger.
					EXPECT().
					WithFields(gomock.Eq(errLog)).
					Return(logger).
					Times(1)

				ctx.
					EXPECT().
					Value(gomock.Any()).
					Return(nil).
					Times(2)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Logf(
						gomock.Eq("error"),
						gomock.Eq("request: %s@%s"),
						gomock.Eq("pkg.Service"),
						gomock.Eq("MethodName"),
					).
					Times(1)

				interceptor = grpclog.StreamServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
				)
				err := interceptor(srv, ss, info, func(srv interface{}, stream grpc.ServerStream) error {
					return expectErr
				})

				Expect(err).To(Equal(expectErr))
			})
		})

		When("warning error occured", func() {
			It("should return result", func() {
				expectErr := status.New(codes.PermissionDenied, "fatal error").Err()

				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				ss.
					EXPECT().
					Context().
					Return(ctx).
					Times(1)

				ctx.
					EXPECT().
					Deadline().
					Return(currentTs, false).
					Times(1)

				errLog := map[string]interface{}{
					"error": expectErr,
				}
				logger.
					EXPECT().
					WithFields(gomock.Eq(errLog)).
					Return(logger).
					Times(1)

				ctx.
					EXPECT().
					Value(gomock.Any()).
					Return(nil).
					Times(2)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Logf(
						gomock.Eq("warn"),
						gomock.Eq("request: %s@%s"),
						gomock.Eq("pkg.Service"),
						gomock.Eq("MethodName"),
					).
					Times(1)

				interceptor = grpclog.StreamServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
				)
				err := interceptor(srv, ss, info, func(srv interface{}, stream grpc.ServerStream) error {
					return expectErr
				})

				Expect(err).To(Equal(expectErr))
			})
		})

		When("network info available", func() {
			It("should return result", func() {
				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				nctx := context.Background()
				pctx := peer.NewContext(nctx, &peer.Peer{
					Addr:     &net.IPAddr{},
					AuthInfo: credentials.TLSInfo{},
				})

				ss.
					EXPECT().
					Context().
					Return(pctx).
					Times(1)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Logf(
						gomock.Eq("info"),
						gomock.Eq("request: %s@%s"),
						gomock.Eq("pkg.Service"),
						gomock.Eq("MethodName"),
					).
					Times(1)

				interceptor = grpclog.StreamServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
				)
				err := interceptor(srv, ss, info, handler)

				Expect(err).To(BeNil())
			})
		})

		When("metadata info available", func() {
			It("should return result", func() {
				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				nctx := context.Background()
				mctx := metadata.NewIncomingContext(nctx, metadata.New(map[string]string{
					"X-Correlation-Id": "123",
				}))

				ss.
					EXPECT().
					Context().
					Return(mctx).
					Times(1)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Logf(
						gomock.Eq("info"),
						gomock.Eq("request: %s@%s"),
						gomock.Eq("pkg.Service"),
						gomock.Eq("MethodName"),
					).
					Times(1)

				interceptor = grpclog.StreamServerInterceptor(
					grpclog.WithLogger(logger),
					grpclog.WithClock(clock),
					grpclog.AllowedMetadata([]string{
						"X-Correlation-Id",
					}),
				)
				err := interceptor(srv, ss, info, handler)

				Expect(err).To(BeNil())
			})
		})
	})
})
