package mysql_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-seidon/hippo/internal/repository"
	repository_mysql "github.com/go-seidon/hippo/internal/repository/mysql"
	mock_datetime "github.com/go-seidon/provider/datetime/mock"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("File Repository", func() {
	Context("NewFileRepository function", Label("unit"), func() {
		When("master db client is not specified", func() {
			It("should return error", func() {
				res, err := repository_mysql.NewFileRepository()

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid db client specified")))
			})
		})

		When("replica db client is not specified", func() {
			It("should return error", func() {
				mOpt := repository_mysql.WithDbMaster(&sql.DB{})
				res, err := repository_mysql.NewFileRepository(mOpt)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid db client specified")))
			})
		})

		When("required parameter is specified", func() {
			It("should return result", func() {
				mOpt := repository_mysql.WithDbMaster(&sql.DB{})
				rOpt := repository_mysql.WithDbReplica(&sql.DB{})
				res, err := repository_mysql.NewFileRepository(mOpt, rOpt)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("clock is specified", func() {
			It("should return result", func() {
				clockOpt := repository_mysql.WithClock(&mock_datetime.MockClock{})
				mOpt := repository_mysql.WithDbMaster(&sql.DB{})
				rOpt := repository_mysql.WithDbReplica(&sql.DB{})
				res, err := repository_mysql.NewFileRepository(clockOpt, mOpt, rOpt)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("DeleteFile function", Label("unit"), func() {
		var (
			ctx              context.Context
			currentTimestamp time.Time
			clock            *mock_datetime.MockClock
			dbClient         sqlmock.Sqlmock
			repo             repository.FileRepository
			p                repository.DeleteFileParam
			findFileQuery    string
			deleteFileQuery  string
			fileRows         *sqlmock.Rows
		)

		BeforeEach(func() {
			ctx = context.Background()
			t := GinkgoT()

			currentTimestamp = time.Now()
			ctrl := gomock.NewController(t)
			clock = mock_datetime.NewMockClock(ctrl)
			clock.EXPECT().Now().Return(currentTimestamp)

			db, mock, err := sqlmock.New()
			if err != nil {
				AbortSuite("failed create db mock: " + err.Error())
			}
			dbClient = mock

			clockOpt := repository_mysql.WithClock(clock)
			mOpt := repository_mysql.WithDbMaster(db)
			rOpt := repository_mysql.WithDbReplica(db)
			repo, _ = repository_mysql.NewFileRepository(clockOpt, mOpt, rOpt)

			p = repository.DeleteFileParam{
				UniqueId: "mock-unique-id",
				DeleteFn: func(ctx context.Context, p repository.DeleteFnParam) error {
					return nil
				},
			}
			findFileQuery = regexp.QuoteMeta(`
				SELECT 
					id, name, path,
					mimetype, extension, size,
					created_at, updated_at, deleted_at
				FROM file
				WHERE id = ?
			`)
			deleteFileQuery = regexp.QuoteMeta(`
				UPDATE file 
				SET deleted_at = ?
				WHERE id = ?
			`)
			fileRows = sqlmock.NewRows([]string{
				"id", "name", "path",
				"mimetype", "extension", "size",
				"created_at", "updated_at", "deleted_at",
			}).AddRow(
				"mock-unique-id",
				"mock-name",
				"mock-path",
				"mock-mimetype",
				"mock-extension",
				0,
				0,
				0,
				nil,
			)
		})

		When("failed start db transaction", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin().
					WillReturnError(fmt.Errorf("failed start db trx"))

				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed start db trx")))
			})
		})

		When("failed scan record", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				rows := sqlmock.NewRows([]string{
					"id", "name", "path",
					"mimetype", "extension", "size",
					"created_at", "updated_at", "deleted_at",
				}).AddRow(
					"mock-unique-id",
					"mock-name",
					"mock-path",
					"mock-mimetype",
					"mock-extension",
					"invalid_int_value", //should be int64
					0,
					0,
					0,
				)
				dbClient.ExpectQuery(findFileQuery).WillReturnRows(rows)
				dbClient.ExpectRollback()

				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("sql: Scan error on column index 5, name \"size\": converting driver.Value type string (\"invalid_int_value\") to a int64: invalid syntax"))
			})
		})

		When("record is not found", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.ExpectQuery(findFileQuery).
					WillReturnError(sql.ErrNoRows)
				dbClient.ExpectRollback()

				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrNotFound))
			})
		})

		When("failed find file record", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.ExpectQuery(findFileQuery).
					WillReturnError(fmt.Errorf("db error"))
				dbClient.ExpectRollback()

				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("db error")))
			})
		})

		When("failed rollback find file trx", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.ExpectQuery(findFileQuery).
					WillReturnError(fmt.Errorf("db error"))
				dbClient.ExpectRollback().
					WillReturnError(fmt.Errorf("failed rollback"))

				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed rollback")))
			})
		})

		When("failed rollback file is deleted", func() {
			It("should return error", func() {
				fileRows = sqlmock.NewRows([]string{
					"id", "name", "path",
					"mimetype", "extension", "size",
					"created_at", "updated_at", "deleted_at",
				}).AddRow(
					"mock-unique-id",
					"mock-name",
					"mock-path",
					"mock-mimetype",
					"mock-extension",
					0,
					0,
					0,
					1, //deleted
				)

				dbClient.ExpectBegin()
				dbClient.ExpectQuery(findFileQuery).WillReturnRows(fileRows)
				dbClient.ExpectRollback().WillReturnError(fmt.Errorf("rollback error"))

				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("file is deleted", func() {
			It("should return error", func() {
				fileRows = sqlmock.NewRows([]string{
					"id", "name", "path",
					"mimetype", "extension", "size",
					"created_at", "updated_at", "deleted_at",
				}).AddRow(
					"mock-unique-id",
					"mock-name",
					"mock-path",
					"mock-mimetype",
					"mock-extension",
					0,
					0,
					0,
					1, //deleted
				)

				dbClient.ExpectBegin()
				dbClient.ExpectQuery(findFileQuery).WillReturnRows(fileRows)
				dbClient.ExpectRollback()

				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrDeleted))
			})
		})

		When("failed update file record", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.ExpectQuery(findFileQuery).WillReturnRows(fileRows)
				dbClient.ExpectExec(deleteFileQuery).WillReturnError(fmt.Errorf("db error"))
				dbClient.ExpectRollback()

				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("db error")))
			})
		})

		When("failed rollback update file db trx", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.ExpectQuery(findFileQuery).WillReturnRows(fileRows)
				dbClient.
					ExpectExec(deleteFileQuery).
					WithArgs(
						currentTimestamp.UnixMilli(),
						p.UniqueId,
					).
					WillReturnError(fmt.Errorf("db error"))
				dbClient.ExpectRollback().WillReturnError(fmt.Errorf("rollback error"))

				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("total affected row is not 1", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.ExpectQuery(findFileQuery).WillReturnRows(fileRows)
				dbClient.
					ExpectExec(deleteFileQuery).
					WithArgs(
						currentTimestamp.UnixMilli(),
						p.UniqueId,
					).
					WillReturnResult(driver.ResultNoRows)
				dbClient.ExpectRollback()

				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("record is not updated")))
			})
		})

		When("failed rollback total affected row trx", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.ExpectQuery(findFileQuery).WillReturnRows(fileRows)
				dbClient.
					ExpectExec(deleteFileQuery).
					WithArgs(
						currentTimestamp.UnixMilli(),
						p.UniqueId,
					).
					WillReturnResult(driver.ResultNoRows)
				dbClient.ExpectRollback().WillReturnError(fmt.Errorf("rollback error"))

				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("failed execute delete function", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.ExpectQuery(findFileQuery).WillReturnRows(fileRows)
				dbClient.
					ExpectExec(deleteFileQuery).
					WithArgs(
						currentTimestamp.UnixMilli(),
						p.UniqueId,
					).
					WillReturnResult(driver.RowsAffected(1))
				p.DeleteFn = func(ctx context.Context, p repository.DeleteFnParam) error {
					return fmt.Errorf("delete fn error")
				}
				dbClient.ExpectRollback()

				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("delete fn error")))
			})
		})

		When("failed rollback delete fn db trx", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.ExpectQuery(findFileQuery).WillReturnRows(fileRows)
				dbClient.
					ExpectExec(deleteFileQuery).
					WithArgs(
						currentTimestamp.UnixMilli(),
						p.UniqueId,
					).
					WillReturnResult(driver.RowsAffected(1))
				p.DeleteFn = func(ctx context.Context, p repository.DeleteFnParam) error {
					return fmt.Errorf("delete fn error")
				}
				dbClient.ExpectRollback().WillReturnError(fmt.Errorf("rollback error"))

				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("failed commit db trx", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.ExpectQuery(findFileQuery).WillReturnRows(fileRows)
				dbClient.
					ExpectExec(deleteFileQuery).
					WithArgs(
						currentTimestamp.UnixMilli(),
						p.UniqueId,
					).
					WillReturnResult(driver.RowsAffected(1))
				dbClient.ExpectCommit().WillReturnError(fmt.Errorf("commit error"))

				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("commit error")))
			})
		})

		When("success delete file", func() {
			It("should return result", func() {
				dbClient.ExpectBegin()
				dbClient.ExpectQuery(findFileQuery).WillReturnRows(fileRows)
				dbClient.
					ExpectExec(deleteFileQuery).
					WithArgs(
						currentTimestamp.UnixMilli(),
						p.UniqueId,
					).
					WillReturnResult(driver.RowsAffected(1))
				dbClient.ExpectCommit()

				res, err := repo.DeleteFile(ctx, p)

				expectedRes := &repository.DeleteFileResult{
					DeletedAt: currentTimestamp,
				}
				Expect(res).To(Equal(expectedRes))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("DeleteFile function", Label("integration"), Ordered, func() {
		var (
			ctx    context.Context
			client *sql.DB
			repo   repository.FileRepository
			p      repository.DeleteFileParam
		)

		BeforeAll(func() {
			dbClient, err := OpenDb("")
			if err != nil {
				AbortSuite("failed open test db: " + err.Error())
			}
			client = dbClient

			err = RunDbMigration(client)
			if err != nil {
				AbortSuite("failed prepare db migration: " + err.Error())
			}

			ctx = context.Background()
			mOpt := repository_mysql.WithDbMaster(client)
			rOpt := repository_mysql.WithDbReplica(client)
			repo, _ = repository_mysql.NewFileRepository(mOpt, rOpt)
		})

		BeforeEach(func() {
			p = repository.DeleteFileParam{
				UniqueId: "mock-unique-id",
				DeleteFn: func(ctx context.Context, p repository.DeleteFnParam) error {
					return nil
				},
			}
			err := InsertDummyFile(client, InsertDummyFileParam{
				UniqueId: p.UniqueId,
			})
			if err != nil {
				AbortSuite("failed prepare seed data: " + err.Error())
			}
		})

		AfterEach(func() {
			client.Exec("TRUNCATE file")
		})

		AfterAll(func() {
			client.Close()
		})

		When("file record is not available", func() {
			It("should return error", func() {
				p.UniqueId = "unavailable-unique-id"
				res, err := repo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrNotFound))
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

	Context("RetrieveFile function", Label("unit"), func() {
		var (
			ctx           context.Context
			dbClient      sqlmock.Sqlmock
			repo          repository.FileRepository
			p             repository.RetrieveFileParam
			findFileQuery string
			fileRows      *sqlmock.Rows
		)

		BeforeEach(func() {
			ctx = context.Background()

			db, mock, err := sqlmock.New()
			if err != nil {
				AbortSuite("failed create db mock: " + err.Error())
			}
			dbClient = mock

			mOpt := repository_mysql.WithDbMaster(db)
			rOpt := repository_mysql.WithDbReplica(db)
			repo, _ = repository_mysql.NewFileRepository(mOpt, rOpt)

			p = repository.RetrieveFileParam{
				UniqueId: "mock-unique-id",
			}
			findFileQuery = regexp.QuoteMeta(`
				SELECT 
					id, name, path,
					mimetype, extension, size,
					created_at, updated_at, deleted_at
				FROM file
				WHERE id = ?
			`)
			fileRows = sqlmock.NewRows([]string{
				"id", "name", "path",
				"mimetype", "extension", "size",
				"created_at", "updated_at", "deleted_at",
			}).AddRow(
				"mock-unique-id",
				"mock-name",
				"mock-path",
				"mock-mimetype",
				"mock-extension",
				0,
				0,
				0,
				nil,
			)
		})

		When("failed scan record", func() {
			It("should return error", func() {
				rows := sqlmock.NewRows([]string{
					"id", "name", "path",
					"mimetype", "extension", "size",
					"created_at", "updated_at", "deleted_at",
				}).AddRow(
					"mock-unique-id",
					"mock-name",
					"mock-path",
					"mock-mimetype",
					"mock-extension",
					"invalid_int_value", //should be int64
					0,
					0,
					0,
				)
				dbClient.ExpectQuery(findFileQuery).WillReturnRows(rows)

				res, err := repo.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("sql: Scan error on column index 5, name \"size\": converting driver.Value type string (\"invalid_int_value\") to a int64: invalid syntax"))
			})
		})

		When("record is not found", func() {
			It("should return error", func() {
				dbClient.ExpectQuery(findFileQuery).
					WillReturnError(sql.ErrNoRows)

				res, err := repo.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrNotFound))
			})
		})

		When("failed find file record", func() {
			It("should return error", func() {
				dbClient.ExpectQuery(findFileQuery).
					WillReturnError(fmt.Errorf("db error"))

				res, err := repo.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("db error")))
			})
		})

		When("file is deleted", func() {
			It("should return error", func() {
				fileRows = sqlmock.NewRows([]string{
					"id", "name", "path",
					"mimetype", "extension", "size",
					"created_at", "updated_at", "deleted_at",
				}).AddRow(
					"mock-unique-id",
					"mock-name",
					"mock-path",
					"mock-mimetype",
					"mock-extension",
					0,
					0,
					0,
					1,
				)
				dbClient.ExpectQuery(findFileQuery).
					WillReturnRows(fileRows)

				res, err := repo.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrDeleted))
			})
		})

		When("success find file", func() {
			It("should return result", func() {
				dbClient.ExpectQuery(findFileQuery).
					WillReturnRows(fileRows)

				res, err := repo.RetrieveFile(ctx, p)

				eRes := &repository.RetrieveFileResult{
					UniqueId:  "mock-unique-id",
					Name:      "mock-name",
					Path:      "mock-path",
					MimeType:  "mock-mimetype",
					Extension: "mock-extension",
				}
				Expect(res).To(Equal(eRes))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("CreateFile function", Label("unit"), func() {
		var (
			ctx              context.Context
			currentTimestamp time.Time
			dbClient         sqlmock.Sqlmock
			repo             repository.FileRepository
			p                repository.CreateFileParam
			insertSqlQuery   string
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			ctx = context.Background()
			currentTimestamp = time.Now()
			clock := mock_datetime.NewMockClock(ctrl)
			clock.
				EXPECT().
				Now().
				Return(currentTimestamp).
				Times(1)

			db, mock, err := sqlmock.New()
			if err != nil {
				AbortSuite("failed create db mock: " + err.Error())
			}
			dbClient = mock

			clockOpt := repository_mysql.WithClock(clock)
			mOpt := repository_mysql.WithDbMaster(db)
			rOpt := repository_mysql.WithDbReplica(db)
			repo, _ = repository_mysql.NewFileRepository(clockOpt, mOpt, rOpt)

			p = repository.CreateFileParam{
				UniqueId:  "mock-unique-id",
				Name:      "mock-name",
				Path:      "/temp",
				Mimetype:  "image/jpeg",
				Extension: "jpg",
				Size:      200,
				CreateFn: func(ctx context.Context, p repository.CreateFnParam) error {
					return nil
				},
			}
			insertSqlQuery = regexp.QuoteMeta(`
				INSERT INTO file (
					id, name, path, 
					mimetype, extension, size, 
					created_at, updated_at
				) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			`)
		})

		When("failed start db trx", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin().
					WillReturnError(fmt.Errorf("db error"))

				res, err := repo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("db error")))
			})
		})

		When("failed rollback insert record", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.
					ExpectExec(insertSqlQuery).
					WithArgs(
						p.UniqueId, p.Name, p.Path,
						p.Mimetype, p.Extension, p.Size,
						currentTimestamp.UnixMilli(),
						currentTimestamp.UnixMilli(),
					).
					WillReturnError(fmt.Errorf("insert error"))
				dbClient.ExpectRollback().
					WillReturnError(fmt.Errorf("rollback error"))

				res, err := repo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("failed insert record", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.
					ExpectExec(insertSqlQuery).
					WithArgs(
						p.UniqueId, p.Name, p.Path,
						p.Mimetype, p.Extension, p.Size,
						currentTimestamp.UnixMilli(),
						currentTimestamp.UnixMilli(),
					).
					WillReturnError(fmt.Errorf("insert error"))
				dbClient.ExpectRollback()

				res, err := repo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("insert error")))
			})
		})

		When("failed rollback execute create fn", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.
					ExpectExec(insertSqlQuery).
					WithArgs(
						p.UniqueId, p.Name, p.Path,
						p.Mimetype, p.Extension, p.Size,
						currentTimestamp.UnixMilli(),
						currentTimestamp.UnixMilli(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				p.CreateFn = func(ctx context.Context, p repository.CreateFnParam) error {
					return fmt.Errorf("execute error")
				}
				dbClient.
					ExpectRollback().
					WillReturnError(fmt.Errorf("rollback error"))

				res, err := repo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("failed execute create fn", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.
					ExpectExec(insertSqlQuery).
					WithArgs(
						p.UniqueId, p.Name, p.Path,
						p.Mimetype, p.Extension, p.Size,
						currentTimestamp.UnixMilli(),
						currentTimestamp.UnixMilli(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				p.CreateFn = func(ctx context.Context, p repository.CreateFnParam) error {
					return fmt.Errorf("execute error")
				}
				dbClient.ExpectRollback()

				res, err := repo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("execute error")))
			})
		})

		When("failed commit db trx", func() {
			It("should return error", func() {
				dbClient.ExpectBegin()
				dbClient.
					ExpectExec(insertSqlQuery).
					WithArgs(
						p.UniqueId, p.Name, p.Path,
						p.Mimetype, p.Extension, p.Size,
						currentTimestamp.UnixMilli(),
						currentTimestamp.UnixMilli(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				dbClient.
					ExpectCommit().
					WillReturnError(fmt.Errorf("commit error"))

				res, err := repo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("commit error")))
			})
		})

		When("success create file", func() {
			It("should return result", func() {
				dbClient.ExpectBegin()
				dbClient.
					ExpectExec(insertSqlQuery).
					WithArgs(
						p.UniqueId, p.Name, p.Path,
						p.Mimetype, p.Extension, p.Size,
						currentTimestamp.UnixMilli(),
						currentTimestamp.UnixMilli(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				dbClient.ExpectCommit()

				res, err := repo.CreateFile(ctx, p)

				expectedRes := &repository.CreateFileResult{
					UniqueId:  p.UniqueId,
					Name:      p.Name,
					Path:      p.Path,
					Mimetype:  p.Mimetype,
					Extension: p.Extension,
					Size:      p.Size,
					CreatedAt: currentTimestamp,
				}
				Expect(res).To(Equal(expectedRes))
				Expect(err).To(BeNil())
			})
		})
	})

})
