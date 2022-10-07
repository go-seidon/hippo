package grpc_meta_test

import (
	"context"
	"testing"

	grpc_meta "github.com/go-seidon/hippo/internal/grpc-meta"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"
)

func TestGrpcMeta(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Grpc Meta Package")
}

var _ = Describe("Metadata Package", func() {

	Describe("Metadata type", func() {
		Context("Get function", Label("unit"), func() {
			When("key exists", func() {
				It("should return result", func() {
					md := grpc_meta.Metadata{
						"key": []string{"value"},
					}
					res := md.Get("key")

					Expect(res).To(Equal("value"))
				})
			})

			When("key is not exists", func() {
				It("should return empty", func() {
					md := grpc_meta.Metadata{}
					res := md.Get("key")

					Expect(res).To(Equal(""))
				})
			})

			When("key exists with empty value", func() {
				It("should return empty", func() {
					md := grpc_meta.Metadata{
						"key": []string{},
					}
					res := md.Get("key")

					Expect(res).To(Equal(""))
				})
			})
		})
	})

	Context("ExtractIncoming function", Label("unit"), func() {
		When("metadata are available", func() {
			It("should return result", func() {
				md := metadata.MD{
					"key": []string{"value"},
				}
				ctx := metadata.NewIncomingContext(context.Background(), md)

				res := grpc_meta.ExtractIncoming(ctx)

				expectRes := grpc_meta.Metadata{
					"key": []string{"value"},
				}
				Expect(res).To(Equal(expectRes))
			})
		})

		When("metadata are not available", func() {
			It("should return result", func() {
				md := grpc_meta.Metadata(metadata.Pairs())
				ctx := context.Background()

				res := grpc_meta.ExtractIncoming(ctx)

				Expect(res).To(Equal(md))
			})
		})
	})

	Context("ExtractOutgoing function", Label("unit"), func() {
		When("metadata are available", func() {
			It("should return result", func() {
				md := metadata.MD{
					"key": []string{"value"},
				}
				ctx := metadata.NewOutgoingContext(context.Background(), md)

				res := grpc_meta.ExtractOutgoing(ctx)

				expectRes := grpc_meta.Metadata{
					"key": []string{"value"},
				}
				Expect(res).To(Equal(expectRes))
			})
		})

		When("metadata are not available", func() {
			It("should return result", func() {
				md := grpc_meta.Metadata(metadata.Pairs())
				ctx := context.Background()

				res := grpc_meta.ExtractOutgoing(ctx)

				Expect(res).To(Equal(md))
			})
		})
	})

})
