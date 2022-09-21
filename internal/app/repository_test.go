package app_test

import (
	"fmt"

	"github.com/go-seidon/local/internal/app"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository Package", func() {

	Context("NewDefaultRepository function", Label("unit"), func() {
		When("config is not specified", func() {
			It("should return error", func() {
				res, err := app.NewDefaultRepository(nil)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid config")))
			})
		})

		When("db provider is not valid", func() {
			It("should return error", func() {
				res, err := app.NewDefaultRepository(&app.Config{
					DBProvider: "invalid",
				})

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid repository provider")))
			})
		})

		Context("mysql repository", func() {
			When("success create repository", func() {
				It("should return result", func() {
					res, err := app.NewDefaultRepository(&app.Config{
						DBProvider: "mysql",
					})

					Expect(res).ToNot(BeNil())
					Expect(err).To(BeNil())
				})
			})
		})

		Context("mongo repository", func() {
			When("db mode is not valid", func() {
				It("should return error", func() {
					res, err := app.NewDefaultRepository(&app.Config{
						DBProvider: "mongo",
						MongoMode:  "invalid",
					})

					Expect(res).To(BeNil())
					Expect(err).ToNot(BeNil())
				})
			})

			When("auth is not valid", func() {
				It("should return error", func() {
					res, err := app.NewDefaultRepository(&app.Config{
						DBProvider: "mongo",
						MongoMode:  "standalone",
					})

					Expect(res).To(BeNil())
					Expect(err).ToNot(BeNil())
				})
			})

			When("success create using standalone", func() {
				It("should return error", func() {
					res, err := app.NewDefaultRepository(&app.Config{
						DBProvider:    "mongo",
						MongoMode:     "standalone",
						MongoAuthMode: "basic",
					})

					Expect(res).ToNot(BeNil())
					Expect(err).To(BeNil())
				})
			})

			When("success create using replication", func() {
				It("should return error", func() {
					res, err := app.NewDefaultRepository(&app.Config{
						DBProvider:    "mongo",
						MongoMode:     "replication",
						MongoAuthMode: "basic",
					})

					Expect(res).ToNot(BeNil())
					Expect(err).To(BeNil())
				})
			})
		})
	})

})
