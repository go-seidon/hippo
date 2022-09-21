package healthcheck_test

import (
	"fmt"

	"github.com/go-seidon/local/internal/healthcheck"
	mock_healthcheck "github.com/go-seidon/local/internal/healthcheck/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Health Check Checker", func() {

	Context("RepoPing Checker", Label("unit"), func() {
		var (
			checker    healthcheck.Checker
			dataSource *mock_healthcheck.MockDataSource
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			dataSource = mock_healthcheck.NewMockDataSource(ctrl)
			checker = healthcheck.NewRepoPingChecker(dataSource)
		})

		When("failed ping repository", func() {
			It("should return error", func() {
				dataSource.
					EXPECT().
					Ping(gomock.Any()).
					Return(fmt.Errorf("db error")).
					Times(1)

				res, err := checker.Status()

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("db error")))
			})
		})

		When("success ping repository", func() {
			It("should return result", func() {
				dataSource.
					EXPECT().
					Ping(gomock.Any()).
					Return(nil).
					Times(1)

				res, err := checker.Status()

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

})
