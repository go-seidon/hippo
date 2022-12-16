package grpcapp_test

import (
	"context"
	"fmt"

	"github.com/go-seidon/hippo/internal/auth"
	mock_auth "github.com/go-seidon/hippo/internal/auth/mock"
	"github.com/go-seidon/hippo/internal/grpcapp"
	"github.com/go-seidon/hippo/internal/grpcauth"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interceptor Package", func() {

	Context("BasicAuth function", Label("unit"), func() {
		var (
			ctx     context.Context
			ba      *mock_auth.MockBasicAuth
			ccParam auth.CheckCredentialParam
			ccRes   *auth.CheckCredentialResult
		)

		BeforeEach(func() {
			ctx = metadata.NewIncomingContext(context.Background(), metadata.MD{
				grpcauth.AuthKey: []string{grpcauth.BasicKey + " token"},
			})
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			ba = mock_auth.NewMockBasicAuth(ctrl)
			ccParam = auth.CheckCredentialParam{
				AuthToken: "token",
			}
			ccRes = &auth.CheckCredentialResult{
				TokenValid: true,
			}
		})

		When("auth is not specified", func() {
			It("should return error", func() {
				ccErr := fmt.Errorf("request unauthenticated with basic")

				cc := grpcapp.BasicAuth(ba)

				ctx := context.Background()
				err := cc(ctx)

				expectErr := status.Errorf(codes.Unauthenticated, ccErr.Error())
				Expect(err).To(Equal(expectErr))
			})
		})

		When("failed check credential", func() {
			It("should return error", func() {
				ccErr := fmt.Errorf("network error")

				ba.
					EXPECT().
					CheckCredential(gomock.Eq(ctx), gomock.Eq(ccParam)).
					Return(nil, ccErr).
					Times(1)

				cc := grpcapp.BasicAuth(ba)

				err := cc(ctx)

				expectErr := status.Errorf(codes.Unknown, ccErr.Error())
				Expect(err).To(Equal(expectErr))
			})
		})

		When("credential is not valid", func() {
			It("should return error", func() {
				ccRes := &auth.CheckCredentialResult{
					TokenValid: false,
				}
				ba.
					EXPECT().
					CheckCredential(gomock.Eq(ctx), gomock.Eq(ccParam)).
					Return(ccRes, nil).
					Times(1)

				cc := grpcapp.BasicAuth(ba)

				err := cc(ctx)

				expectErr := status.Errorf(codes.Unauthenticated, grpcauth.ErrorInvalidCredential.Error())
				Expect(err).To(Equal(expectErr))
			})
		})

		When("credential is valid", func() {
			It("should return result", func() {
				ba.
					EXPECT().
					CheckCredential(gomock.Eq(ctx), gomock.Eq(ccParam)).
					Return(ccRes, nil).
					Times(1)

				cc := grpcapp.BasicAuth(ba)

				err := cc(ctx)

				Expect(err).To(BeNil())
			})
		})
	})

})
