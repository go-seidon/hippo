package healthcheck_test

import (
	"context"
	"fmt"
	"time"

	"github.com/InVisionApp/go-health"
	"github.com/go-seidon/hippo/internal/healthcheck"
	mock_healthcheck "github.com/go-seidon/hippo/internal/healthcheck/mock"
	mock_logging "github.com/go-seidon/hippo/internal/logging/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Go Health Check", func() {

	Context("NewGoHealthCheck function", Label("unit"), func() {
		var (
			jobs   []*healthcheck.HealthJob
			logger *mock_logging.MockLogger
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			logger = mock_logging.NewMockLogger(ctrl)
			jobs = []*healthcheck.HealthJob{
				{
					Name:     "mock-job",
					Checker:  nil,
					Interval: 1,
				},
			}
		})

		When("jobs are not specified", func() {
			It("should return error", func() {
				r, err := healthcheck.NewGoHealthCheck(healthcheck.WithJobs(nil))

				Expect(r).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid jobs specified")))
			})
		})

		When("jobs are empty", func() {
			It("should return error", func() {
				jobs := []*healthcheck.HealthJob{}
				r, err := healthcheck.NewGoHealthCheck(
					healthcheck.WithJobs(jobs),
					healthcheck.WithLogger(logger),
				)

				Expect(r).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid jobs specified")))
			})
		})

		When("logger is not specified", func() {
			It("should return error", func() {
				r, err := healthcheck.NewGoHealthCheck(
					healthcheck.WithJobs(jobs),
				)

				Expect(r).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid logger specified")))
			})
		})

		When("all params are specified", func() {
			It("should return result", func() {
				r, err := healthcheck.NewGoHealthCheck(
					healthcheck.WithJobs(jobs),
					healthcheck.WithLogger(logger),
				)

				Expect(r).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("Start function", Label("unit"), func() {
		var (
			ctx    context.Context
			client *mock_healthcheck.MockHealthClient
			s      healthcheck.HealthCheck
			logger *mock_logging.MockLogger
		)

		BeforeEach(func() {
			ctx = context.Background()
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			client = mock_healthcheck.NewMockHealthClient(ctrl)
			jobs := []*healthcheck.HealthJob{
				{
					Name:     "mock-job",
					Checker:  nil,
					Interval: 1,
				},
			}
			logger = mock_logging.NewMockLogger(ctrl)
			s, _ = healthcheck.NewGoHealthCheck(
				healthcheck.WithJobs(jobs),
				healthcheck.WithLogger(logger),
				healthcheck.WithClient(client),
			)
		})

		When("failed add checkers", func() {
			It("should return error", func() {
				client.
					EXPECT().
					AddChecks(gomock.Any()).
					Return(fmt.Errorf("failed add checkers")).
					Times(1)

				err := s.Start(ctx)

				Expect(err).To(Equal(fmt.Errorf("failed add checkers")))
			})
		})

		When("failed start app", func() {
			It("should return error", func() {
				client.
					EXPECT().
					AddChecks(gomock.Any()).
					Return(nil).
					Times(1)

				client.
					EXPECT().
					Start().
					Return(fmt.Errorf("failed start app")).
					Times(1)

				err := s.Start(ctx)

				Expect(err).To(Equal(fmt.Errorf("failed start app")))
			})
		})

		When("success start app", func() {
			It("should return result", func() {
				client.
					EXPECT().
					AddChecks(gomock.Any()).
					Return(nil).
					Times(1)

				client.
					EXPECT().
					Start().
					Return(nil).
					Times(1)

				err := s.Start(ctx)

				Expect(err).To(BeNil())
			})
		})

		When("app is already started", func() {
			It("should return result", func() {
				client.
					EXPECT().
					AddChecks(gomock.Any()).
					Return(nil).
					Times(1)

				client.
					EXPECT().
					Start().
					Return(nil).
					Times(1)

				err1 := s.Start(ctx)
				err2 := s.Start(ctx)

				Expect(err1).To(BeNil())
				Expect(err2).To(BeNil())
			})
		})
	})

	Context("Stop function", Label("unit"), func() {
		var (
			ctx    context.Context
			client *mock_healthcheck.MockHealthClient
			s      healthcheck.HealthCheck
			logger *mock_logging.MockLogger
		)

		BeforeEach(func() {
			ctx = context.Background()
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			client = mock_healthcheck.NewMockHealthClient(ctrl)
			jobs := []*healthcheck.HealthJob{
				{
					Name:     "mock-job",
					Checker:  nil,
					Interval: 1,
				},
			}
			logger = mock_logging.NewMockLogger(ctrl)
			s, _ = healthcheck.NewGoHealthCheck(
				healthcheck.WithJobs(jobs),
				healthcheck.WithLogger(logger),
				healthcheck.WithClient(client),
			)
		})

		When("failed stop app", func() {
			It("should return error", func() {
				client.
					EXPECT().
					Stop().
					Return(fmt.Errorf("failed stop app")).
					Times(1)

				err := s.Stop(ctx)

				Expect(err).To(Equal(fmt.Errorf("failed stop app")))
			})
		})

		When("success stop app", func() {
			It("should return result", func() {
				client.
					EXPECT().
					Stop().
					Return(nil).
					Times(1)

				err := s.Stop(ctx)

				Expect(err).To(BeNil())
			})
		})
	})

	Context("Check function", Label("unit"), func() {
		var (
			ctx              context.Context
			client           *mock_healthcheck.MockHealthClient
			s                healthcheck.HealthCheck
			currentTimestamp time.Time
			logger           *mock_logging.MockLogger
		)

		BeforeEach(func() {
			ctx = context.Background()
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			client = mock_healthcheck.NewMockHealthClient(ctrl)
			jobs := []*healthcheck.HealthJob{
				{
					Name:     "mock-job",
					Checker:  nil,
					Interval: 1,
				},
			}
			logger = mock_logging.NewMockLogger(ctrl)
			s, _ = healthcheck.NewGoHealthCheck(
				healthcheck.WithJobs(jobs),
				healthcheck.WithLogger(logger),
				healthcheck.WithClient(client),
			)
			currentTimestamp = time.Now()
		})

		When("error occured", func() {
			It("should return error", func() {
				client.
					EXPECT().
					State().
					Return(nil, true, fmt.Errorf("network error")).
					Times(1)

				res, err := s.Check(ctx)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("failed get state", func() {
			It("should return result", func() {
				client.
					EXPECT().
					State().
					Return(nil, true, nil).
					Times(1)

				res, err := s.Check(ctx)

				expected := &healthcheck.CheckResult{
					Status: "FAILED",
					Items:  make(map[string]healthcheck.CheckResultItem),
				}
				Expect(res).To(Equal(expected))
				Expect(err).To(BeNil())
			})
		})

		When("all check is ok", func() {
			It("should return result", func() {
				states := map[string]health.State{
					"mock-job": {
						Name:      "mock-job",
						Status:    "ok",
						Err:       "",
						Fatal:     false,
						Details:   nil,
						CheckTime: currentTimestamp,
					},
				}

				client.
					EXPECT().
					State().
					Return(states, false, nil).
					Times(1)

				res, err := s.Check(ctx)

				expected := &healthcheck.CheckResult{
					Status: "OK",
					Items: map[string]healthcheck.CheckResultItem{
						"mock-job": {
							Name:      "mock-job",
							Status:    "OK",
							Error:     "",
							Fatal:     false,
							CheckedAt: currentTimestamp.UTC(),
						},
					},
				}
				Expect(res).To(Equal(expected))
				Expect(err).To(BeNil())
			})
		})

		When("all check is failed", func() {
			It("should return result", func() {
				states := map[string]health.State{
					"mock-job": {
						Name:      "mock-job",
						Status:    "failed",
						Err:       "some error",
						Fatal:     false,
						Details:   nil,
						CheckTime: currentTimestamp,
					},
				}

				client.
					EXPECT().
					State().
					Return(states, false, nil).
					Times(1)

				res, err := s.Check(ctx)

				expected := &healthcheck.CheckResult{
					Status: "FAILED",
					Items: map[string]healthcheck.CheckResultItem{
						"mock-job": {
							Name:      "mock-job",
							Status:    "FAILED",
							Error:     "some error",
							Fatal:     false,
							CheckedAt: currentTimestamp.UTC(),
						},
					},
				}
				Expect(res).To(Equal(expected))
				Expect(err).To(BeNil())
			})
		})

		When("some check is failed", func() {
			It("should return result", func() {
				states := map[string]health.State{
					"mock-job": {
						Name:      "mock-job",
						Status:    "failed",
						Err:       "some error",
						Fatal:     false,
						Details:   nil,
						CheckTime: currentTimestamp,
					},
					"mock-job-2": {
						Name:      "mock-job-2",
						Status:    "ok",
						Err:       "",
						Fatal:     false,
						Details:   nil,
						CheckTime: currentTimestamp,
					},
				}

				client.
					EXPECT().
					State().
					Return(states, false, nil).
					Times(1)

				res, err := s.Check(ctx)

				expected := &healthcheck.CheckResult{
					Status: "WARNING",
					Items: map[string]healthcheck.CheckResultItem{
						"mock-job": {
							Name:      "mock-job",
							Status:    "FAILED",
							Error:     "some error",
							Fatal:     false,
							CheckedAt: currentTimestamp.UTC(),
						},
						"mock-job-2": {
							Name:      "mock-job-2",
							Status:    "OK",
							Error:     "",
							Fatal:     false,
							CheckedAt: currentTimestamp.UTC(),
						},
					},
				}
				Expect(res).To(Equal(expected))
				Expect(err).To(BeNil())
			})
		})
	})

})
