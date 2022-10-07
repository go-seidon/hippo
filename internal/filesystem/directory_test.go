package filesystem_test

import (
	"context"
	"os"
	"strings"

	"github.com/go-seidon/hippo/internal/filesystem"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Directory Manager", func() {
	Context("NewDirectoryManager function", Label("unit"), func() {
		When("success create directory manager", func() {
			It("should return result", func() {
				res := filesystem.NewDirectoryManager()

				Expect(res).ToNot(BeNil())
			})
		})
	})

	Describe("Directory Manager", Label("integration"), func() {
		var (
			ctx context.Context
			dm  filesystem.DirectoryManager
		)

		BeforeEach(func() {
			ctx = context.Background()
			dm = filesystem.NewDirectoryManager()
		})

		Context("IsDirectoryExists function", func() {
			When("folder is available", func() {
				It("should return result", func() {
					res, err := dm.IsDirectoryExists(ctx, filesystem.IsDirectoryExistsParam{
						Path: "/",
					})

					Expect(res).To(BeTrue())
					Expect(err).To(BeNil())
				})
			})

			When("folder is unavailable", func() {
				It("should return result", func() {
					res, err := dm.IsDirectoryExists(ctx, filesystem.IsDirectoryExistsParam{
						Path: "",
					})

					Expect(res).To(BeFalse())
					Expect(err).To(BeNil())
				})
			})

			When("unexpected error happened", func() {
				It("should return error", func() {
					res, err := dm.IsDirectoryExists(ctx, filesystem.IsDirectoryExistsParam{
						Path: "\000",
					})

					Expect(res).To(BeFalse())
					Expect(err).ToNot(BeNil())
				})
			})
		})

		Context("CreateDir function", Ordered, func() {
			var (
				tempDir string
			)

			BeforeAll(func() {
				wd, _ := os.Getwd()
				tempDir = strings.ReplaceAll(wd, "internal\\filesystem", "") + "\\temp"
			})

			BeforeEach(func() {
				ctx = context.Background()
				dm = filesystem.NewDirectoryManager()
			})

			AfterAll(func() {
				err := os.RemoveAll(tempDir + "\\one")
				if err != nil {
					AbortSuite("failed cleaningup dir: " + err.Error())
				}
			})

			When("failed create dir", func() {
				It("should return error", func() {
					res, err := dm.CreateDir(ctx, filesystem.CreateDirParam{
						Path: "directory.go", //directory.go is exists and it's a file
					})

					Expect(res).To(BeNil())
					Expect(err).ToNot(BeNil())
				})
			})

			When("success create one dir", func() {
				It("should return error", func() {
					res, err := dm.CreateDir(ctx, filesystem.CreateDirParam{
						Path: tempDir + "\\one",
					})

					Expect(res).ToNot(BeNil())
					Expect(err).To(BeNil())
				})
			})

			When("success create nested dir", func() {
				It("should return error", func() {
					res, err := dm.CreateDir(ctx, filesystem.CreateDirParam{
						Path: tempDir + "\\one\\two\\three",
					})

					Expect(res).ToNot(BeNil())
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
