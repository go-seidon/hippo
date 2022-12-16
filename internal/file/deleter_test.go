package file_test

import (
	"context"
	"fmt"
	"time"

	"github.com/go-seidon/hippo/internal/file"
	mock_file "github.com/go-seidon/hippo/internal/file/mock"
	"github.com/go-seidon/hippo/internal/filesystem"
	mock_filesystem "github.com/go-seidon/hippo/internal/filesystem/mock"
	"github.com/go-seidon/hippo/internal/repository"
	mock_repository "github.com/go-seidon/hippo/internal/repository/mock"
	mock_identifier "github.com/go-seidon/provider/identifier/mock"
	mock_logging "github.com/go-seidon/provider/logging/mock"
	mock_validation "github.com/go-seidon/provider/validation/mock"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deleter", func() {

	Context("DeleteFile function", Label("unit"), func() {
		var (
			ctx         context.Context
			p           file.DeleteFileParam
			fileRepo    *mock_repository.MockFileRepository
			fileManager *mock_filesystem.MockFileManager
			log         *mock_logging.MockLogger
			validator   *mock_validation.MockValidator
			s           file.File
			deleteRes   *repository.DeleteFileResult
			finalRes    *file.DeleteFileResult
		)

		BeforeEach(func() {
			currentTimestamp := time.Now()
			ctx = context.Background()
			p = file.DeleteFileParam{
				FileId: "mock-file-id",
			}
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileRepo = mock_repository.NewMockFileRepository(ctrl)
			fileManager = mock_filesystem.NewMockFileManager(ctrl)
			dirManager := mock_filesystem.NewMockDirectoryManager(ctrl)
			identifier := mock_identifier.NewMockIdentifier(ctrl)
			locator := mock_file.NewMockUploadLocation(ctrl)
			log = mock_logging.NewMockLogger(ctrl)
			validator = mock_validation.NewMockValidator(ctrl)
			s, _ = file.NewFile(file.NewFileParam{
				FileRepo:    fileRepo,
				FileManager: fileManager,
				DirManager:  dirManager,
				Logger:      log,
				Identifier:  identifier,
				Locator:     locator,
				Validator:   validator,
				Config: &file.FileConfig{
					UploadDir: "temp",
				},
			})
			deleteRes = &repository.DeleteFileResult{
				DeletedAt: currentTimestamp,
			}
			finalRes = &file.DeleteFileResult{
				DeletedAt: currentTimestamp,
			}

			log.EXPECT().
				Debug("In function: DeleteFile").
				Times(1)
			log.EXPECT().
				Debug("Returning function: DeleteFile").
				Times(1)
		})

		When("parameter are not valid", func() {
			It("should return error", func() {
				p.FileId = ""

				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(fmt.Errorf("invalid data")).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid data")))
			})
		})

		When("failed delete file", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

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
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

				fileRepo.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, repository.ErrNotFound).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(file.ErrorNotFound))
			})
		})

		When("file is deleted", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

				fileRepo.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, repository.ErrDeleted).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(file.ErrorNotFound))
			})
		})

		When("failed success file", func() {
			It("should return result", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

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
			fn = file.NewDeleteFn(fileManager)
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

				Expect(err).To(Equal(file.ErrorNotFound))
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
