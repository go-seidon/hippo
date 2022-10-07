package db_mysql_test

import (
	"testing"

	db_mysql "github.com/go-seidon/hippo/internal/db-mysql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDeleting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DB MySQL Package")
}

var _ = Describe("DB MySQL Package", func() {
	Context("NewClient function", Label("unit"), func() {
		When("success create client", func() {
			It("should return result", func() {
				res, err := db_mysql.NewClient(
					db_mysql.WithAuth("user", "pw"),
					db_mysql.WithConfig(db_mysql.ClientConfig{
						DbName: "db_name",
					}),
					db_mysql.WithLocation("host", 3306),
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("parse time is used", func() {
			It("should return result", func() {
				res, err := db_mysql.NewClient(
					db_mysql.WithAuth("user", "pw"),
					db_mysql.WithConfig(db_mysql.ClientConfig{
						DbName: "db_name",
					}),
					db_mysql.WithLocation("host", 3306),
					db_mysql.ParseTime(),
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})
})
