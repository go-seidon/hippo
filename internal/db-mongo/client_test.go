package db_mongo_test

import (
	"fmt"
	"testing"

	db_mongo "github.com/go-seidon/hippo/internal/db-mongo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDeleting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DB MySQL Package")
}

var _ = Describe("DB MySQL Package", func() {
	Context("NewClient function", Label("unit"), func() {
		When("db mode is not valid", func() {
			It("should return error", func() {
				res, err := db_mongo.NewClient(
					db_mongo.WithConfig(db_mongo.ClientConfig{
						DbName: "db_name",
						DbMode: "invalid",
					}),
				)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("mode is not supported")))
			})
		})

		When("auth mode is not valid", func() {
			It("should return error", func() {
				res, err := db_mongo.NewClient(
					db_mongo.WithConfig(db_mongo.ClientConfig{
						DbName:   "db_name",
						AuthMode: "invalid",
					}),
					db_mongo.UsingStandalone("host", 27017),
				)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("auth is not supported")))
			})
		})

		When("success create standalone client", func() {
			It("should return result", func() {
				res, err := db_mongo.NewClient(
					db_mongo.WithBasicAuth("user", "pw", "db_name"),
					db_mongo.WithConfig(db_mongo.ClientConfig{
						DbName: "db_name",
					}),
					db_mongo.UsingStandalone("host", 27017),
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("success create replica client", func() {
			It("should return result", func() {
				res, err := db_mongo.NewClient(
					db_mongo.WithBasicAuth("user", "pw", "db_name"),
					db_mongo.WithConfig(db_mongo.ClientConfig{
						DbName: "db_name",
					}),
					db_mongo.UsingReplication("rs-name", []string{"h1:27030", "h2:27031"}),
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})
})
