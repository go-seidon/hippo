package healthcheck_test

import (
	"testing"
	"time"

	"github.com/go-seidon/local/internal/healthcheck"
	mock_logging "github.com/go-seidon/local/internal/logging/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHealthCheck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HealthCheck Package")
}

var _ = Describe("Health Check Job", func() {

	Context("WithLogger function", Label("unit"), func() {
		var (
			logger *mock_logging.MockLogger
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			logger = mock_logging.NewMockLogger(ctrl)
		})

		When("logger is specified", func() {
			It("should append logger", func() {
				opt := healthcheck.WithLogger(logger)

				var option healthcheck.HealthCheckParam
				opt(&option)

				Expect(option.Logger).To(Equal(logger))
				Expect(option.Jobs).To(BeNil())
			})
		})
	})

	Context("AddJob function", Label("unit"), func() {
		var (
			job *healthcheck.HealthJob
		)

		BeforeEach(func() {
			job = &healthcheck.HealthJob{
				Name:     "mock-name",
				Interval: 1 * time.Second,
			}
		})

		When("jobs are empty", func() {
			It("should append job", func() {
				opt := healthcheck.AddJob(job)

				var option healthcheck.HealthCheckParam
				opt(&option)

				Expect(option.Logger).To(BeNil())
				Expect(len(option.Jobs)).To(Equal(1))
			})
		})

		When("jobs are not empty", func() {
			It("should append job", func() {
				opt := healthcheck.AddJob(job)

				var option healthcheck.HealthCheckParam
				option.Jobs = append(option.Jobs, job)
				opt(&option)

				Expect(option.Logger).To(BeNil())
				Expect(len(option.Jobs)).To(Equal(2))
			})
		})
	})

	Context("WithJobs function", Label("unit"), func() {
		var (
			job  *healthcheck.HealthJob
			jobs []*healthcheck.HealthJob
		)

		BeforeEach(func() {
			job = &healthcheck.HealthJob{
				Name:     "mock-name",
				Interval: 1 * time.Second,
			}
			jobs = []*healthcheck.HealthJob{job, job}
		})

		When("jobs are empty", func() {
			It("should add jobs", func() {
				opt := healthcheck.WithJobs(jobs)

				var option healthcheck.HealthCheckParam

				opt(&option)

				Expect(option.Logger).To(BeNil())
				Expect(len(option.Jobs)).To(Equal(2))
			})
		})

		When("jobs are not empty", func() {
			It("should replace jobs", func() {
				opt := healthcheck.WithJobs(jobs)

				var option healthcheck.HealthCheckParam
				option.Jobs = []*healthcheck.HealthJob{job}

				opt(&option)

				Expect(option.Logger).To(BeNil())
				Expect(len(option.Jobs)).To(Equal(2))
			})
		})

	})
})
