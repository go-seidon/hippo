package mongo_test

import (
	"context"
	"time"

	"github.com/go-seidon/hippo/internal/repository"
	repository_mongo "github.com/go-seidon/hippo/internal/repository/mongo"
	"github.com/go-seidon/provider/typeconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth Repository", func() {
	Context("CreateClient function", Label("integration"), Ordered, func() {
		var (
			ctx       context.Context
			currentTs time.Time
			client    *mongo.Client
			repo      repository.Auth
			p         repository.CreateClientParam
			r         *repository.CreateClientResult
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
			currentTs = time.Now().UTC()
			p = repository.CreateClientParam{
				Id:           "create-id",
				Name:         "create-name",
				Type:         "basic",
				Status:       "active",
				ClientId:     "create-client-id",
				ClientSecret: "create-client-secret",
				CreatedAt:    currentTs,
			}
			r = &repository.CreateClientResult{
				Id:           "create-id",
				Name:         "create-name",
				Type:         "basic",
				Status:       "active",
				ClientId:     "create-client-id",
				ClientSecret: "create-client-secret",
				CreatedAt:    time.UnixMilli(currentTs.UnixMilli()).UTC(),
			}
			err := InsertAuthClient(client, InsertAuthClientParam{
				Id:           "exists-id",
				Name:         "exists-client-name",
				ClientId:     "exists-client-id",
				ClientSecret: "exists-client-secret",
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
						Key: "_id",
						Value: bson.D{
							{
								Key:   "$in",
								Value: []string{"create-id", "exists-id"},
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

		When("client is already exists", func() {
			It("should return error", func() {
				p := repository.CreateClientParam{
					Id:           "exists-id",
					Name:         "exists-name",
					Type:         "basic",
					Status:       "active",
					ClientId:     "exists-client-id",
					ClientSecret: "exists-client-secret",
					CreatedAt:    currentTs,
				}
				res, err := repo.CreateClient(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrExists))
			})
		})

		When("success create client", func() {
			It("should return result", func() {
				res, err := repo.CreateClient(ctx, p)

				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})
	})

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

	Context("UpdateClient function", Label("integration"), Ordered, func() {
		var (
			ctx       context.Context
			currentTs time.Time
			client    *mongo.Client
			repo      repository.Auth
			p         repository.UpdateClientParam
			r         *repository.UpdateClientResult
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
			currentTs = time.Now().UTC()
			p = repository.UpdateClientParam{
				Id:        "update-id",
				Name:      "update-name",
				Type:      "basic",
				Status:    "active",
				ClientId:  "update-client-id",
				UpdatedAt: currentTs,
			}
			r = &repository.UpdateClientResult{
				Id:           "update-id",
				Name:         "update-name",
				Type:         "basic",
				Status:       "active",
				ClientId:     "update-client-id",
				ClientSecret: "update-client-secret",
				CreatedAt:    time.UnixMilli(currentTs.UnixMilli()).UTC(),
				UpdatedAt:    time.UnixMilli(currentTs.UnixMilli()).UTC(),
			}
			err := InsertAuthClient(client, InsertAuthClientParam{
				Id:           "update-id",
				Name:         "update-client-name",
				ClientId:     "update-client-id",
				ClientSecret: "update-client-secret",
				Type:         "basic",
				Status:       "active",
				CreatedAt:    currentTs,
				UpdatedAt:    currentTs,
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
						Key: "_id",
						Value: bson.D{
							{
								Key:   "$in",
								Value: []string{"update-id"},
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

		When("client is not available", func() {
			It("should return error", func() {
				p := repository.UpdateClientParam{
					Id:       "invalid-id",
					Name:     "invalid-name",
					Type:     "basic",
					Status:   "active",
					ClientId: "invalid-client-id",
				}
				res, err := repo.UpdateClient(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrNotFound))
			})
		})

		When("success update client", func() {
			It("should return result", func() {
				res, err := repo.UpdateClient(ctx, p)

				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("SearchClient function", Label("integration"), Ordered, func() {
		var (
			ctx       context.Context
			currentTs time.Time
			client    *mongo.Client
			repo      repository.Auth
			p         repository.SearchClientParam
			r         *repository.SearchClientResult
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
			currentTs = time.Now().UTC()
			p = repository.SearchClientParam{
				Limit:    24,
				Offset:   0,
				Keyword:  "search",
				Statuses: []string{"active"},
			}
			r = &repository.SearchClientResult{
				Summary: repository.SearchClientSummary{
					TotalItems: 2,
				},
				Items: []repository.SearchClientItem{
					{
						Id:           "search-1",
						Name:         "search-1",
						ClientId:     "search-1",
						ClientSecret: "search-1",
						Type:         "basic",
						Status:       "active",
						CreatedAt:    time.UnixMilli(currentTs.UnixMilli()).UTC(),
						UpdatedAt:    typeconv.Time(time.UnixMilli(currentTs.UnixMilli()).UTC()),
					},
					{
						Id:           "search-2",
						Name:         "search-2",
						ClientId:     "search-2",
						ClientSecret: "search-2",
						Type:         "basic",
						Status:       "active",
						CreatedAt:    time.UnixMilli(currentTs.UnixMilli()).UTC(),
						UpdatedAt:    typeconv.Time(time.UnixMilli(currentTs.UnixMilli()).UTC()),
					},
				},
			}
			err := InsertAuthClient(client, InsertAuthClientParam{
				Id:           "search-1",
				Name:         "search-1",
				ClientId:     "search-1",
				ClientSecret: "search-1",
				Type:         "basic",
				Status:       "active",
				CreatedAt:    currentTs,
				UpdatedAt:    currentTs,
				DbName:       "hippo_test",
			})
			if err != nil {
				AbortSuite("failed prepare seed data: " + err.Error())
			}
			err = InsertAuthClient(client, InsertAuthClientParam{
				Id:           "search-2",
				Name:         "search-2",
				ClientId:     "search-2",
				ClientSecret: "search-2",
				Type:         "basic",
				Status:       "active",
				CreatedAt:    currentTs,
				UpdatedAt:    currentTs,
				DbName:       "hippo_test",
			})
			if err != nil {
				AbortSuite("failed prepare seed data: " + err.Error())
			}
			err = InsertAuthClient(client, InsertAuthClientParam{
				Id:           "inactive-3",
				Name:         "inactive-3",
				ClientId:     "inactive-3",
				ClientSecret: "inactive-3",
				Type:         "basic",
				Status:       "inactive",
				CreatedAt:    currentTs,
				UpdatedAt:    currentTs,
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
						Key: "_id",
						Value: bson.D{
							{
								Key:   "$in",
								Value: []string{"search-1", "search-2", "inactive-3"},
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

		When("there are no items matched", func() {
			It("should return result", func() {
				p := repository.SearchClientParam{
					Limit:    24,
					Offset:   0,
					Keyword:  "unavailable",
					Statuses: []string{"active"},
				}
				r := &repository.SearchClientResult{
					Summary: repository.SearchClientSummary{
						TotalItems: 0,
					},
					Items: []repository.SearchClientItem{},
				}
				res, err := repo.SearchClient(ctx, p)

				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})

		When("there is one items matched", func() {
			It("should return result", func() {
				p := repository.SearchClientParam{
					Limit:    24,
					Offset:   0,
					Keyword:  "inactive",
					Statuses: []string{},
				}
				r := &repository.SearchClientResult{
					Summary: repository.SearchClientSummary{
						TotalItems: 1,
					},
					Items: []repository.SearchClientItem{
						{
							Id:           "inactive-3",
							Name:         "inactive-3",
							ClientId:     "inactive-3",
							ClientSecret: "inactive-3",
							Type:         "basic",
							Status:       "inactive",
							CreatedAt:    time.UnixMilli(currentTs.UnixMilli()).UTC(),
							UpdatedAt:    typeconv.Time(time.UnixMilli(currentTs.UnixMilli()).UTC()),
						},
					},
				}
				res, err := repo.SearchClient(ctx, p)

				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})

		When("there are some items matched", func() {
			It("should return result", func() {
				res, err := repo.SearchClient(ctx, p)

				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})
	})

})
