package restapp_test

import (
	"context"
	"fmt"

	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/go-seidon/hippo/internal/app"
	mock_restapp "github.com/go-seidon/hippo/internal/restapp/mock"
	mock_healthcheck "github.com/go-seidon/provider/health/mock"
	mock_logging "github.com/go-seidon/provider/logging/mock"

	"github.com/go-seidon/hippo/internal/repository"
	mock_repository "github.com/go-seidon/hippo/internal/repository/mock"
	"github.com/go-seidon/hippo/internal/restapp"
)

func TestRestApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rest App Package")
}

var _ = Describe("App Package", func() {

	Context("NewRestApp function", Label("unit"), func() {
		var (
			log *mock_logging.MockLogger
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			log = mock_logging.NewMockLogger(ctrl)
		})

		When("config is not specified", func() {
			It("should return error", func() {
				res, err := restapp.NewRestApp(
					restapp.WithLogger(log),
				)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid config")))
			})
		})

		When("logger is not specified", func() {
			It("should return result", func() {
				res, err := restapp.NewRestApp(
					restapp.WithConfig(&app.Config{
						RepositoryProvider: repository.PROVIDER_MYSQL,
					}),
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("debug is enabled", func() {
			It("should return result", func() {
				res, err := restapp.NewRestApp(
					restapp.WithConfig(&app.Config{
						RepositoryProvider: repository.PROVIDER_MYSQL,
						AppDebug:           true,
					}),
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("env is specified", func() {
			It("should return result", func() {
				res, err := restapp.NewRestApp(
					restapp.WithConfig(&app.Config{
						RepositoryProvider: repository.PROVIDER_MYSQL,
						AppDebug:           true,
						AppEnv:             "local",
					}),
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("parameter is specified", func() {
			It("should return result", func() {
				res, err := restapp.NewRestApp(
					restapp.WithLogger(log),
					restapp.WithConfig(&app.Config{
						RepositoryProvider: repository.PROVIDER_MONGO,
						AppEnv:             "local",
						MongoMode:          "standalone",
						MongoAuthMode:      "basic",
					}),
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("RestAppConfig", Label("unit"), func() {
		var (
			cfg *restapp.RestAppConfig
		)

		BeforeEach(func() {
			cfg = &restapp.RestAppConfig{
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
			logger        *mock_logging.MockLogger
			server        *mock_restapp.MockServer
			healthService *mock_healthcheck.MockHealthCheck
			repo          *mock_repository.MockRepository
			ctx           context.Context
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			logger = mock_logging.NewMockLogger(ctrl)
			healthService = mock_healthcheck.NewMockHealthCheck(ctrl)
			server = mock_restapp.NewMockServer(ctrl)
			repo = mock_repository.NewMockRepository(ctrl)
			fileRepo := mock_repository.NewMockFile(ctrl)
			authRepo := mock_repository.NewMockAuth(ctrl)
			repo.EXPECT().GetFile().Return(fileRepo).AnyTimes()
			repo.EXPECT().GetAuth().Return(authRepo).AnyTimes()

			ra, _ = restapp.NewRestApp(
				restapp.WithConfig(&app.Config{
					AppName:            "mock-name",
					AppVersion:         "mock-version",
					RESTAppHost:        "localhost",
					RESTAppPort:        4949,
					RepositoryProvider: "mysql",
				}),
				restapp.WithLogger(logger),
				restapp.WithServer(server),
				restapp.WithService(healthService),
				restapp.WithRepository(repo),
			)

			ctx = context.Background()
		})

		When("failed start healthcehck", func() {
			It("should return error", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Running %s:%s"), gomock.Eq("mock-name-rest"), gomock.Eq("mock-version")).
					Times(1)

				healthService.
					EXPECT().
					Start(gomock.Eq(ctx)).
					Return(fmt.Errorf("healthcheck error")).
					Times(1)

				err := ra.Run(ctx)

				Expect(err).To(Equal(fmt.Errorf("healthcheck error")))
			})
		})

		When("failed init repo", func() {
			It("should return error", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Running %s:%s"), gomock.Eq("mock-name-rest"), gomock.Eq("mock-version")).
					Times(1)

				healthService.
					EXPECT().
					Start(gomock.Eq(ctx)).
					Return(nil).
					Times(1)

				repo.
					EXPECT().
					Init(gomock.Eq(ctx)).
					Return(fmt.Errorf("db error")).
					Times(1)

				err := ra.Run(ctx)

				Expect(err).To(Equal(fmt.Errorf("db error")))
			})
		})

		When("failed listen and serve", func() {
			It("should return error", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Running %s:%s"), gomock.Eq("mock-name-rest"), gomock.Eq("mock-version")).
					Times(1)

				healthService.
					EXPECT().
					Start(gomock.Eq(ctx)).
					Return(nil).
					Times(1)

				repo.
					EXPECT().
					Init(gomock.Eq(ctx)).
					Return(nil).
					Times(1)

				logger.
					EXPECT().
					Infof(gomock.Eq("Listening on: %s"), gomock.Eq("localhost:4949")).
					Times(1)

				server.
					EXPECT().
					Start(gomock.Any()).
					Return(fmt.Errorf("port already used")).
					Times(1)

				err := ra.Run(ctx)

				Expect(err).To(Equal(fmt.Errorf("port already used")))
			})
		})

		When("server is closed", func() {
			It("should return result", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Running %s:%s"), gomock.Eq("mock-name-rest"), gomock.Eq("mock-version")).
					Times(1)

				healthService.
					EXPECT().
					Start(gomock.Eq(ctx)).
					Return(nil).
					Times(1)

				repo.
					EXPECT().
					Init(gomock.Eq(ctx)).
					Return(nil).
					Times(1)

				logger.
					EXPECT().
					Infof(gomock.Eq("Listening on: %s"), gomock.Eq("localhost:4949")).
					Times(1)

				server.
					EXPECT().
					Start(gomock.Any()).
					Return(http.ErrServerClosed).
					Times(1)

				err := ra.Run(ctx)

				Expect(err).To(BeNil())
			})
		})
	})

	Context("Stop function", Label("unit"), func() {
		var (
			ra            app.App
			logger        *mock_logging.MockLogger
			server        *mock_restapp.MockServer
			healthService *mock_healthcheck.MockHealthCheck
			ctx           context.Context
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			logger = mock_logging.NewMockLogger(ctrl)
			healthService = mock_healthcheck.NewMockHealthCheck(ctrl)
			server = mock_restapp.NewMockServer(ctrl)
			ra, _ = restapp.NewRestApp(
				restapp.WithConfig(&app.Config{
					AppName:            "mock-name",
					AppVersion:         "mock-version",
					RESTAppHost:        "localhost",
					RESTAppPort:        4949,
					RepositoryProvider: "mysql",
				}),
				restapp.WithLogger(logger),
				restapp.WithServer(server),
				restapp.WithService(healthService),
			)
			ctx = context.Background()
		})

		When("failed stop app", func() {
			It("should return error", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Stopping %s on: %s"), gomock.Eq("mock-name-rest"), gomock.Eq("localhost:4949")).
					Times(1)

				healthService.
					EXPECT().
					Stop(gomock.Eq(ctx)).
					Return(nil).
					Times(1)

				server.
					EXPECT().
					Shutdown(gomock.Eq(context.Background())).
					Return(fmt.Errorf("cant stop app")).
					Times(1)

				err := ra.Stop(ctx)

				Expect(err).To(Equal(fmt.Errorf("cant stop app")))
			})
		})

		When("success stop app", func() {
			It("should return result", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Stopping %s on: %s"), gomock.Eq("mock-name-rest"), gomock.Eq("localhost:4949")).
					Times(1)

				healthService.
					EXPECT().
					Stop(gomock.Eq(ctx)).
					Return(nil).
					Times(1)

				server.
					EXPECT().
					Shutdown(gomock.Eq(context.Background())).
					Return(nil).
					Times(1)

				err := ra.Stop(ctx)

				Expect(err).To(BeNil())
			})
		})

		When("failed stop healthcheck", func() {
			It("should log the error", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Stopping %s on: %s"), gomock.Eq("mock-name-rest"), gomock.Eq("localhost:4949")).
					Times(1)

				healthService.
					EXPECT().
					Stop(gomock.Eq(ctx)).
					Return(fmt.Errorf("routine error")).
					Times(1)

				logger.
					EXPECT().
					Errorf(gomock.Eq("Failed stopping healthcheck, err: %s"), gomock.Eq("routine error")).
					Times(1)

				server.
					EXPECT().
					Shutdown(gomock.Eq(context.Background())).
					Return(fmt.Errorf("cant stop app")).
					Times(1)

				err := ra.Stop(ctx)

				Expect(err).To(Equal(fmt.Errorf("cant stop app")))
			})
		})
	})
})
