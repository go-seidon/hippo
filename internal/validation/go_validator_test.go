package validation_test

import (
	"github.com/go-seidon/local/internal/validation"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Go Validator", func() {

	Context("Validate function", Label("unit"), func() {
		var (
			validator validation.Validator
		)

		BeforeEach(func() {
			validator = validation.NewGoValidator()
		})

		When("there is no invalid data", func() {
			It("should return error", func() {
				i := struct {
					Key string `validate:"required,min=3"`
				}{
					Key: "key",
				}

				err := validator.Validate(i)

				Expect(err).To(BeNil())
			})
		})

		When("data is not a struct", func() {
			It("should return error", func() {
				i := make(chan bool)

				err := validator.Validate(i)

				expectErr := validation.Error(
					"validator: (nil chan bool)",
				)
				Expect(err.Error()).To(Equal(expectErr.Error()))
			})
		})

		When("there are invalid data", func() {
			It("should return error", func() {
				i := struct {
					Key string `validate:"required,min=3"`
				}{
					Key: "",
				}

				err := validator.Validate(i)

				expectErr := validation.Error(
					"Key: 'Key' Error:Field validation for 'Key' failed on the 'required' tag",
				)
				Expect(err.Error()).To(Equal(expectErr.Error()))
			})
		})

	})

})
