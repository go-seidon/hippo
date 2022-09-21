package app_test

import (
	"github.com/go-seidon/local/internal/app"
	mock_logging "github.com/go-seidon/local/internal/logging/mock"
	mock_repository "github.com/go-seidon/local/internal/repository/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Healthcheck Package", func() {

	Context("NewDefaultHealthCheck function", Label("unit"), func() {
		var (
			logger     *mock_logging.MockLogger
			repository *mock_repository.MockProvider
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			logger = mock_logging.NewMockLogger(ctrl)
			repository = mock_repository.NewMockProvider(ctrl)
		})

		When("success create default healthcheck", func() {
			It("should return result", func() {
				res, err := app.NewDefaultHealthCheck(logger, repository)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("failed create default healthcheck", func() {
			It("should return error", func() {
				res, err := app.NewDefaultHealthCheck(nil, nil)

				Expect(res).To(BeNil())
				Expect(err).ToNot(BeNil())
			})
		})
	})

})
