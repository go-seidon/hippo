package mysql_test

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-seidon/hippo/internal/repository"
	repository_mysql "github.com/go-seidon/hippo/internal/repository/mysql"
	"github.com/go-seidon/provider/typeconv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ = Describe("File Repository", func() {
	Context("CreateFile function", Label("unit"), func() {
		var (
			ctx        context.Context
			currentTs  time.Time
			dbClient   sqlmock.Sqlmock
			fileRepo   repository.File
			p          repository.CreateFileParam
			checkStmt  string
			insertStmt string
			findStmt   string
		)

		BeforeEach(func() {
			var (
				db  *sql.DB
				err error
			)

			ctx = context.Background()
			currentTs = time.Now().UTC()
			db, dbClient, err = sqlmock.New()
			if err != nil {
				AbortSuite("failed create db mock: " + err.Error())
			}

			gormClient, err := gorm.Open(gorm_mysql.New(gorm_mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing: true,
			})
			if err != nil {
				AbortSuite("failed create gorm client: " + err.Error())
			}
			fileRepo = repository_mysql.NewFile(repository_mysql.FileParam{
				GormClient: gormClient,
			})

			p = repository.CreateFileParam{
				UniqueId:  "id",
				Path:      "storage/id",
				Name:      "dolphin",
				Mimetype:  "image/jpeg",
				Extension: "jpg",
				Size:      2334,
				CreateFn: func(ctx context.Context, p repository.CreateFnParam) error {
					return nil
				},
				CreatedAt: currentTs,
			}
			checkStmt = regexp.QuoteMeta("SELECT `id` FROM `file` WHERE id = ? ORDER BY `file`.`id` LIMIT 1")
			insertStmt = regexp.QuoteMeta("INSERT INTO `file` (`id`,`path`,`name`,`mimetype`,`extension`,`size`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")
			findStmt = regexp.QuoteMeta("SELECT id, name, path, mimetype, extension, size, created_at FROM `file` WHERE id = ? ORDER BY `file`.`id` LIMIT 1")
		})

		AfterEach(func() {
			err := dbClient.ExpectationsWereMet()
			if err != nil {
				AbortSuite("some expectations were not met " + err.Error())
			}
		})

		When("failed begin trx", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin().
					WillReturnError(fmt.Errorf("begin error"))

				res, err := fileRepo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("begin error")))
			})
		})

		When("failed rollback during file existance checking", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(fmt.Errorf("network error"))

				dbClient.
					ExpectRollback().
					WillReturnError(fmt.Errorf("rollback error"))

				res, err := fileRepo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("failed check file existance", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(fmt.Errorf("network error"))

				dbClient.
					ExpectRollback()

				res, err := fileRepo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("file is already exists", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				checkRows := sqlmock.
					NewRows([]string{"id"}).
					AddRow("id")
				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnRows(checkRows)

				dbClient.
					ExpectRollback()

				res, err := fileRepo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrExists))
			})
		})

		When("failed rollback during failed create file", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(gorm.ErrRecordNotFound)

				dbClient.
					ExpectExec(insertStmt).
					WithArgs(
						p.UniqueId,
						p.Path,
						p.Name,
						p.Mimetype,
						p.Extension,
						p.Size,
						p.CreatedAt.UnixMilli(),
						p.CreatedAt.UnixMilli(),
					).
					WillReturnError(fmt.Errorf("network error"))

				dbClient.
					ExpectRollback().
					WillReturnError(fmt.Errorf("rollback error"))

				res, err := fileRepo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("failed create file", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(gorm.ErrRecordNotFound)

				dbClient.
					ExpectExec(insertStmt).
					WithArgs(
						p.UniqueId,
						p.Path,
						p.Name,
						p.Mimetype,
						p.Extension,
						p.Size,
						p.CreatedAt.UnixMilli(),
						p.CreatedAt.UnixMilli(),
					).
					WillReturnError(fmt.Errorf("network error"))

				dbClient.
					ExpectRollback()

				res, err := fileRepo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("failed rollback during check inserted file", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(gorm.ErrRecordNotFound)

				dbClient.
					ExpectExec(insertStmt).
					WithArgs(
						p.UniqueId,
						p.Path,
						p.Name,
						p.Mimetype,
						p.Extension,
						p.Size,
						p.CreatedAt.UnixMilli(),
						p.CreatedAt.UnixMilli(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(fmt.Errorf("network error"))

				dbClient.
					ExpectRollback().
					WillReturnError(fmt.Errorf("rollback error"))

				res, err := fileRepo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("failed check inserted file", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(gorm.ErrRecordNotFound)

				dbClient.
					ExpectExec(insertStmt).
					WithArgs(
						p.UniqueId,
						p.Path,
						p.Name,
						p.Mimetype,
						p.Extension,
						p.Size,
						p.CreatedAt.UnixMilli(),
						p.CreatedAt.UnixMilli(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(fmt.Errorf("network error"))

				dbClient.
					ExpectRollback()

				res, err := fileRepo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("failed rollback during failure execute callback", func() {
			It("should return error", func() {
				p := repository.CreateFileParam{
					UniqueId:  "id",
					Path:      "storage/id",
					Name:      "dolphin",
					Mimetype:  "image/jpeg",
					Extension: "jpg",
					Size:      2334,
					CreateFn: func(ctx context.Context, p repository.CreateFnParam) error {
						return fmt.Errorf("callback error")
					},
					CreatedAt: currentTs,
				}

				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(gorm.ErrRecordNotFound)

				dbClient.
					ExpectExec(insertStmt).
					WithArgs(
						p.UniqueId,
						p.Path,
						p.Name,
						p.Mimetype,
						p.Extension,
						p.Size,
						p.CreatedAt.UnixMilli(),
						p.CreatedAt.UnixMilli(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				findRows := sqlmock.
					NewRows([]string{
						"id", "name", "path", "mimetype",
						"extension", "size", "created_at",
					}).
					AddRow(
						p.UniqueId,
						p.Name,
						p.Path,
						p.Mimetype,
						p.Extension,
						p.Size,
						p.CreatedAt.UnixMilli(),
					)

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnRows(findRows)

				dbClient.
					ExpectRollback().
					WillReturnError(fmt.Errorf("rollback error"))

				res, err := fileRepo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("failed execute callback", func() {
			It("should return error", func() {
				p := repository.CreateFileParam{
					UniqueId:  "id",
					Path:      "storage/id",
					Name:      "dolphin",
					Mimetype:  "image/jpeg",
					Extension: "jpg",
					Size:      2334,
					CreateFn: func(ctx context.Context, p repository.CreateFnParam) error {
						return fmt.Errorf("callback error")
					},
					CreatedAt: currentTs,
				}

				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(gorm.ErrRecordNotFound)

				dbClient.
					ExpectExec(insertStmt).
					WithArgs(
						p.UniqueId,
						p.Path,
						p.Name,
						p.Mimetype,
						p.Extension,
						p.Size,
						p.CreatedAt.UnixMilli(),
						p.CreatedAt.UnixMilli(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				findRows := sqlmock.
					NewRows([]string{
						"id", "name", "path", "mimetype",
						"extension", "size", "created_at",
					}).
					AddRow(
						p.UniqueId,
						p.Name,
						p.Path,
						p.Mimetype,
						p.Extension,
						p.Size,
						p.CreatedAt.UnixMilli(),
					)

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnRows(findRows)

				dbClient.
					ExpectRollback()

				res, err := fileRepo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("callback error")))
			})
		})

		When("failed execute callback", func() {
			It("should return error", func() {
				p := repository.CreateFileParam{
					UniqueId:  "id",
					Path:      "storage/id",
					Name:      "dolphin",
					Mimetype:  "image/jpeg",
					Extension: "jpg",
					Size:      2334,
					CreateFn: func(ctx context.Context, p repository.CreateFnParam) error {
						return fmt.Errorf("callback error")
					},
					CreatedAt: currentTs,
				}

				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(gorm.ErrRecordNotFound)

				dbClient.
					ExpectExec(insertStmt).
					WithArgs(
						p.UniqueId,
						p.Path,
						p.Name,
						p.Mimetype,
						p.Extension,
						p.Size,
						p.CreatedAt.UnixMilli(),
						p.CreatedAt.UnixMilli(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				findRows := sqlmock.
					NewRows([]string{
						"id", "name", "path", "mimetype",
						"extension", "size", "created_at",
					}).
					AddRow(
						p.UniqueId,
						p.Name,
						p.Path,
						p.Mimetype,
						p.Extension,
						p.Size,
						p.CreatedAt.UnixMilli(),
					)

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnRows(findRows)

				dbClient.
					ExpectRollback()

				res, err := fileRepo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("callback error")))
			})
		})

		When("failed commit trx", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(gorm.ErrRecordNotFound)

				dbClient.
					ExpectExec(insertStmt).
					WithArgs(
						p.UniqueId,
						p.Path,
						p.Name,
						p.Mimetype,
						p.Extension,
						p.Size,
						p.CreatedAt.UnixMilli(),
						p.CreatedAt.UnixMilli(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				findRows := sqlmock.
					NewRows([]string{
						"id", "name", "path", "mimetype",
						"extension", "size", "created_at",
					}).
					AddRow(
						p.UniqueId,
						p.Name,
						p.Path,
						p.Mimetype,
						p.Extension,
						p.Size,
						p.CreatedAt.UnixMilli(),
					)

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnRows(findRows)

				dbClient.
					ExpectCommit().
					WillReturnError(fmt.Errorf("commit error"))

				res, err := fileRepo.CreateFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("commit error")))
			})
		})

		When("success create file", func() {
			It("should return result", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(gorm.ErrRecordNotFound)

				dbClient.
					ExpectExec(insertStmt).
					WithArgs(
						p.UniqueId,
						p.Path,
						p.Name,
						p.Mimetype,
						p.Extension,
						p.Size,
						p.CreatedAt.UnixMilli(),
						p.CreatedAt.UnixMilli(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				findRows := sqlmock.
					NewRows([]string{
						"id", "name", "path", "mimetype",
						"extension", "size", "created_at",
					}).
					AddRow(
						p.UniqueId,
						p.Name,
						p.Path,
						p.Mimetype,
						p.Extension,
						p.Size,
						p.CreatedAt.UnixMilli(),
					)

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnRows(findRows)

				dbClient.
					ExpectCommit()

				res, err := fileRepo.CreateFile(ctx, p)

				Expect(err).To(BeNil())
				Expect(res).To(Equal(&repository.CreateFileResult{
					UniqueId:  p.UniqueId,
					Name:      p.Name,
					Path:      p.Path,
					Mimetype:  p.Mimetype,
					Extension: p.Extension,
					Size:      p.Size,
					CreatedAt: time.UnixMilli(p.CreatedAt.UnixMilli()).UTC(),
				}))
			})
		})
	})

	Context("RetrieveFile function", Label("unit"), func() {
		var (
			ctx       context.Context
			currentTs time.Time
			dbClient  sqlmock.Sqlmock
			fileRepo  repository.File
			p         repository.RetrieveFileParam
			r         *repository.RetrieveFileResult
			findStmt  string
		)

		BeforeEach(func() {
			var (
				db  *sql.DB
				err error
			)

			ctx = context.Background()
			currentTs = time.Now().UTC()
			db, dbClient, err = sqlmock.New()
			if err != nil {
				AbortSuite("failed create db mock: " + err.Error())
			}

			gormClient, err := gorm.Open(gorm_mysql.New(gorm_mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing: true,
			})
			if err != nil {
				AbortSuite("failed create gorm client: " + err.Error())
			}
			fileRepo = repository_mysql.NewFile(repository_mysql.FileParam{
				GormClient: gormClient,
			})

			p = repository.RetrieveFileParam{
				UniqueId: "id",
			}
			r = &repository.RetrieveFileResult{
				UniqueId:  "id",
				CreatedAt: time.UnixMilli(currentTs.UnixMilli()).UTC(),
				DeletedAt: typeconv.Time(time.UnixMilli(currentTs.UnixMilli()).UTC()),
			}
			findStmt = regexp.QuoteMeta("SELECT id, name, path, mimetype, extension, size, created_at, deleted_at FROM `file` WHERE id = ? ORDER BY `file`.`id` LIMIT 1")
		})

		AfterEach(func() {
			err := dbClient.ExpectationsWereMet()
			if err != nil {
				AbortSuite("some expectations were not met " + err.Error())
			}
		})

		When("failed find file", func() {
			It("should return error", func() {
				dbClient.
					ExpectQuery(findStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(fmt.Errorf("network error"))

				res, err := fileRepo.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("file is not available", func() {
			It("should return error", func() {
				dbClient.
					ExpectQuery(findStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnError(gorm.ErrRecordNotFound)

				res, err := fileRepo.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrNotFound))
			})
		})

		When("success find file", func() {
			It("should return result", func() {
				findRows := sqlmock.
					NewRows([]string{
						"id", "name", "path", "mimetype",
						"extension", "size",
						"created_at", "deleted_at",
					}).
					AddRow(
						r.UniqueId,
						r.Name,
						r.Path,
						r.Mimetype,
						r.Extension,
						r.Size,
						currentTs.UnixMilli(),
						currentTs.UnixMilli(),
					)

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(
						p.UniqueId,
					).
					WillReturnRows(findRows)

				res, err := fileRepo.RetrieveFile(ctx, p)

				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("DeleteFile function", Label("unit"), func() {
		var (
			ctx        context.Context
			currentTs  time.Time
			dbClient   sqlmock.Sqlmock
			fileRepo   repository.File
			p          repository.DeleteFileParam
			findStmt   string
			deleteStmt string
			checkStmt  string
			findRows   *sqlmock.Rows
			checkRows  *sqlmock.Rows
		)

		BeforeEach(func() {
			var (
				db  *sql.DB
				err error
			)

			ctx = context.Background()
			currentTs = time.Now().UTC()
			db, dbClient, err = sqlmock.New()
			if err != nil {
				AbortSuite("failed create db mock: " + err.Error())
			}

			gormClient, err := gorm.Open(gorm_mysql.New(gorm_mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing: true,
			})
			if err != nil {
				AbortSuite("failed create gorm client: " + err.Error())
			}
			fileRepo = repository_mysql.NewFile(repository_mysql.FileParam{
				GormClient: gormClient,
			})

			p = repository.DeleteFileParam{
				UniqueId:  "id",
				DeletedAt: currentTs,
				DeleteFn: func(ctx context.Context, p repository.DeleteFnParam) error {
					return nil
				},
			}
			findStmt = regexp.QuoteMeta("SELECT `id`,`deleted_at` FROM `file` WHERE id = ? ORDER BY `file`.`id` LIMIT 1")
			deleteStmt = regexp.QuoteMeta("UPDATE `file` SET `deleted_at`=?,`updated_at`=? WHERE id = ?")
			checkStmt = regexp.QuoteMeta("SELECT id, path, deleted_at FROM `file` WHERE id = ? ORDER BY `file`.`id")
			findRows = sqlmock.
				NewRows([]string{"id", "deleted_at"}).
				AddRow("id", nil)
			checkRows = sqlmock.
				NewRows([]string{"id", "path", "deleted_at"}).
				AddRow("id", "path", currentTs.UnixMilli())
		})

		AfterEach(func() {
			err := dbClient.ExpectationsWereMet()
			if err != nil {
				AbortSuite("some expectations were not met " + err.Error())
			}
		})

		When("failed begin trx", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin().
					WillReturnError(fmt.Errorf("begin error"))

				res, err := fileRepo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("begin error")))
			})
		})

		When("failed rollback during check file existance", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(p.UniqueId).
					WillReturnError(fmt.Errorf("network error"))

				dbClient.
					ExpectRollback().
					WillReturnError(fmt.Errorf("rollback error"))

				res, err := fileRepo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("failed check file existance", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(p.UniqueId).
					WillReturnError(fmt.Errorf("network error"))

				dbClient.
					ExpectRollback()

				res, err := fileRepo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("file is not found", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(p.UniqueId).
					WillReturnError(gorm.ErrRecordNotFound)

				dbClient.
					ExpectRollback()

				res, err := fileRepo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrNotFound))
			})
		})

		When("failed rollback during file deleted", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				findRows := sqlmock.
					NewRows([]string{"id", "deleted_at"}).
					AddRow("id", currentTs.UnixMilli())
				dbClient.
					ExpectQuery(findStmt).
					WithArgs(p.UniqueId).
					WillReturnRows(findRows)

				dbClient.
					ExpectRollback().
					WillReturnError(fmt.Errorf("rollback error"))

				res, err := fileRepo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("file is already deleted", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				findRows := sqlmock.
					NewRows([]string{"id", "deleted_at"}).
					AddRow("id", currentTs.UnixMilli())
				dbClient.
					ExpectQuery(findStmt).
					WithArgs(p.UniqueId).
					WillReturnRows(findRows)

				dbClient.
					ExpectRollback()

				res, err := fileRepo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrDeleted))
			})
		})

		When("failed rollback during delete file", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(p.UniqueId).
					WillReturnRows(findRows)

				dbClient.
					ExpectExec(deleteStmt).
					WithArgs(
						p.DeletedAt.UnixMilli(),
						p.DeletedAt.UnixMilli(),
						p.UniqueId,
					).
					WillReturnError(fmt.Errorf("network error"))

				dbClient.
					ExpectRollback().
					WillReturnError(fmt.Errorf("rollback error"))

				res, err := fileRepo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("failed delete file", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(p.UniqueId).
					WillReturnRows(findRows)

				dbClient.
					ExpectExec(deleteStmt).
					WithArgs(
						p.DeletedAt.UnixMilli(),
						p.DeletedAt.UnixMilli(),
						p.UniqueId,
					).
					WillReturnError(fmt.Errorf("network error"))

				dbClient.
					ExpectRollback()

				res, err := fileRepo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("failed rollback during check deleted file", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(p.UniqueId).
					WillReturnRows(findRows)

				dbClient.
					ExpectExec(deleteStmt).
					WithArgs(
						p.DeletedAt.UnixMilli(),
						p.DeletedAt.UnixMilli(),
						p.UniqueId,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(p.UniqueId).
					WillReturnError(fmt.Errorf("network error"))

				dbClient.
					ExpectRollback().
					WillReturnError(fmt.Errorf("rollback error"))

				res, err := fileRepo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("failed check deleted file", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(p.UniqueId).
					WillReturnRows(findRows)

				dbClient.
					ExpectExec(deleteStmt).
					WithArgs(
						p.DeletedAt.UnixMilli(),
						p.DeletedAt.UnixMilli(),
						p.UniqueId,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(p.UniqueId).
					WillReturnError(fmt.Errorf("network error"))

				dbClient.
					ExpectRollback()

				res, err := fileRepo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("network error")))
			})
		})

		When("failed rollback during execute callback", func() {
			It("should return error", func() {
				p := repository.DeleteFileParam{
					UniqueId:  "id",
					DeletedAt: currentTs,
					DeleteFn: func(ctx context.Context, p repository.DeleteFnParam) error {
						return fmt.Errorf("callback error")
					},
				}

				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(p.UniqueId).
					WillReturnRows(findRows)

				dbClient.
					ExpectExec(deleteStmt).
					WithArgs(
						p.DeletedAt.UnixMilli(),
						p.DeletedAt.UnixMilli(),
						p.UniqueId,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(p.UniqueId).
					WillReturnRows(checkRows)

				dbClient.
					ExpectRollback().
					WillReturnError(fmt.Errorf("rollback error"))

				res, err := fileRepo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("rollback error")))
			})
		})

		When("failed execute callback", func() {
			It("should return error", func() {
				p := repository.DeleteFileParam{
					UniqueId:  "id",
					DeletedAt: currentTs,
					DeleteFn: func(ctx context.Context, p repository.DeleteFnParam) error {
						return fmt.Errorf("callback error")
					},
				}

				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(p.UniqueId).
					WillReturnRows(findRows)

				dbClient.
					ExpectExec(deleteStmt).
					WithArgs(
						p.DeletedAt.UnixMilli(),
						p.DeletedAt.UnixMilli(),
						p.UniqueId,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(p.UniqueId).
					WillReturnRows(checkRows)

				dbClient.
					ExpectRollback()

				res, err := fileRepo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("callback error")))
			})
		})

		When("failed commit trx", func() {
			It("should return error", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(p.UniqueId).
					WillReturnRows(findRows)

				dbClient.
					ExpectExec(deleteStmt).
					WithArgs(
						p.DeletedAt.UnixMilli(),
						p.DeletedAt.UnixMilli(),
						p.UniqueId,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(p.UniqueId).
					WillReturnRows(checkRows)

				dbClient.
					ExpectCommit().
					WillReturnError(fmt.Errorf("commit error"))

				res, err := fileRepo.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("commit error")))
			})
		})

		When("success delete file", func() {
			It("should return result", func() {
				dbClient.
					ExpectBegin()

				dbClient.
					ExpectQuery(findStmt).
					WithArgs(p.UniqueId).
					WillReturnRows(findRows)

				dbClient.
					ExpectExec(deleteStmt).
					WithArgs(
						p.DeletedAt.UnixMilli(),
						p.DeletedAt.UnixMilli(),
						p.UniqueId,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				dbClient.
					ExpectQuery(checkStmt).
					WithArgs(p.UniqueId).
					WillReturnRows(checkRows)

				dbClient.
					ExpectCommit()

				res, err := fileRepo.DeleteFile(ctx, p)

				Expect(err).To(BeNil())
				Expect(res).To(Equal(&repository.DeleteFileResult{
					DeletedAt: time.UnixMilli(currentTs.UnixMilli()).UTC(),
				}))
			})
		})
	})
})
