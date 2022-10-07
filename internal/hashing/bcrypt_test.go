package hashing_test

import (
	"fmt"

	"github.com/go-seidon/hippo/internal/hashing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bcrypt Hasher Package", func() {
	Context("NewBcryptHasher function", Label("unit"), func() {
		When("function is called", func() {
			It("should return result", func() {
				res := hashing.NewBcryptHasher()

				Expect(res).ToNot(BeNil())
			})
		})
	})

	Context("Generate function", Label("unit", "slow"), func() {
		var (
			h    hashing.Hasher
			text string
		)

		BeforeEach(func() {
			h = hashing.NewBcryptHasher()
			text = "some-secret"
		})

		When("success generate hash", func() {
			It("should return result", func() {
				res, err := h.Generate(text)

				equalIfNil := h.Verify(string(res), text)

				Expect(res).ToNot(BeEmpty())
				Expect(err).To(BeNil())
				Expect(equalIfNil).To(BeNil())
			})
		})
	})

	Context("Verify function", Label("unit", "slow"), func() {
		var (
			h    hashing.Hasher
			text string
			hash string
		)

		BeforeEach(func() {
			h = hashing.NewBcryptHasher()
			text = "some-secret"
			hash = "$2a$10$xA9.FPfIYi2ZI6V5/jw5leFVUCjsgN4lBS5iS8loLv1hngJj1ys/2"
		})

		When("hash is equal", func() {
			It("should return nil", func() {
				err := h.Verify(hash, text)

				Expect(err).To(BeNil())
			})
		})

		When("hash is not equal", func() {
			It("should return nil", func() {
				err := h.Verify(hash, "other-secret")

				Expect(err).To(Equal(fmt.Errorf("crypto/bcrypt: hashedPassword is not the hash of the given password")))
			})
		})
	})

})
