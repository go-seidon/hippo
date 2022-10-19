package healthcheck_test

import (
	"fmt"

	"github.com/go-seidon/hippo/internal/healthcheck"
	mock_logging "github.com/go-seidon/provider/logging/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Health Check Log", func() {

	Context("NewGoHealthLog function", Label("unit"), func() {
		When("logger is not specified", func() {
			It("should return error", func() {
				res, err := healthcheck.NewGoHealthLog(nil)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid logger")))
			})
		})

		When("parameters are specified", func() {
			It("should return result", func() {
				t := GinkgoT()
				ctrl := gomock.NewController(t)
				logger := mock_logging.NewMockLogger(ctrl)
				res, err := healthcheck.NewGoHealthLog(logger)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("Go Health Log", Label("unit"), func() {
		var (
			log    *healthcheck.GoHealthLog
			client *mock_logging.MockLogger
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			client = mock_logging.NewMockLogger(ctrl)
			log, _ = healthcheck.NewGoHealthLog(client)
		})

		Context("Info function", Label("unit"), func() {
			When("success send log", func() {
				It("should log the message", func() {
					client.
						EXPECT().
						Info(gomock.Eq("mock-log-1"), gomock.Eq("mock-log-2")).
						Times(1)

					log.Info("mock-log-1", "mock-log-2")
				})
			})
		})

		Context("Debug function", Label("unit"), func() {
			When("success send log", func() {
				It("should log the message", func() {
					client.
						EXPECT().
						Debug(gomock.Eq("mock-log-1"), gomock.Eq("mock-log-2")).
						Times(1)

					log.Debug("mock-log-1", "mock-log-2")
				})
			})
		})

		Context("Error function", Label("unit"), func() {
			When("success send log", func() {
				It("should log the message", func() {
					client.
						EXPECT().
						Error(gomock.Eq("mock-log-1"), gomock.Eq("mock-log-2")).
						Times(1)

					log.Error("mock-log-1", "mock-log-2")
				})
			})
		})

		Context("Warn function", Label("unit"), func() {
			When("success send log", func() {
				It("should log the message", func() {
					client.
						EXPECT().
						Warn(gomock.Eq("mock-log-1"), gomock.Eq("mock-log-2")).
						Times(1)

					log.Warn("mock-log-1", "mock-log-2")
				})
			})
		})

		Context("Infof function", Label("unit"), func() {
			When("success send log", func() {
				It("should log the message", func() {
					client.
						EXPECT().
						Infof(
							gomock.Eq("%s = %s"),
							gomock.Eq("mock-log-1"),
							gomock.Eq("mock-log-2"),
						).
						Times(1)

					log.Infof("%s = %s", "mock-log-1", "mock-log-2")
				})
			})
		})

		Context("Debugf function", Label("unit"), func() {
			When("success send log", func() {
				It("should log the message", func() {
					client.
						EXPECT().
						Debugf(
							gomock.Eq("%s = %s"),
							gomock.Eq("mock-log-1"),
							gomock.Eq("mock-log-2"),
						).
						Times(1)

					log.Debugf("%s = %s", "mock-log-1", "mock-log-2")
				})
			})
		})

		Context("Errorf function", Label("unit"), func() {
			When("success send log", func() {
				It("should log the message", func() {
					client.
						EXPECT().
						Errorf(
							gomock.Eq("%s = %s"),
							gomock.Eq("mock-log-1"),
							gomock.Eq("mock-log-2"),
						).
						Times(1)

					log.Errorf("%s = %s", "mock-log-1", "mock-log-2")
				})
			})
		})

		Context("Warnf function", Label("unit"), func() {
			When("success send log", func() {
				It("should log the message", func() {
					client.
						EXPECT().
						Warnf(
							gomock.Eq("%s = %s"),
							gomock.Eq("mock-log-1"),
							gomock.Eq("mock-log-2"),
						).
						Times(1)

					log.Warnf("%s = %s", "mock-log-1", "mock-log-2")
				})
			})
		})

		Context("Infoln function", Label("unit"), func() {
			When("success send log", func() {
				It("should log the message", func() {
					client.
						EXPECT().
						Infoln(gomock.Eq("mock-log-1"), gomock.Eq("mock-log-2")).
						Times(1)

					log.Infoln("mock-log-1", "mock-log-2")
				})
			})
		})

		Context("Debugln function", Label("unit"), func() {
			When("success send log", func() {
				It("should log the message", func() {
					client.
						EXPECT().
						Debugln(gomock.Eq("mock-log-1"), gomock.Eq("mock-log-2")).
						Times(1)

					log.Debugln("mock-log-1", "mock-log-2")
				})
			})
		})

		Context("Errorln function", Label("unit"), func() {
			When("success send log", func() {
				It("should log the message", func() {
					client.
						EXPECT().
						Errorln(gomock.Eq("mock-log-1"), gomock.Eq("mock-log-2")).
						Times(1)

					log.Errorln("mock-log-1", "mock-log-2")
				})
			})
		})

		Context("Warnln function", Label("unit"), func() {
			When("success send log", func() {
				It("should log the message", func() {
					client.
						EXPECT().
						Warnln(gomock.Eq("mock-log-1"), gomock.Eq("mock-log-2")).
						Times(1)

					log.Warnln("mock-log-1", "mock-log-2")
				})
			})
		})

		Context("WithFields function", Label("unit"), func() {
			When("success set fields", func() {
				It("should return result", func() {
					p := map[string]interface{}{
						"key": "value",
					}
					client.
						EXPECT().
						WithFields(gomock.Eq(p)).
						Return(client).
						Times(1)

					res := log.WithFields(p)

					Expect(res).To(Equal(log))
				})
			})
		})
	})

})
