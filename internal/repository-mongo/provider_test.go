package repository_mongo_test

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/mongo"

	mock_db_mongo "github.com/go-seidon/local/internal/db-mongo/mock"
	"github.com/go-seidon/local/internal/repository"
	repository_mongo "github.com/go-seidon/local/internal/repository-mongo"
)

var _ = Describe("Repository Provider", func() {
	Context("NewRepository function", Label("unit"), func() {
		When("db client is not specified", func() {
			It("should return error", func() {
				res, err := repository_mongo.NewRepository()

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid db client specified")))
			})
		})

		When("db config is not specified", func() {
			It("should return error", func() {
				mOpt := repository_mongo.WithDbClient(&mongo.Client{})
				res, err := repository_mongo.NewRepository(mOpt)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid db config specified")))
			})
		})

		When("required parameters are specified", func() {
			It("should return result", func() {
				mOpt := repository_mongo.WithDbClient(&mongo.Client{})
				dbCfgOpt := repository_mongo.WithDbConfig(&repository_mongo.DbConfig{
					DbName: "db_name",
				})
				res, err := repository_mongo.NewRepository(mOpt, dbCfgOpt)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("GetAuthRepo function", Label("unit"), func() {
		var (
			provider repository.Provider
		)

		BeforeEach(func() {
			mOpt := repository_mongo.WithDbClient(&mongo.Client{})
			dbCfgOpt := repository_mongo.WithDbConfig(&repository_mongo.DbConfig{
				DbName: "db_name",
			})
			provider, _ = repository_mongo.NewRepository(mOpt, dbCfgOpt)
		})

		When("function is called", func() {
			It("should return result", func() {
				res := provider.GetAuthRepo()

				Expect(res).ToNot(BeNil())
			})
		})
	})

	Context("GetFileRepo function", Label("unit"), func() {
		var (
			provider repository.Provider
		)

		BeforeEach(func() {
			mOpt := repository_mongo.WithDbClient(&mongo.Client{})
			dbCfgOpt := repository_mongo.WithDbConfig(&repository_mongo.DbConfig{
				DbName: "db_name",
			})
			provider, _ = repository_mongo.NewRepository(mOpt, dbCfgOpt)
		})

		When("function is called", func() {
			It("should return result", func() {
				res := provider.GetFileRepo()

				Expect(res).ToNot(BeNil())
			})
		})
	})

	Context("Init function", Label("unit"), func() {
		var (
			provider repository.Provider
			ctx      context.Context
			dbClient *mock_db_mongo.MockClient
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			dbClient = mock_db_mongo.NewMockClient(ctrl)
			mOpt := repository_mongo.WithDbClient(dbClient)
			dbCfgOpt := repository_mongo.WithDbConfig(&repository_mongo.DbConfig{
				DbName: "db_name",
			})
			provider, _ = repository_mongo.NewRepository(mOpt, dbCfgOpt)
			ctx = context.Background()
		})

		When("success init", func() {
			It("should return result", func() {
				dbClient.
					EXPECT().
					Connect(gomock.Eq(ctx)).
					Return(nil).
					Times(1)

				err := provider.Init(ctx)

				Expect(err).To(BeNil())
			})
		})

		When("failed init", func() {
			It("should return error", func() {
				dbClient.
					EXPECT().
					Connect(gomock.Eq(ctx)).
					Return(fmt.Errorf("db error")).
					Times(1)

				err := provider.Init(ctx)

				Expect(err).To(Equal(fmt.Errorf("db error")))
			})
		})
	})

	Context("Ping function", Label("unit"), func() {
		var (
			provider repository.Provider
			ctx      context.Context
			dbClient *mock_db_mongo.MockClient
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			dbClient = mock_db_mongo.NewMockClient(ctrl)
			mOpt := repository_mongo.WithDbClient(dbClient)
			dbCfgOpt := repository_mongo.WithDbConfig(&repository_mongo.DbConfig{
				DbName: "db_name",
			})
			provider, _ = repository_mongo.NewRepository(mOpt, dbCfgOpt)
			ctx = context.Background()
		})

		When("success ping", func() {
			It("should return result", func() {
				dbClient.
					EXPECT().
					Ping(gomock.Eq(ctx), gomock.Any()).
					Return(nil).
					Times(1)

				err := provider.Ping(ctx)

				Expect(err).To(BeNil())
			})
		})

		When("failed ping", func() {
			It("should return error", func() {
				dbClient.
					EXPECT().
					Ping(gomock.Eq(ctx), gomock.Any()).
					Return(fmt.Errorf("ping error")).
					Times(1)

				err := provider.Ping(ctx)

				Expect(err).To(Equal(fmt.Errorf("ping error")))
			})
		})
	})
})
