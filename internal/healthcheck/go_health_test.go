package healthcheck_test

import (
	"fmt"
	"time"

	"github.com/InVisionApp/go-health"
	"github.com/go-seidon/local/internal/healthcheck"
	"github.com/go-seidon/local/internal/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Go Health Check", func() {

	Context("NewGoHealthCheck function", Label("unit"), func() {
		var (
			jobs   []*healthcheck.HealthJob
			logger *mock.MockLogger
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			logger = mock.NewMockLogger(ctrl)
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
				r, err := healthcheck.NewGoHealthCheck(nil, logger)

				Expect(r).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid jobs specified")))
			})
		})

		When("jobs are empty", func() {
			It("should return error", func() {
				jobs := []*healthcheck.HealthJob{}
				r, err := healthcheck.NewGoHealthCheck(jobs, logger)

				Expect(r).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid jobs specified")))
			})
		})

		When("logger is not specified", func() {
			It("should return error", func() {
				r, err := healthcheck.NewGoHealthCheck(jobs, nil)

				Expect(r).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid logger specified")))
			})
		})

		When("all params are specified", func() {
			It("should return result", func() {
				r, err := healthcheck.NewGoHealthCheck(jobs, logger)

				Expect(r).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("Start function", Label("unit"), func() {
		var (
			client *mock.MockHealthClient
			s      *healthcheck.GoHealthCheck
			logger *mock.MockLogger
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			client = mock.NewMockHealthClient(ctrl)
			jobs := []*healthcheck.HealthJob{
				{
					Name:     "mock-job",
					Checker:  nil,
					Interval: 1,
				},
			}
			logger = mock.NewMockLogger(ctrl)
			s, _ = healthcheck.NewGoHealthCheck(jobs, logger)
			s.Client = client
		})

		When("failed add checkers", func() {
			It("should return error", func() {
				client.
					EXPECT().
					AddChecks(gomock.Any()).
					Return(fmt.Errorf("failed add checkers")).
					Times(1)

				err := s.Start()

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

				err := s.Start()

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

				err := s.Start()

				Expect(err).To(BeNil())
			})
		})
	})

	Context("Stop function", Label("unit"), func() {
		var (
			client *mock.MockHealthClient
			s      *healthcheck.GoHealthCheck
			logger *mock.MockLogger
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			client = mock.NewMockHealthClient(ctrl)
			jobs := []*healthcheck.HealthJob{
				{
					Name:     "mock-job",
					Checker:  nil,
					Interval: 1,
				},
			}
			logger = mock.NewMockLogger(ctrl)
			s, _ = healthcheck.NewGoHealthCheck(jobs, logger)
			s.Client = client
		})

		When("failed stop app", func() {
			It("should return error", func() {
				client.
					EXPECT().
					Stop().
					Return(fmt.Errorf("failed stop app")).
					Times(1)

				err := s.Stop()

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

				err := s.Stop()

				Expect(err).To(BeNil())
			})
		})
	})

	Context("Check function", Label("unit"), func() {
		var (
			client           *mock.MockHealthClient
			s                *healthcheck.GoHealthCheck
			currentTimestamp time.Time
			logger           *mock.MockLogger
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			client = mock.NewMockHealthClient(ctrl)
			jobs := []*healthcheck.HealthJob{
				{
					Name:     "mock-job",
					Checker:  nil,
					Interval: 1,
				},
			}
			logger = mock.NewMockLogger(ctrl)
			s, _ = healthcheck.NewGoHealthCheck(jobs, logger)
			s.Client = client
			currentTimestamp = time.Now()
		})

		When("error occured", func() {
			It("should return error", func() {
				client.
					EXPECT().
					State().
					Return(nil, true, fmt.Errorf("network error")).
					Times(1)

				res, err := s.Check()

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

				res, err := s.Check()

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
					"mock-job": health.State{
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

				res, err := s.Check()

				expected := &healthcheck.CheckResult{
					Status: "OK",
					Items: map[string]healthcheck.CheckResultItem{
						"mock-job": healthcheck.CheckResultItem{
							Name:      "mock-job",
							Status:    "OK",
							Error:     "",
							Fatal:     false,
							Metadata:  nil,
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
					"mock-job": health.State{
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

				res, err := s.Check()

				expected := &healthcheck.CheckResult{
					Status: "FAILED",
					Items: map[string]healthcheck.CheckResultItem{
						"mock-job": healthcheck.CheckResultItem{
							Name:      "mock-job",
							Status:    "FAILED",
							Error:     "some error",
							Fatal:     false,
							Metadata:  nil,
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

				res, err := s.Check()

				expected := &healthcheck.CheckResult{
					Status: "WARNING",
					Items: map[string]healthcheck.CheckResultItem{
						"mock-job": {
							Name:      "mock-job",
							Status:    "FAILED",
							Error:     "some error",
							Fatal:     false,
							Metadata:  nil,
							CheckedAt: currentTimestamp.UTC(),
						},
						"mock-job-2": {
							Name:      "mock-job-2",
							Status:    "OK",
							Error:     "",
							Fatal:     false,
							Metadata:  nil,
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
