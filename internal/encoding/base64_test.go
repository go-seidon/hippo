package encoding_test

import (
	"encoding/base64"

	"github.com/go-seidon/hippo/internal/encoding"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Base64 Encoder Package", func() {
	Context("NewBase64Encoder function", Label("unit"), func() {
		When("function is called", func() {
			It("should return result", func() {
				res := encoding.NewBase64Encoder()

				Expect(res).ToNot(BeNil())
			})
		})
	})

	Context("Encode function", Label("unit"), func() {
		var (
			encoder encoding.Encoder
			d       []byte
		)

		BeforeEach(func() {
			encoder = encoding.NewBase64Encoder()
			d = []byte("some-data")
		})

		When("success encode data", func() {
			It("should return result", func() {
				res, err := encoder.Encode(d)

				Expect(res).To(Equal("c29tZS1kYXRh"))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("Decode function", Label("unit"), func() {
		var (
			encoder encoding.Encoder
			d       string
		)

		BeforeEach(func() {
			encoder = encoding.NewBase64Encoder()
			d = "c29tZS1kYXRh"
		})

		When("failed decode data", func() {
			It("should return error", func() {
				res, err := encoder.Decode("\\")

				Expect(res).To(Equal([]byte{}))
				Expect(err).To(Equal(base64.CorruptInputError(0)))
			})
		})

		When("success decode data", func() {
			It("should return result", func() {
				res, err := encoder.Decode(d)

				Expect(string(res)).To(Equal("some-data"))
				Expect(err).To(BeNil())
			})
		})
	})

})
