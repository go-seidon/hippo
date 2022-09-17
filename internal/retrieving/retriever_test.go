package retrieving_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-seidon/local/internal/filesystem"
	mock_filesystem "github.com/go-seidon/local/internal/filesystem/mock"
	mock_logging "github.com/go-seidon/local/internal/logging/mock"
	"github.com/go-seidon/local/internal/repository"
	mock_repository "github.com/go-seidon/local/internal/repository/mock"
	"github.com/go-seidon/local/internal/retrieving"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRetrieving(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Retrieving Package")
}

var _ = Describe("Retriever Service", func() {
	Context("NewRetriever function", Label("unit"), func() {
		var (
			fileRepo    *mock_repository.MockFileRepository
			fileManager *mock_filesystem.MockFileManager
			logger      *mock_logging.MockLogger
			p           retrieving.NewRetrieverParam
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileRepo = mock_repository.NewMockFileRepository(ctrl)
			fileManager = mock_filesystem.NewMockFileManager(ctrl)
			logger = mock_logging.NewMockLogger(ctrl)
			p = retrieving.NewRetrieverParam{
				FileRepo:    fileRepo,
				FileManager: fileManager,
				Logger:      logger,
			}
		})

		When("success create service", func() {
			It("should return result", func() {
				res, err := retrieving.NewRetriever(p)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("file repo is not specified", func() {
			It("should return error", func() {
				p.FileRepo = nil
				res, err := retrieving.NewRetriever(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("file repo is not specified")))
			})
		})

		When("file manager is not specified", func() {
			It("should return error", func() {
				p.FileManager = nil
				res, err := retrieving.NewRetriever(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("file manager is not specified")))
			})
		})

		When("logger is not specified", func() {
			It("should return error", func() {
				p.Logger = nil
				res, err := retrieving.NewRetriever(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("logger is not specified")))
			})
		})
	})

	Context("RetrieveFile function", Label("unit"), func() {
		var (
			ctx           context.Context
			p             retrieving.RetrieveFileParam
			r             *retrieving.RetrieveFileResult
			fileRepo      *mock_repository.MockFileRepository
			fileManager   *mock_filesystem.MockFileManager
			log           *mock_logging.MockLogger
			s             retrieving.Retriever
			retrieveParam repository.RetrieveFileParam
			retrieveRes   *repository.RetrieveFileResult
			openParam     filesystem.OpenFileParam
			openRes       *filesystem.OpenFileResult
		)

		BeforeEach(func() {
			ctx = context.Background()
			p = retrieving.RetrieveFileParam{
				FileId: "mock-file-id",
			}
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileRepo = mock_repository.NewMockFileRepository(ctrl)
			fileManager = mock_filesystem.NewMockFileManager(ctrl)
			log = mock_logging.NewMockLogger(ctrl)
			s, _ = retrieving.NewRetriever(retrieving.NewRetrieverParam{
				FileRepo:    fileRepo,
				FileManager: fileManager,
				Logger:      log,
			})
			retrieveParam = repository.RetrieveFileParam{
				UniqueId: p.FileId,
			}
			retrieveRes = &repository.RetrieveFileResult{
				UniqueId:  p.FileId,
				Name:      "mock-name",
				Path:      "mock-path",
				MimeType:  "mock-mimetype",
				Extension: "mock-extension",
			}
			openParam = filesystem.OpenFileParam{
				Path: retrieveRes.Path,
			}
			file := &os.File{}
			openRes = &filesystem.OpenFileResult{
				File: file,
			}
			r = &retrieving.RetrieveFileResult{
				Data:      file,
				UniqueId:  retrieveRes.UniqueId,
				Name:      retrieveRes.Name,
				Path:      retrieveRes.Path,
				MimeType:  retrieveRes.MimeType,
				Extension: retrieveRes.Extension,
			}

			log.EXPECT().
				Debug("In function: RetrieveFile").
				Times(1)
			log.EXPECT().
				Debug("Returning function: RetrieveFile").
				Times(1)
		})

		When("file id is not specified", func() {
			It("should return error", func() {
				p.FileId = ""
				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid file id parameter")))
			})
		})

		When("file record is not found", func() {
			It("should return error", func() {
				fileRepo.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retrieveParam)).
					Return(nil, repository.ErrorRecordNotFound).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(retrieving.ErrorResourceNotFound))
			})
		})

		When("failed find file record", func() {
			It("should return error", func() {
				fileRepo.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retrieveParam)).
					Return(nil, fmt.Errorf("db error")).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("db error")))
			})
		})

		When("file is not available in disk", func() {
			It("should return error", func() {
				fileRepo.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retrieveParam)).
					Return(retrieveRes, nil).
					Times(1)

				fileManager.
					EXPECT().
					OpenFile(gomock.Eq(ctx), gomock.Eq(openParam)).
					Return(nil, filesystem.ErrorFileNotFound).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(retrieving.ErrorResourceNotFound))
			})
		})

		When("failed open file in disk", func() {
			It("should return error", func() {
				fileRepo.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retrieveParam)).
					Return(retrieveRes, nil).
					Times(1)

				fileManager.
					EXPECT().
					OpenFile(gomock.Eq(ctx), gomock.Eq(openParam)).
					Return(nil, fmt.Errorf("disk error")).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("disk error")))
			})
		})

		When("success retrieve file", func() {
			It("should return result", func() {
				fileRepo.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retrieveParam)).
					Return(retrieveRes, nil).
					Times(1)

				fileManager.
					EXPECT().
					OpenFile(gomock.Eq(ctx), gomock.Eq(openParam)).
					Return(openRes, nil).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(Equal(r))
				Expect(err).To(BeNil())
			})
		})
	})
})
