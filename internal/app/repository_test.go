package app_test

import (
	"fmt"

	"github.com/go-seidon/local/internal/app"
	"github.com/go-seidon/local/internal/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/go-sql-driver/mysql"
)

var _ = Describe("Repository Package", func() {
	Context("NewRepository function", Label("unit"), func() {
		var (
			opt *mock.MockRepositoryOption
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			opt = mock.NewMockRepositoryOption(ctrl)
		})

		When("option is not specified", func() {
			It("should return error", func() {
				res, err := app.NewRepository()

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid repository option")))
			})
		})

		When("provider is not supported", func() {
			It("should return error", func() {
				opt.EXPECT().Apply(&app.NewRepositoryOption{})
				res, err := app.NewRepository(opt)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("db provider is not supported")))
			})
		})

		When("success create mysql repository", func() {
			It("should return result", func() {
				masterOpt := app.WithMySQL(app.MySQLConn{
					Host:     "mock-host",
					Port:     3306,
					User:     "mock-username",
					Password: "mock-password",
					DbName:   "mock-db",
				}, app.MySQLConn{
					Host:     "mock-host",
					Port:     3306,
					User:     "mock-username",
					Password: "mock-password",
					DbName:   "mock-db",
				})
				res, err := app.NewRepository(masterOpt)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("success create mongo repository", func() {
			It("should return result", func() {
				dbClientOpt := app.WithMongo(app.MongoConn{
					Host:     "mock-host",
					Port:     27017,
					User:     "mock-user",
					Password: "mock-password",
					DbName:   "mock-db",
				})
				res, err := app.NewRepository(dbClientOpt)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})
})
