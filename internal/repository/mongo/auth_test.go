package mongo_test

import (
	"context"

	"github.com/go-seidon/hippo/internal/repository"
	repository_mongo "github.com/go-seidon/hippo/internal/repository/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth Repository", func() {
	Context("FindClient function", Label("integration"), Ordered, func() {
		var (
			ctx    context.Context
			client *mongo.Client
			repo   repository.Auth
			p      repository.FindClientParam
			r      *repository.FindClientResult
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
			repo = repository_mongo.NewAuth(dbClientOpt, dbCfgOpt)
		})

		BeforeEach(func() {
			p = repository.FindClientParam{
				ClientId: "mock-client-id",
			}
			r = &repository.FindClientResult{
				Id:           "mock-id",
				Name:         "mock-client-name",
				ClientId:     "mock-client-id",
				ClientSecret: "mock-client-secret",
				Type:         "basic",
				Status:       "active",
			}
			err := InsertAuthClient(client, InsertAuthClientParam{
				Id:           "mock-id",
				Name:         "mock-client-name",
				ClientId:     "mock-client-id",
				ClientSecret: "mock-client-secret",
				Type:         "basic",
				Status:       "active",
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
				DeleteMany(ctx, bson.D{
					{
						Key: "client_id",
						Value: bson.D{
							{
								Key:   "$in",
								Value: []string{"mock-client-id"},
							},
						},
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

		When("client is available using client_id", func() {
			It("should return result", func() {
				res, err := repo.FindClient(ctx, p)

				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})

		When("client is not available", func() {
			It("should return error", func() {
				p.ClientId = "invalid-client-id"
				p.Id = ""
				res, err := repo.FindClient(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrNotFound))
			})
		})

		When("client is available using id", func() {
			It("should return result", func() {
				p.ClientId = ""
				p.Id = "mock-id"
				res, err := repo.FindClient(ctx, p)

				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})
	})

})
