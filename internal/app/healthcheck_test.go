package app_test

import (
	"github.com/go-seidon/local/internal/app"
	"github.com/go-seidon/local/internal/healthcheck"
	mock_logging "github.com/go-seidon/local/internal/logging/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Healthcheck Package", func() {

	Context("NewDefaultHealthCheck function", func() {
		var (
			logger *mock_logging.MockLogger
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			logger = mock_logging.NewMockLogger(ctrl)
		})

		When("success create default healthcheck", func() {
			It("should return result", func() {
				res, err := app.NewDefaultHealthCheck(
					healthcheck.WithLogger(logger),
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("failed create default healthcheck", func() {
			It("should return error", func() {
				res, err := app.NewDefaultHealthCheck()

				Expect(res).To(BeNil())
				Expect(err).ToNot(BeNil())
			})
		})
	})

})
