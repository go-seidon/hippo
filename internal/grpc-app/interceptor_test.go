package grpc_app_test

import (
	"context"
	"fmt"

	"github.com/go-seidon/hippo/internal/auth"
	mock_auth "github.com/go-seidon/hippo/internal/auth/mock"
	grpc_app "github.com/go-seidon/hippo/internal/grpc-app"
	grpc_auth "github.com/go-seidon/hippo/internal/grpc-auth"
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
				grpc_auth.AuthKey: []string{grpc_auth.BasicKey + " token"},
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

				cc := grpc_app.BasicAuth(ba)

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

				cc := grpc_app.BasicAuth(ba)

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

				cc := grpc_app.BasicAuth(ba)

				err := cc(ctx)

				expectErr := status.Errorf(codes.Unauthenticated, grpc_auth.ErrorInvalidCredential.Error())
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

				cc := grpc_app.BasicAuth(ba)

				err := cc(ctx)

				Expect(err).To(BeNil())
			})
		})
	})

})
