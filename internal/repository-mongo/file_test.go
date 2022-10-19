package repository_mongo_test

import (
	"context"
	"fmt"
	"time"

	"github.com/go-seidon/hippo/internal/repository"
	repository_mongo "github.com/go-seidon/hippo/internal/repository-mongo"
	mock_datetime "github.com/go-seidon/provider/datetime/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ = Describe("File Repository", func() {
	Context("NewFileRepository function", Label("unit"), func() {
		When("db client is not specified", func() {
			It("should return error", func() {
				res, err := repository_mongo.NewFileRepository()

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid db client specified")))
			})
		})

		When("db config is not specified", func() {
			It("should return error", func() {

				mOpt := repository_mongo.WithDbClient(&mongo.Client{})
				res, err := repository_mongo.NewFileRepository(mOpt)

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
				res, err := repository_mongo.NewFileRepository(mOpt, dbCfgOpt)

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
				res, err := repository_mongo.NewFileRepository(clockOpt, mOpt, dbCfgOpt)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("DeleteFile function", Label("integration"), Ordered, func() {
		var (
			ctx    context.Context
			client *mongo.Client
			repo   repository.FileRepository
			p      repository.DeleteFileParam
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
			repo, _ = repository_mongo.NewFileRepository(dbClientOpt, dbCfgOpt)
		})

		BeforeEach(func() {
			p = repository.DeleteFileParam{
				UniqueId: "mock-unique-id",
				DeleteFn: func(ctx context.Context, p repository.DeleteFnParam) error {
					return nil
				},
			}
			err := InsertFile(client, InsertFileParam{
				Id:        "mock-unique-id",
				Name:      "image",
				Path:      "/file/2022",
				Mimetype:  "image/jpeg",
				Extension: "jpeg",
				Size:      200,
				CreatedAt: 1660380011999,
				UpdatedAt: 1660380011999,
				DbName:    "hippo_test",
			})
			if err != nil {
				AbortSuite("failed prepare seed data: " + err.Error())
			}
		})

		AfterEach(func() {
			_, err := client.
				Database("hippo_test").
				Collection("file").
				DeleteOne(ctx, bson.D{
					{
						Key:   "_id",
						Value: "mock-unique-id",
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

		When("file is not available", func() {
			It("should return error", func() {
				p.UniqueId = "invalid-file-id"
				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrorRecordNotFound))
			})
		})

		When("failed proceed callback", func() {
			It("should return error", func() {
				p.DeleteFn = func(ctx context.Context, p repository.DeleteFnParam) error {
					return fmt.Errorf("failed proceed callback")
				}
				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed proceed callback")))
			})
		})

		When("success delete file", func() {
			It("should return result", func() {
				res, err := repo.DeleteFile(ctx, p)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("RetrieveFile function", Label("integration"), Ordered, func() {
		var (
			ctx    context.Context
			client *mongo.Client
			repo   repository.FileRepository
			p      repository.RetrieveFileParam
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
			repo, _ = repository_mongo.NewFileRepository(dbClientOpt, dbCfgOpt)
		})

		BeforeEach(func() {
			p = repository.RetrieveFileParam{
				UniqueId: "mock-unique-id",
			}
			err := InsertFile(client, InsertFileParam{
				Id:        "mock-unique-id",
				Name:      "image",
				Path:      "/file/2022",
				Mimetype:  "image/jpeg",
				Extension: "jpeg",
				Size:      200,
				CreatedAt: 1660380011999,
				UpdatedAt: 1660380011999,
				DbName:    "hippo_test",
			})
			if err != nil {
				AbortSuite("failed prepare seed data: " + err.Error())
			}
		})

		AfterEach(func() {
			_, err := client.
				Database("hippo_test").
				Collection("file").
				DeleteOne(ctx, bson.D{
					{
						Key:   "_id",
						Value: "mock-unique-id",
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

		When("file is not available", func() {
			It("should return error", func() {
				p.UniqueId = "invalid-file-id"
				res, err := repo.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrorRecordNotFound))
			})
		})

		When("success retrieve file", func() {
			It("should return result", func() {
				res, err := repo.RetrieveFile(ctx, p)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("CreateFile function", Label("integration"), Ordered, func() {
		var (
			ctx    context.Context
			client *mongo.Client
			repo   repository.FileRepository
			p      repository.CreateFileParam
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
			repo, _ = repository_mongo.NewFileRepository(dbClientOpt, dbCfgOpt)
		})

		BeforeEach(func() {
			p = repository.CreateFileParam{
				UniqueId:  "mock-unique-id",
				Name:      "image",
				Path:      "/file/2022",
				Mimetype:  "image/jpeg",
				Extension: "jpeg",
				Size:      200,
				CreateFn: func(ctx context.Context, p repository.CreateFnParam) error {
					return nil
				},
			}
		})

		AfterEach(func() {
			_, err := client.
				Database("hippo_test").
				Collection("file").
				DeleteOne(ctx, bson.D{
					{
						Key:   "_id",
						Value: "mock-unique-id",
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

		When("failed proceed callback", func() {
			It("shold return error", func() {
				p.CreateFn = func(ctx context.Context, p repository.CreateFnParam) error {
					return fmt.Errorf("failed proceed callback")
				}
				res, err := repo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed proceed callback")))
			})
		})

		When("success create file", func() {
			It("shold return error", func() {
				res, err := repo.CreateFile(ctx, p)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("failed create file", func() {
			It("shold return error", func() {
				currentTimestamp := time.Now()
				_, err := client.
					Database("hippo_test").
					Collection("file").
					InsertOne(ctx, bson.D{
						{
							Key:   "_id",
							Value: p.UniqueId,
						},
						{
							Key:   "name",
							Value: p.Name,
						},
						{
							Key:   "path",
							Value: p.Path,
						},
						{
							Key:   "mimetype",
							Value: p.Mimetype,
						},
						{
							Key:   "extension",
							Value: p.Extension,
						},
						{
							Key:   "size",
							Value: p.Size,
						},
						{
							Key:   "created_at",
							Value: currentTimestamp,
						},
						{
							Key:   "updated_at",
							Value: currentTimestamp,
						},
					})
				if err != nil {
					AbortSuite("failed prepare dummy data: " + err.Error())
				}

				res, err := repo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
