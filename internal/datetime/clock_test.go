package datetime_test

import (
	"github.com/go-seidon/hippo/internal/datetime"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Datetime Package", func() {

	Context("NewClock function", Label("unit"), func() {
		When("success create clock", func() {
			It("should return result", func() {
				res := datetime.NewClock()

				Expect(res).ToNot(BeNil())
			})
		})
	})

	Context("Now function", Label("unit"), func() {
		var (
			c datetime.Clock
		)

		BeforeEach(func() {
			c = datetime.NewClock()
		})

		When("success get current time", func() {
			It("should return result", func() {
				res := c.Now()

				Expect(res).ToNot(BeNil())
			})
		})
	})

})
