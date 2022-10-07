package grpc_auth_test

import (
	"context"

	grpc_auth "github.com/go-seidon/hippo/internal/grpc-auth"
	mock_grpc "github.com/go-seidon/hippo/internal/grpc/mock"
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
				expectErr := grpc_auth.ErrorInvalidCredential
				cc := func(ctx context.Context) error {
					return expectErr
				}
				interceptor := grpc_auth.UnaryServerInterceptor(
					grpc_auth.WithAuth(cc),
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
				interceptor := grpc_auth.UnaryServerInterceptor(
					grpc_auth.WithAuth(cc),
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
				expectErr := grpc_auth.ErrorInvalidCredential
				cc := func(ctx context.Context) error {
					return expectErr
				}
				interceptor := grpc_auth.StreamServerInterceptor(
					grpc_auth.WithAuth(cc),
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
				interceptor := grpc_auth.StreamServerInterceptor(
					grpc_auth.WithAuth(cc),
				)

				err := interceptor(srv, ss, info, handler)

				Expect(err).To(BeNil())
			})
		})
	})

})
