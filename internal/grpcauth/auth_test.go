package grpcauth_test

import (
	"context"
	"fmt"
	"testing"

	grpc_auth "github.com/go-seidon/hippo/internal/grpcauth"
	"google.golang.org/grpc/metadata"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGrpcAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Grpc Auth Package")
}

var _ = Describe("Auth Package", func() {

	Context("AuthFromMD function", Label("unit"), func() {
		var (
			ctx    context.Context
			scheme string
		)

		BeforeEach(func() {
			scheme = grpc_auth.BasicKey
			ctx = metadata.NewIncomingContext(context.Background(), metadata.MD{
				grpc_auth.AuthKey: []string{scheme + " token"},
			})
		})

		When("auth is not available", func() {
			It("should return error", func() {
				ctx := context.Background()
				res, err := grpc_auth.AuthFromMD(ctx, scheme)

				Expect(res).To(Equal(""))
				Expect(err).To(Equal(fmt.Errorf("request unauthenticated with %s", scheme)))
			})
		})

		When("auth string is not available", func() {
			It("should return error", func() {
				ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{
					grpc_auth.AuthKey: []string{scheme + ""},
				})
				res, err := grpc_auth.AuthFromMD(ctx, scheme)

				Expect(res).To(Equal(""))
				Expect(err).To(Equal(fmt.Errorf("bad authorization string")))
			})
		})

		When("scheme is not match", func() {
			It("should return error", func() {
				scheme := "invalid"
				res, err := grpc_auth.AuthFromMD(ctx, scheme)

				Expect(res).To(Equal(""))
				Expect(err).To(Equal(fmt.Errorf("invalid scheme of %s", scheme)))
			})
		})

		When("scheme is match", func() {
			It("should return result", func() {
				res, err := grpc_auth.AuthFromMD(ctx, scheme)

				Expect(res).To(Equal("token"))
				Expect(err).To(BeNil())
			})
		})
	})
})
