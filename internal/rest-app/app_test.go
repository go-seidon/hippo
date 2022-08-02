package rest_app_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/go-seidon/local/internal/app"
	"github.com/go-seidon/local/internal/mock"
	rest_app "github.com/go-seidon/local/internal/rest-app"
)

func TestRestApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rest App Package")
}

var _ = Describe("Response Package", func() {

	Context("NewRestApp function", Label("unit"), func() {
		var (
			log *mock.MockLogger
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			log = mock.NewMockLogger(ctrl)
		})

		When("config is not specified", func() {
			It("should return error", func() {
				res, err := rest_app.NewRestApp(
					rest_app.WithLogger(log),
				)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid rest app config")))
			})
		})

		When("db provider is not supported", func() {
			It("should return result", func() {
				res, err := rest_app.NewRestApp(
					rest_app.WithConfig(app.Config{
						DBProvider: "invalid db provider",
					}),
				)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("unsupported db provider")))
			})
		})

		When("logger is not specified", func() {
			It("should return result", func() {
				res, err := rest_app.NewRestApp(
					rest_app.WithConfig(app.Config{
						DBProvider: app.DB_PROVIDER_MYSQL,
						AppDebug:   true,
					}),
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("parameter is specified", func() {
			It("should return result", func() {
				log.
					EXPECT().
					WriterLevel(gomock.Eq("error")).
					Times(1)

				res, err := rest_app.NewRestApp(
					rest_app.WithLogger(log),
					rest_app.WithConfig(app.Config{
						DBProvider: app.DB_PROVIDER_MYSQL,
						AppEnv:     "local",
					}),
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("RestAppConfig", Label("unit"), func() {
		var (
			cfg *rest_app.RestAppConfig
		)

		BeforeEach(func() {
			cfg = &rest_app.RestAppConfig{
				AppName:    "mock-name",
				AppVersion: "mock-version",
				AppHost:    "host",
				AppPort:    3000,
			}
		})

		When("GetAppName function is called", func() {
			It("should return app name", func() {
				r := cfg.GetAppName()

				Expect(r).To(Equal("mock-name"))
			})
		})

		When("GetAppVersion function is called", func() {
			It("should return app name", func() {
				r := cfg.GetAppVersion()

				Expect(r).To(Equal("mock-version"))
			})
		})

		When("GetAddress function is called", func() {
			It("should return app name", func() {
				r := cfg.GetAddress()

				Expect(r).To(Equal("host:3000"))
			})
		})
	})

	Context("Run function", Label("unit"), func() {
		var (
			ra            app.App
			logger        *mock.MockLogger
			server        *mock.MockServer
			healthService *mock.MockHealthCheck
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			logger = mock.NewMockLogger(ctrl)
			healthService = mock.NewMockHealthCheck(ctrl)
			server = mock.NewMockServer(ctrl)
			ra, _ = rest_app.NewRestApp(
				rest_app.WithConfig(app.Config{
					AppName:     "mock-name",
					AppVersion:  "mock-version",
					RESTAppHost: "localhost",
					RESTAppPort: 4949,
					DBProvider:  "mysql",
				}),
				rest_app.WithLogger(logger),
				rest_app.WithServer(server),
				rest_app.WithService(healthService),
			)
		})

		When("failed start healthcehck", func() {
			It("should return error", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Running %s:%s"), gomock.Eq("mock-name"), gomock.Eq("mock-version")).
					Times(1)

				healthService.
					EXPECT().
					Start().
					Return(fmt.Errorf("healthcheck error")).
					Times(1)

				err := ra.Run()

				Expect(err).To(Equal(fmt.Errorf("healthcheck error")))
			})
		})

		When("failed listen and serve", func() {
			It("should return error", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Running %s:%s"), gomock.Eq("mock-name"), gomock.Eq("mock-version")).
					Times(1)

				healthService.
					EXPECT().
					Start().
					Return(nil).
					Times(1)

				logger.
					EXPECT().
					Infof(gomock.Eq("Listening on: %s"), gomock.Eq("localhost:4949")).
					Times(1)

				server.
					EXPECT().
					ListenAndServe().
					Return(fmt.Errorf("port already used")).
					Times(1)

				err := ra.Run()

				Expect(err).To(Equal(fmt.Errorf("port already used")))
			})
		})

		When("server is closed", func() {
			It("should return result", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Running %s:%s"), gomock.Eq("mock-name"), gomock.Eq("mock-version")).
					Times(1)

				healthService.
					EXPECT().
					Start().
					Return(nil).
					Times(1)

				logger.
					EXPECT().
					Infof(gomock.Eq("Listening on: %s"), gomock.Eq("localhost:4949")).
					Times(1)

				server.
					EXPECT().
					ListenAndServe().
					Return(http.ErrServerClosed).
					Times(1)

				err := ra.Run()

				Expect(err).To(BeNil())
			})
		})

	})

	Context("Stop function", Label("unit"), func() {
		var (
			ra            app.App
			logger        *mock.MockLogger
			server        *mock.MockServer
			healthService *mock.MockHealthCheck
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			logger = mock.NewMockLogger(ctrl)
			healthService = mock.NewMockHealthCheck(ctrl)
			server = mock.NewMockServer(ctrl)
			ra, _ = rest_app.NewRestApp(
				rest_app.WithConfig(app.Config{
					AppName:     "mock-name",
					AppVersion:  "mock-version",
					RESTAppHost: "localhost",
					RESTAppPort: 4949,
					DBProvider:  "mysql",
				}),
				rest_app.WithLogger(logger),
				rest_app.WithServer(server),
				rest_app.WithService(healthService),
			)
		})

		When("failed stop app", func() {
			It("should return error", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Stopping %s on: %s"), gomock.Eq("mock-name"), gomock.Eq("localhost:4949")).
					Times(1)

				server.
					EXPECT().
					Shutdown(gomock.Eq(context.Background())).
					Return(fmt.Errorf("cant stop app")).
					Times(1)

				err := ra.Stop()

				Expect(err).To(Equal(fmt.Errorf("cant stop app")))
			})
		})

		When("success stop app", func() {
			It("should return result", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Stopping %s on: %s"), gomock.Eq("mock-name"), gomock.Eq("localhost:4949")).
					Times(1)

				server.
					EXPECT().
					Shutdown(gomock.Eq(context.Background())).
					Return(nil).
					Times(1)

				err := ra.Stop()

				Expect(err).To(BeNil())
			})
		})
	})
})
