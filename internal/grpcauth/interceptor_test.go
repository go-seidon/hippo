package grpcauth_test

import (
	"context"

	"github.com/go-seidon/hippo/internal/grpcauth"
	mock_grpc "github.com/go-seidon/provider/grpc/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

var _ = Describe("Auth Package", func() {

	Context("UnaryServerInterceptor function", Label("unit"), func() {
		var (
			ctx     context.Context
			req     interface{}
			info    *grpc.UnaryServerInfo
			handler func(ctx context.Context, req interface{}) (interface{}, error)
		)

		BeforeEach(func() {
			ctx = context.Background()
			req = struct{}{}
			info = &grpc.UnaryServerInfo{}
			handler = func(ctx context.Context, req interface{}) (interface{}, error) {
				res := struct{}{}
				return res, nil
			}
		})

		When("credential is not valid", func() {
			It("should return error", func() {
				expectErr := grpcauth.ErrorInvalidCredential
				cc := func(ctx context.Context) error {
					return expectErr
				}
				interceptor := grpcauth.UnaryServerInterceptor(
					grpcauth.WithAuth(cc),
				)

				res, err := interceptor(ctx, req, info, handler)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(expectErr))
			})
		})

		When("credential is valid", func() {
			It("should return result", func() {
				cc := func(ctx context.Context) error {
					return nil
				}
				interceptor := grpcauth.UnaryServerInterceptor(
					grpcauth.WithAuth(cc),
				)

				res, err := interceptor(ctx, req, info, handler)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("StreamServerInterceptor function", Label("unit"), func() {
		var (
			srv     interface{}
			ss      *mock_grpc.MockServerStream
			info    *grpc.StreamServerInfo
			handler func(srv interface{}, stream grpc.ServerStream) error
		)

		BeforeEach(func() {
			srv = struct{}{}
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			ss = mock_grpc.NewMockServerStream(ctrl)
			info = &grpc.StreamServerInfo{}
			handler = func(srv interface{}, stream grpc.ServerStream) error {
				return nil
			}
			ss.
				EXPECT().
				Context().
				Return(context.Background()).
				Times(1)
		})

		When("credential is not valid", func() {
			It("should return error", func() {
				expectErr := grpcauth.ErrorInvalidCredential
				cc := func(ctx context.Context) error {
					return expectErr
				}
				interceptor := grpcauth.StreamServerInterceptor(
					grpcauth.WithAuth(cc),
				)

				err := interceptor(srv, ss, info, handler)

				Expect(err).To(Equal(expectErr))
			})
		})

		When("credential is valid", func() {
			It("should return result", func() {
				cc := func(ctx context.Context) error {
					return nil
				}
				interceptor := grpcauth.StreamServerInterceptor(
					grpcauth.WithAuth(cc),
				)

				err := interceptor(srv, ss, info, handler)

				Expect(err).To(BeNil())
			})
		})
	})

})
