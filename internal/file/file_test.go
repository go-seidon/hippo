package file_test

import (
	"fmt"
	"testing"

	"github.com/go-seidon/local/internal/file"
	mock_filesystem "github.com/go-seidon/local/internal/filesystem/mock"
	mock_logging "github.com/go-seidon/local/internal/logging/mock"
	mock_repository "github.com/go-seidon/local/internal/repository/mock"
	mock_text "github.com/go-seidon/local/internal/text/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "File Package")
}

var _ = Describe("File", func() {

	Context("NewFile function", Label("unit"), func() {
		var (
			fileRepo    *mock_repository.MockFileRepository
			fileManager *mock_filesystem.MockFileManager
			dirManager  *mock_filesystem.MockDirectoryManager
			logger      *mock_logging.MockLogger
			identifier  *mock_text.MockIdentifier
			p           file.NewFileParam
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileRepo = mock_repository.NewMockFileRepository(ctrl)
			fileManager = mock_filesystem.NewMockFileManager(ctrl)
			dirManager = mock_filesystem.NewMockDirectoryManager(ctrl)
			logger = mock_logging.NewMockLogger(ctrl)
			identifier = mock_text.NewMockIdentifier(ctrl)
			p = file.NewFileParam{
				FileRepo:    fileRepo,
				FileManager: fileManager,
				DirManager:  dirManager,
				Logger:      logger,
				Identifier:  identifier,
			}
		})

		When("success create service", func() {
			It("should return result", func() {
				res, err := file.NewFile(p)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("file repo is not specified", func() {
			It("should return error", func() {
				p.FileRepo = nil
				res, err := file.NewFile(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("file repo is not specified")))
			})
		})

		When("file manager is not specified", func() {
			It("should return error", func() {
				p.FileManager = nil
				res, err := file.NewFile(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("file manager is not specified")))
			})
		})

		When("directory manager is not specified", func() {
			It("should return error", func() {
				p.DirManager = nil
				res, err := file.NewFile(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("directory manager is not specified")))
			})
		})

		When("logger is not specified", func() {
			It("should return error", func() {
				p.Logger = nil
				res, err := file.NewFile(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("logger is not specified")))
			})
		})

		When("identifier is not specified", func() {
			It("should return error", func() {
				p.Identifier = nil
				res, err := file.NewFile(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("identifier is not specified")))
			})
		})
	})

})
