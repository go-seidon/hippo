package mysql_test

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-seidon/hippo/internal/repository"
	repository_mysql "github.com/go-seidon/hippo/internal/repository/mysql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth Repository", func() {

	Context("FindClient function", Label("unit"), func() {
		var (
			ctx             context.Context
			dbClient        sqlmock.Sqlmock
			repo            repository.Auth
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
			repo = repository_mysql.NewAuth(dbMOpt, dbROpt)
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
				Expect(err).To(Equal(repository.ErrNotFound))
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
