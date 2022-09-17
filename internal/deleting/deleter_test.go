package deleting_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-seidon/local/internal/deleting"
	"github.com/go-seidon/local/internal/filesystem"
	mock_filesystem "github.com/go-seidon/local/internal/filesystem/mock"
	mock_logging "github.com/go-seidon/local/internal/logging/mock"
	"github.com/go-seidon/local/internal/repository"
	mock_repository "github.com/go-seidon/local/internal/repository/mock"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDeleting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Deleting Package")
}

var _ = Describe("Deleter Service", func() {
	Context("NewDeleter function", Label("unit"), func() {
		var (
			fileRepo    *mock_repository.MockFileRepository
			fileManager *mock_filesystem.MockFileManager
			logger      *mock_logging.MockLogger
			p           deleting.NewDeleterParam
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileRepo = mock_repository.NewMockFileRepository(ctrl)
			fileManager = mock_filesystem.NewMockFileManager(ctrl)
			logger = mock_logging.NewMockLogger(ctrl)
			p = deleting.NewDeleterParam{
				FileRepo:    fileRepo,
				FileManager: fileManager,
				Logger:      logger,
			}
		})

		When("success create service", func() {
			It("should return result", func() {
				res, err := deleting.NewDeleter(p)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("file repo is not specified", func() {
			It("should return error", func() {
				p.FileRepo = nil
				res, err := deleting.NewDeleter(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("file repo is not specified")))
			})
		})

		When("file manager is not specified", func() {
			It("should return error", func() {
				p.FileManager = nil
				res, err := deleting.NewDeleter(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("file manager is not specified")))
			})
		})

		When("logger is not specified", func() {
			It("should return error", func() {
				p.Logger = nil
				res, err := deleting.NewDeleter(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("logger is not specified")))
			})
		})
	})

	Context("DeleteFile function", Label("unit"), func() {
		var (
			ctx         context.Context
			p           deleting.DeleteFileParam
			fileRepo    *mock_repository.MockFileRepository
			fileManager *mock_filesystem.MockFileManager
			log         *mock_logging.MockLogger
			s           deleting.Deleter
			deleteRes   *repository.DeleteFileResult
			finalRes    *deleting.DeleteFileResult
		)

		BeforeEach(func() {
			currentTimestamp := time.Now()
			ctx = context.Background()
			p = deleting.DeleteFileParam{
				FileId: "mock-file-id",
			}
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileRepo = mock_repository.NewMockFileRepository(ctrl)
			fileManager = mock_filesystem.NewMockFileManager(ctrl)
			log = mock_logging.NewMockLogger(ctrl)
			s, _ = deleting.NewDeleter(deleting.NewDeleterParam{
				FileRepo:    fileRepo,
				FileManager: fileManager,
				Logger:      log,
			})
			deleteRes = &repository.DeleteFileResult{
				DeletedAt: currentTimestamp,
			}
			finalRes = &deleting.DeleteFileResult{
				DeletedAt: currentTimestamp,
			}

			log.EXPECT().
				Debug("In function: DeleteFile").
				Times(1)
			log.EXPECT().
				Debug("Returning function: DeleteFile").
				Times(1)
		})

		When("file id is not specified", func() {
			It("should return error", func() {
				p.FileId = ""
				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid file id parameter")))
			})
		})

		When("failed delete file", func() {
			It("should return error", func() {
				fileRepo.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, fmt.Errorf("failed delete file")).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed delete file")))
			})
		})

		When("file is not available", func() {
			It("should return error", func() {
				fileRepo.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, repository.ErrorRecordNotFound).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(deleting.ErrorResourceNotFound))
			})
		})

		When("failed success file", func() {
			It("should return result", func() {
				fileRepo.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Any()).
					Return(deleteRes, nil).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(Equal(finalRes))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("NewDeleteFn function", Label("unit"), func() {
		var (
			ctx               context.Context
			fileManager       *mock_filesystem.MockFileManager
			fn                repository.DeleteFn
			deleteFnParam     repository.DeleteFnParam
			isFileExistsParam filesystem.IsFileExistsParam
			removeParam       filesystem.RemoveFileParam
			removeRes         *filesystem.RemoveFileResult
		)

		BeforeEach(func() {
			currentTimestamp := time.Now()
			ctx = context.Background()
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileManager = mock_filesystem.NewMockFileManager(ctrl)
			fn = deleting.NewDeleteFn(fileManager)
			deleteFnParam = repository.DeleteFnParam{
				FilePath: "mock/path",
			}
			isFileExistsParam = filesystem.IsFileExistsParam{
				Path: deleteFnParam.FilePath,
			}
			removeParam = filesystem.RemoveFileParam{
				Path: deleteFnParam.FilePath,
			}
			removeRes = &filesystem.RemoveFileResult{
				RemovedAt: currentTimestamp,
			}
		})

		When("failed check file existstance", func() {
			It("should return error", func() {
				fileManager.
					EXPECT().
					IsFileExists(gomock.Eq(ctx), gomock.Eq(isFileExistsParam)).
					Return(false, fmt.Errorf("failed read disk")).
					Times(1)

				err := fn(ctx, deleteFnParam)

				Expect(err).To(Equal(fmt.Errorf("failed read disk")))
			})
		})

		When("file is not available in disk", func() {
			It("should return error", func() {
				fileManager.
					EXPECT().
					IsFileExists(gomock.Eq(ctx), gomock.Eq(isFileExistsParam)).
					Return(false, nil).
					Times(1)

				err := fn(ctx, deleteFnParam)

				Expect(err).To(Equal(deleting.ErrorResourceNotFound))
			})
		})

		When("failed remove file from disk", func() {
			It("should return error", func() {
				fileManager.
					EXPECT().
					IsFileExists(gomock.Eq(ctx), gomock.Eq(isFileExistsParam)).
					Return(true, nil).
					Times(1)

				fileManager.
					EXPECT().
					RemoveFile(gomock.Eq(ctx), gomock.Eq(removeParam)).
					Return(nil, fmt.Errorf("disk error")).
					Times(1)

				err := fn(ctx, deleteFnParam)

				Expect(err).To(Equal(fmt.Errorf("disk error")))
			})
		})

		When("success remove file from disk", func() {
			It("should return result", func() {
				fileManager.
					EXPECT().
					IsFileExists(gomock.Eq(ctx), gomock.Eq(isFileExistsParam)).
					Return(true, nil).
					Times(1)

				fileManager.
					EXPECT().
					RemoveFile(gomock.Eq(ctx), gomock.Eq(removeParam)).
					Return(removeRes, nil).
					Times(1)

				err := fn(ctx, deleteFnParam)

				Expect(err).To(BeNil())
			})
		})
	})
})
