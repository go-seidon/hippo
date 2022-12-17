package grpcapp_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-seidon/hippo/internal/app"
	"github.com/go-seidon/hippo/internal/grpcapp"
	mock_grpcapp "github.com/go-seidon/hippo/internal/grpcapp/mock"
	mock_healthcheck "github.com/go-seidon/hippo/internal/healthcheck/mock"
	mock_repository "github.com/go-seidon/hippo/internal/repository/mock"
	mock_logging "github.com/go-seidon/provider/logging/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

func TestGrpcApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Grpc App Package")
}

var _ = Describe("App Package", func() {

	Context("NewGrpcApp function", Label("unit"), func() {
		var (
			cfg           *app.Config
			logger        *mock_logging.MockLogger
			repository    *mock_repository.MockProvider
			healthService *mock_healthcheck.MockHealthCheck
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			cfg = &app.Config{
				AppDebug:           true,
				AppEnv:             "local",
				RepositoryProvider: "mongo",
			}
			logger = mock_logging.NewMockLogger(ctrl)
			repository = mock_repository.NewMockProvider(ctrl)
			fileRepo := mock_repository.NewMockFile(ctrl)
			authRepo := mock_repository.NewMockAuth(ctrl)
			repository.EXPECT().GetFileRepo().Return(fileRepo).AnyTimes()
			repository.EXPECT().GetAuthRepo().Return(authRepo).AnyTimes()
			healthService = mock_healthcheck.NewMockHealthCheck(ctrl)
		})

		When("config is not specified", func() {
			It("should return error", func() {
				res, err := grpcapp.NewGrpcApp()

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid config")))
			})
		})

		When("logger is specified", func() {
			It("should return result", func() {
				res, err := grpcapp.NewGrpcApp(
					grpcapp.WithConfig(cfg),
					grpcapp.WithLogger(logger),
					grpcapp.WithRepository(repository),
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("health service is specified", func() {
			It("should return result", func() {
				res, err := grpcapp.NewGrpcApp(
					grpcapp.WithConfig(cfg),
					grpcapp.WithService(healthService),
					grpcapp.WithRepository(repository),
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("all parameters are specified", func() {
			It("should return result", func() {
				res, err := grpcapp.NewGrpcApp(
					grpcapp.WithConfig(cfg),
					grpcapp.WithLogger(logger),
					grpcapp.WithRepository(repository),
					grpcapp.WithService(healthService),
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("GrpcAppConfig", Label("unit"), func() {
		var (
			cfg *grpcapp.GrpcAppConfig
		)

		BeforeEach(func() {
			cfg = &grpcapp.GrpcAppConfig{
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
			grpcApp       app.App
			ctx           context.Context
			logger        *mock_logging.MockLogger
			healthService *mock_healthcheck.MockHealthCheck
			server        *mock_grpcapp.MockServer
			repository    *mock_repository.MockProvider
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			cfg := &app.Config{
				AppDebug:    true,
				AppEnv:      "local",
				AppName:     "mock-name",
				AppVersion:  "mock-version",
				GRPCAppHost: "localhost",
				GRPCAppPort: 4949,
			}

			logger = mock_logging.NewMockLogger(ctrl)
			healthService = mock_healthcheck.NewMockHealthCheck(ctrl)
			server = mock_grpcapp.NewMockServer(ctrl)
			repository = mock_repository.NewMockProvider(ctrl)
			fileRepo := mock_repository.NewMockFile(ctrl)
			authRepo := mock_repository.NewMockAuth(ctrl)
			repository.EXPECT().GetFileRepo().Return(fileRepo).AnyTimes()
			repository.EXPECT().GetAuthRepo().Return(authRepo).AnyTimes()

			grpcApp, _ = grpcapp.NewGrpcApp(
				grpcapp.WithConfig(cfg),
				grpcapp.WithLogger(logger),
				grpcapp.WithService(healthService),
				grpcapp.WithServer(server),
				grpcapp.WithRepository(repository),
			)
			ctx = context.Background()
		})

		When("failed start health service", func() {
			It("should return error", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Running %s:%s"), gomock.Eq("mock-name-grpc"), gomock.Eq("mock-version")).
					Times(1)

				healthService.
					EXPECT().
					Start(gomock.Eq(ctx)).
					Return(fmt.Errorf("failed start healthcheck")).
					Times(1)

				err := grpcApp.Run(ctx)

				Expect(err).To(Equal(fmt.Errorf("failed start healthcheck")))
			})
		})

		When("failed init repository", func() {
			It("should return error", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Running %s:%s"), gomock.Eq("mock-name-grpc"), gomock.Eq("mock-version")).
					Times(1)

				healthService.
					EXPECT().
					Start(gomock.Eq(ctx)).
					Return(nil).
					Times(1)

				repository.
					EXPECT().
					Init(gomock.Eq(ctx)).
					Return(fmt.Errorf("db error")).
					Times(1)

				err := grpcApp.Run(ctx)

				Expect(err).To(Equal(fmt.Errorf("db error")))
			})
		})

		When("failed listen and serve", func() {
			It("should return error", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Running %s:%s"), gomock.Eq("mock-name-grpc"), gomock.Eq("mock-version")).
					Times(1)

				healthService.
					EXPECT().
					Start(gomock.Eq(ctx)).
					Return(nil).
					Times(1)

				repository.
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
					ListenAndServe().
					Return(fmt.Errorf("port already used")).
					Times(1)

				err := grpcApp.Run(ctx)

				Expect(err).To(Equal(fmt.Errorf("port already used")))
			})
		})

		When("server is closed", func() {
			It("should return result", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Running %s:%s"), gomock.Eq("mock-name-grpc"), gomock.Eq("mock-version")).
					Times(1)

				healthService.
					EXPECT().
					Start(gomock.Eq(ctx)).
					Return(nil).
					Times(1)

				repository.
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
					ListenAndServe().
					Return(grpc.ErrServerStopped).
					Times(1)

				err := grpcApp.Run(ctx)

				Expect(err).To(BeNil())
			})
		})
	})

	Context("Stop function", Label("unit"), func() {
		var (
			grpcApp       app.App
			ctx           context.Context
			logger        *mock_logging.MockLogger
			healthService *mock_healthcheck.MockHealthCheck
			server        *mock_grpcapp.MockServer
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			cfg := &app.Config{
				AppDebug:    true,
				AppEnv:      "local",
				AppName:     "mock-name",
				AppVersion:  "mock-version",
				GRPCAppHost: "localhost",
				GRPCAppPort: 4949,
			}

			logger = mock_logging.NewMockLogger(ctrl)
			healthService = mock_healthcheck.NewMockHealthCheck(ctrl)
			server = mock_grpcapp.NewMockServer(ctrl)
			repository := mock_repository.NewMockProvider(ctrl)
			fileRepo := mock_repository.NewMockFile(ctrl)
			authRepo := mock_repository.NewMockAuth(ctrl)
			repository.EXPECT().GetFileRepo().Return(fileRepo).AnyTimes()
			repository.EXPECT().GetAuthRepo().Return(authRepo).AnyTimes()

			grpcApp, _ = grpcapp.NewGrpcApp(
				grpcapp.WithConfig(cfg),
				grpcapp.WithLogger(logger),
				grpcapp.WithService(healthService),
				grpcapp.WithServer(server),
				grpcapp.WithRepository(repository),
			)
			ctx = context.Background()
		})

		When("failed stop app", func() {
			It("should return error", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Stopping %s on: %s"), gomock.Eq("mock-name-grpc"), gomock.Eq("localhost:4949")).
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

				err := grpcApp.Stop(ctx)

				Expect(err).To(Equal(fmt.Errorf("cant stop app")))
			})
		})

		When("success stop app", func() {
			It("should return result", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Stopping %s on: %s"), gomock.Eq("mock-name-grpc"), gomock.Eq("localhost:4949")).
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

				err := grpcApp.Stop(ctx)

				Expect(err).To(BeNil())
			})
		})

		When("failed stop healthcheck", func() {
			It("should log the error", func() {
				logger.
					EXPECT().
					Infof(gomock.Eq("Stopping %s on: %s"), gomock.Eq("mock-name-grpc"), gomock.Eq("localhost:4949")).
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

				err := grpcApp.Stop(ctx)

				Expect(err).To(Equal(fmt.Errorf("cant stop app")))
			})
		})
	})
})
