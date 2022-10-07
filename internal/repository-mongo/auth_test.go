package repository_mongo_test

import (
	"context"
	"fmt"

	mock_datetime "github.com/go-seidon/hippo/internal/datetime/mock"
	"github.com/go-seidon/hippo/internal/repository"
	repository_mongo "github.com/go-seidon/hippo/internal/repository-mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth Repository", func() {

	Context("NewAuthRepository function", Label("unit"), func() {
		When("db client is not specified", func() {
			It("should return error", func() {
				res, err := repository_mongo.NewAuthRepository()

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid db client specified")))
			})
		})

		When("db config is not specified", func() {
			It("should return error", func() {

				mOpt := repository_mongo.WithDbClient(&mongo.Client{})
				res, err := repository_mongo.NewAuthRepository(mOpt)

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
				res, err := repository_mongo.NewAuthRepository(mOpt, dbCfgOpt)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("clock is specified", func() {
			It("should return result", func() {
				clockOpt := repository_mongo.WithClock(&mock_datetime.MockClock{})
				dbCfgOpt := repository_mongo.WithDbConfig(&repository_mongo.DbConfig{
					DbName: "db_name",
				})
				mOpt := repository_mongo.WithDbClient(&mongo.Client{})
				res, err := repository_mongo.NewAuthRepository(clockOpt, mOpt, dbCfgOpt)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("FindClient function", Label("integration"), Ordered, func() {
		var (
			ctx    context.Context
			client *mongo.Client
			repo   repository.AuthRepository
			p      repository.FindClientParam
		)

		BeforeAll(func() {
			dbClient, err := OpenDb("")
			if err != nil {
				AbortSuite("failed open test db: " + err.Error())
			}
			client = dbClient

			err = RunDbMigration(dbClient, RunDbMigrationParam{
				DbName: "hippo_test",
			})
			if err != nil {
				AbortSuite("failed prepare db migration: " + err.Error())
			}
			ctx = context.Background()
			dbCfgOpt := repository_mongo.WithDbConfig(&repository_mongo.DbConfig{
				DbName: "hippo_test",
			})
			dbClientOpt := repository_mongo.WithDbClient(client)
			repo, _ = repository_mongo.NewAuthRepository(dbClientOpt, dbCfgOpt)
		})

		BeforeEach(func() {
			p = repository.FindClientParam{
				ClientId: "mock-client-id",
			}
			err := InsertAuthClient(client, InsertAuthClientParam{
				Id:           "mock-id",
				Name:         "mock-client-name",
				ClientId:     "mock-client-id",
				ClientSecret: "mock-client-secret",
				DbName:       "hippo_test",
			})
			if err != nil {
				AbortSuite("failed prepare seed data: " + err.Error())
			}
		})

		AfterEach(func() {
			_, err := client.
				Database("hippo_test").
				Collection("auth_client").
				DeleteOne(ctx, bson.D{
					{
						Key:   "client_id",
						Value: "mock-client-id",
					},
				})
			if err != nil {
				AbortSuite("failed cleaning seed data: " + err.Error())
			}
		})

		AfterAll(func() {
			err := client.Disconnect(ctx)
			if err != nil {
				AbortSuite("failed close test db: " + err.Error())
			}
		})

		When("client is available", func() {
			It("should return result", func() {
				res, err := repo.FindClient(ctx, p)

				expectedRes := &repository.FindClientResult{
					ClientId:     "mock-client-id",
					ClientSecret: "mock-client-secret",
				}
				Expect(res).To(Equal(expectedRes))
				Expect(err).To(BeNil())
			})
		})

		When("client is not available", func() {
			It("should return error", func() {
				p.ClientId = "invalid-client-id"
				res, err := repo.FindClient(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrorRecordNotFound))
			})
		})
	})

})
