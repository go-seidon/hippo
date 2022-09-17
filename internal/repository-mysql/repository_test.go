package repository_mysql_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/go-seidon/local/internal/repository"
	repository_mysql "github.com/go-seidon/local/internal/repository-mysql"
)

func TestRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Package")
}

var _ = Describe("Repository Provider", func() {
	Context("NewRepository function", Label("unit"), func() {
		When("master db client is not specified", func() {
			It("should return error", func() {
				res, err := repository_mysql.NewRepository()

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid db client specified")))
			})
		})

		When("replica db client is not specified", func() {
			It("should return error", func() {
				mOpt := repository_mysql.WithDbMaster(&sql.DB{})
				res, err := repository_mysql.NewRepository(mOpt)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid db client specified")))
			})
		})

		When("required parameters are specified", func() {
			It("should return result", func() {
				mOpt := repository_mysql.WithDbMaster(&sql.DB{})
				rOpt := repository_mysql.WithDbReplica(&sql.DB{})
				res, err := repository_mysql.NewRepository(mOpt, rOpt)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("GetAuthRepo function", Label("unit"), func() {
		var (
			provider repository.Provider
		)

		BeforeEach(func() {
			mOpt := repository_mysql.WithDbMaster(&sql.DB{})
			rOpt := repository_mysql.WithDbReplica(&sql.DB{})
			provider, _ = repository_mysql.NewRepository(mOpt, rOpt)
		})

		When("function is called", func() {
			It("should return result", func() {
				res := provider.GetAuthRepo()

				Expect(res).ToNot(BeNil())
			})
		})
	})

	Context("GetFileRepo function", Label("unit"), func() {
		var (
			provider repository.Provider
		)

		BeforeEach(func() {
			mOpt := repository_mysql.WithDbMaster(&sql.DB{})
			rOpt := repository_mysql.WithDbReplica(&sql.DB{})
			provider, _ = repository_mysql.NewRepository(mOpt, rOpt)
		})

		When("function is called", func() {
			It("should return result", func() {
				res := provider.GetFileRepo()

				Expect(res).ToNot(BeNil())
			})
		})
	})

	Context("Init function", Label("unit"), func() {
		var (
			provider repository.Provider
			ctx      context.Context
		)

		BeforeEach(func() {
			mOpt := repository_mysql.WithDbMaster(&sql.DB{})
			rOpt := repository_mysql.WithDbReplica(&sql.DB{})
			provider, _ = repository_mysql.NewRepository(mOpt, rOpt)
			ctx = context.Background()
		})

		When("success init", func() {
			It("should return result", func() {
				res := provider.Init(ctx)

				Expect(res).To(BeNil())
			})
		})
	})
})
