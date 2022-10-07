package repository_mysql_test

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	mock_datetime "github.com/go-seidon/hippo/internal/datetime/mock"
	"github.com/go-seidon/hippo/internal/repository"
	repository_mysql "github.com/go-seidon/hippo/internal/repository-mysql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth Repository", func() {

	Context("NewAuthRepository function", Label("unit"), func() {
		When("master db client is not specified", func() {
			It("should return error", func() {
				res, err := repository_mysql.NewAuthRepository()

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid db client specified")))
			})
		})

		When("replica db client is not specified", func() {
			It("should return error", func() {
				mOpt := repository_mysql.WithDbMaster(&sql.DB{})
				res, err := repository_mysql.NewAuthRepository(mOpt)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid db client specified")))
			})
		})

		When("required parameter is specified", func() {
			It("should return result", func() {
				mOpt := repository_mysql.WithDbMaster(&sql.DB{})
				rOpt := repository_mysql.WithDbReplica(&sql.DB{})
				res, err := repository_mysql.NewAuthRepository(mOpt, rOpt)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("clock is specified", func() {
			It("should return result", func() {
				clockOpt := repository_mysql.WithClock(&mock_datetime.MockClock{})
				mOpt := repository_mysql.WithDbMaster(&sql.DB{})
				rOpt := repository_mysql.WithDbReplica(&sql.DB{})
				res, err := repository_mysql.NewAuthRepository(clockOpt, mOpt, rOpt)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("FindClient function", Label("unit"), func() {
		var (
			ctx             context.Context
			dbClient        sqlmock.Sqlmock
			repo            repository.AuthRepository
			p               repository.FindClientParam
			findClientQuery string
		)

		BeforeEach(func() {
			ctx = context.Background()

			db, mock, err := sqlmock.New()
			if err != nil {
				AbortSuite("failed create db mock: " + err.Error())
			}
			dbClient = mock

			dbMOpt := repository_mysql.WithDbMaster(db)
			dbROpt := repository_mysql.WithDbReplica(db)
			repo, _ = repository_mysql.NewAuthRepository(dbMOpt, dbROpt)
			p = repository.FindClientParam{
				ClientId: "client_id",
			}

			findClientQuery = regexp.QuoteMeta(`
				SELECT 
					client_id, client_secret
				FROM auth_client
				WHERE client_id = ?
			`)
		})

		When("failed client not found", func() {
			It("should return error", func() {
				dbClient.
					ExpectQuery(findClientQuery).
					WillReturnError(sql.ErrNoRows)

				res, err := repo.FindClient(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(repository.ErrorRecordNotFound))
			})
		})

		When("unexpected error happened", func() {
			It("should return error", func() {
				dbClient.
					ExpectQuery(findClientQuery).
					WillReturnError(fmt.Errorf("error"))

				res, err := repo.FindClient(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("error")))
			})
		})

		When("client is available", func() {
			It("should return result", func() {
				rows := sqlmock.NewRows([]string{
					"client_id", "client_secret",
				}).AddRow(
					"mock-client-id",
					"mock-client-client_secret",
				)
				dbClient.
					ExpectQuery(findClientQuery).
					WillReturnRows(rows)

				res, err := repo.FindClient(ctx, p)

				expectedRes := &repository.FindClientResult{
					ClientId:     "mock-client-id",
					ClientSecret: "mock-client-client_secret",
				}
				Expect(res).To(Equal(expectedRes))
				Expect(err).To(BeNil())
			})
		})
	})

})
