package filesystem_test

import (
	"context"
	"io/fs"
	"os"
	"testing"

	"github.com/go-seidon/hippo/internal/filesystem"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFilesystem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filesystem Package")
}

var _ = Describe("File Manager", func() {
	Context("NewFileManager function", Label("unit"), func() {
		When("success create file manager", func() {
			It("should return result", func() {
				res := filesystem.NewFileManager()

				Expect(res).ToNot(BeNil())
			})
		})
	})

	Describe("File Manager", Label("integration"), func() {
		var (
			ctx context.Context
			fm  filesystem.FileManager
		)

		BeforeEach(func() {
			ctx = context.Background()
			fm = filesystem.NewFileManager()
		})

		Context("IsFileExists function", func() {
			When("file is available", func() {
				It("should return result", func() {
					res, err := fm.IsFileExists(ctx, filesystem.IsFileExistsParam{
						Path: "file.go",
					})

					Expect(res).To(BeTrue())
					Expect(err).To(BeNil())
				})
			})

			When("file is unavailable", func() {
				It("should return result", func() {
					res, err := fm.IsFileExists(ctx, filesystem.IsFileExistsParam{
						Path: "unavailable-file",
					})

					Expect(res).To(BeFalse())
					Expect(err).To(BeNil())
				})
			})

			When("unexpected error happened", func() {
				It("should return error", func() {
					res, err := fm.IsFileExists(ctx, filesystem.IsFileExistsParam{
						Path: "\000",
					})

					Expect(res).To(BeFalse())
					Expect(err).ToNot(BeNil())
				})
			})
		})

		Context("OpenFile function", func() {
			When("file is available", func() {
				It("should return result", func() {
					res, err := fm.OpenFile(ctx, filesystem.OpenFileParam{
						Path: "file.go",
					})

					Expect(res).ToNot(BeNil())
					Expect(err).To(BeNil())
				})
			})

			When("file is unavailable", func() {
				It("should return error", func() {
					res, err := fm.OpenFile(ctx, filesystem.OpenFileParam{
						Path: "unavailable-file",
					})

					Expect(res).To(BeNil())
					Expect(err).To(Equal(filesystem.ErrorFileNotFound))
				})
			})

			When("failed open file", func() {
				It("should return error", func() {
					res, err := fm.OpenFile(ctx, filesystem.OpenFileParam{
						Path: "\000",
					})

					Expect(res).To(BeNil())
					Expect(err).ToNot(BeNil())
				})
			})
		})

		Context("SaveFile function", Ordered, func() {
			var (
				fileName string
			)

			BeforeAll(func() {
				fileName = "temp-save-file.txt"
			})

			AfterAll(func() {
				err := os.Remove(fileName)
				if err != nil {
					AbortSuite("failed cleaningup temp file: " + err.Error())
				}
			})

			When("failed save file", func() {
				It("should return error", func() {
					res, err := fm.SaveFile(ctx, filesystem.SaveFileParam{
						Name:       "", //should specify file name
						Data:       nil,
						Permission: fs.ModeTemporary,
					})

					Expect(res).To(BeNil())
					Expect(err).ToNot(BeNil())
				})
			})

			When("success save file", func() {
				It("should return result", func() {
					res, err := fm.SaveFile(ctx, filesystem.SaveFileParam{
						Name:       fileName,
						Data:       nil,
						Permission: fs.ModeTemporary,
					})

					Expect(res).ToNot(BeNil())
					Expect(err).To(BeNil())
				})
			})
		})

		Context("RemoveFile function", Ordered, func() {
			var (
				fileName string
			)

			BeforeAll(func() {
				fileName = "temp-remove-file.txt"
				err := os.WriteFile(fileName, nil, fs.ModeTemporary)
				if err != nil {
					AbortSuite("failed settingup temp file: " + err.Error())
				}
			})

			When("failed remove file", func() {
				It("should return error", func() {
					res, err := fm.RemoveFile(ctx, filesystem.RemoveFileParam{
						Path: "\000",
					})

					Expect(res).To(BeNil())
					Expect(err).ToNot(BeNil())
				})
			})

			When("file is unavailable", func() {
				It("should return result", func() {
					res, err := fm.RemoveFile(ctx, filesystem.RemoveFileParam{
						Path: "unavailable-file",
					})

					Expect(res).To(BeNil())
					Expect(err).To(Equal(filesystem.ErrorFileNotFound))
				})
			})

			When("file is available", func() {
				It("should return result", func() {
					res, err := fm.RemoveFile(ctx, filesystem.RemoveFileParam{
						Path: fileName,
					})

					Expect(res).ToNot(BeNil())
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
