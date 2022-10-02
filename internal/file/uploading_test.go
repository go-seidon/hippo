package file_test

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/go-seidon/local/internal/file"
	mock_file "github.com/go-seidon/local/internal/file/mock"
	"github.com/go-seidon/local/internal/filesystem"
	mock_filesystem "github.com/go-seidon/local/internal/filesystem/mock"
	mock_io "github.com/go-seidon/local/internal/io/mock"
	mock_logging "github.com/go-seidon/local/internal/logging/mock"
	"github.com/go-seidon/local/internal/repository"
	mock_repository "github.com/go-seidon/local/internal/repository/mock"
	mock_text "github.com/go-seidon/local/internal/text/mock"
	mock_validation "github.com/go-seidon/local/internal/validation/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Uploader", func() {

	Context("UploadFile function", Label("unit"), func() {
		var (
			ctx              context.Context
			currentTimestamp time.Time
			fileRepo         *mock_repository.MockFileRepository
			fileManager      *mock_filesystem.MockFileManager
			dirManager       *mock_filesystem.MockDirectoryManager
			logger           *mock_logging.MockLogger
			reader           *mock_io.MockReader
			identifier       *mock_text.MockIdentifier
			locator          *mock_file.MockUploadLocation
			validator        *mock_validation.MockValidator
			s                file.File
			dirExistsParam   filesystem.IsDirectoryExistsParam
			createDirParam   filesystem.CreateDirParam
			createFileRes    *repository.CreateFileResult
			opts             []file.UploadFileOption
		)

		BeforeEach(func() {
			currentTimestamp = time.Now()
			ctx = context.Background()
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileRepo = mock_repository.NewMockFileRepository(ctrl)
			fileManager = mock_filesystem.NewMockFileManager(ctrl)
			dirManager = mock_filesystem.NewMockDirectoryManager(ctrl)
			logger = mock_logging.NewMockLogger(ctrl)
			identifier = mock_text.NewMockIdentifier(ctrl)
			locator = mock_file.NewMockUploadLocation(ctrl)
			validator = mock_validation.NewMockValidator(ctrl)
			reader = mock_io.NewMockReader(ctrl)
			s, _ = file.NewFile(file.NewFileParam{
				FileRepo:    fileRepo,
				FileManager: fileManager,
				DirManager:  dirManager,
				Logger:      logger,
				Identifier:  identifier,
				Locator:     locator,
				Validator:   validator,
				Config: &file.FileConfig{
					UploadDir: "temp",
				},
			})
			dirExistsParam = filesystem.IsDirectoryExistsParam{
				Path: "temp/2022/08/22",
			}
			createDirParam = filesystem.CreateDirParam{
				Path:       "temp/2022/08/22",
				Permission: 0644,
			}
			createFileRes = &repository.CreateFileResult{
				UniqueId:  "mock-unique-id",
				Name:      "mock-name",
				Path:      "mock-path",
				Mimetype:  "mock-mimetype",
				Extension: "mock-extension",
				Size:      200,
				CreatedAt: currentTimestamp,
			}
			dataOpt := file.WithData([]byte{})
			infoOpt := file.WithFileInfo("mock-name", "image/jpeg", "jpg", 100)
			opts = append(opts, dataOpt)
			opts = append(opts, infoOpt)

			logger.
				EXPECT().
				Debug("In function: UploadFile").
				Times(1)
			logger.
				EXPECT().
				Debug("Returning function: UploadFile").
				Times(1)
		})

		When("file data is not specified", func() {
			It("should return error", func() {
				res, err := s.UploadFile(ctx)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("file is not specified")))
			})
		})

		When("failed check directory existance", func() {
			It("should return error", func() {
				locator.
					EXPECT().
					GetLocation().
					Return("2022/08/22").
					Times(1)

				dirManager.
					EXPECT().
					IsDirectoryExists(gomock.Eq(ctx), gomock.Eq(dirExistsParam)).
					Return(false, fmt.Errorf("disk error")).
					Times(1)

				res, err := s.UploadFile(ctx, opts...)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("disk error")))
			})
		})

		When("failed create upload directory", func() {
			It("should return error", func() {
				locator.
					EXPECT().
					GetLocation().
					Return("2022/08/22").
					Times(1)

				dirManager.
					EXPECT().
					IsDirectoryExists(gomock.Eq(ctx), gomock.Eq(dirExistsParam)).
					Return(false, nil).
					Times(1)
				dirManager.
					EXPECT().
					CreateDir(gomock.Eq(ctx), gomock.Eq(createDirParam)).
					Return(nil, fmt.Errorf("r/w error")).
					Times(1)

				res, err := s.UploadFile(ctx, opts...)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("r/w error")))
			})
		})

		When("failed read from file reader", func() {
			It("should return error", func() {
				locator.
					EXPECT().
					GetLocation().
					Return("2022/08/22").
					Times(1)

				dirManager.
					EXPECT().
					IsDirectoryExists(gomock.Eq(ctx), gomock.Eq(dirExistsParam)).
					Return(true, nil).
					Times(1)
				dirManager.
					EXPECT().
					CreateDir(gomock.Eq(ctx), gomock.Eq(createDirParam)).
					Times(0)
				reader.
					EXPECT().
					Read(gomock.Any()).
					Return(0, fmt.Errorf("disk error")).
					Times(1)

				fwOpt := file.WithReader(reader)
				copts := opts
				copts = append(copts, fwOpt)

				res, err := s.UploadFile(ctx, copts...)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("disk error")))
			})
		})

		When("failed generate file id", func() {
			It("should return error", func() {
				locator.
					EXPECT().
					GetLocation().
					Return("2022/08/22").
					Times(1)

				dirManager.
					EXPECT().
					IsDirectoryExists(gomock.Eq(ctx), gomock.Eq(dirExistsParam)).
					Return(true, nil).
					Times(1)
				dirManager.
					EXPECT().
					CreateDir(gomock.Eq(ctx), gomock.Eq(createDirParam)).
					Times(0)
				reader.
					EXPECT().
					Read(gomock.Any()).
					Return(0, io.EOF).
					Times(1)
				identifier.
					EXPECT().
					GenerateId().
					Return("", fmt.Errorf("generate error")).
					Times(1)

				fwOpt := file.WithReader(reader)
				copts := opts
				copts = append(copts, fwOpt)

				res, err := s.UploadFile(ctx, copts...)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("generate error")))
			})
		})

		When("failed create file", func() {
			It("should return error", func() {
				locator.
					EXPECT().
					GetLocation().
					Return("2022/08/22").
					Times(1)

				dirManager.
					EXPECT().
					IsDirectoryExists(gomock.Eq(ctx), gomock.Eq(dirExistsParam)).
					Return(true, nil).
					Times(1)
				dirManager.
					EXPECT().
					CreateDir(gomock.Eq(ctx), gomock.Eq(createDirParam)).
					Times(0)
				identifier.
					EXPECT().
					GenerateId().
					Return("mock-unique-id", nil).
					Times(1)
				fileRepo.
					EXPECT().
					CreateFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, fmt.Errorf("db error")).
					Times(1)

				res, err := s.UploadFile(ctx, opts...)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("db error")))
			})
		})

		When("success upload file", func() {
			It("should return result", func() {
				locator.
					EXPECT().
					GetLocation().
					Return("2022/08/22").
					Times(1)

				dirManager.
					EXPECT().
					IsDirectoryExists(gomock.Eq(ctx), gomock.Eq(dirExistsParam)).
					Return(true, nil).
					Times(1)
				dirManager.
					EXPECT().
					CreateDir(gomock.Eq(ctx), gomock.Eq(createDirParam)).
					Times(0)
				identifier.
					EXPECT().
					GenerateId().
					Return("mock-unique-id", nil).
					Times(1)
				fileRepo.
					EXPECT().
					CreateFile(gomock.Eq(ctx), gomock.Any()).
					Return(createFileRes, nil).
					Times(1)

				res, err := s.UploadFile(ctx, opts...)

				expectedRes := &file.UploadFileResult{
					UniqueId:   "mock-unique-id",
					Name:       "mock-name",
					Path:       "mock-path",
					Mimetype:   "mock-mimetype",
					Extension:  "mock-extension",
					Size:       200,
					UploadedAt: currentTimestamp,
				}
				Expect(res).To(Equal(expectedRes))
				Expect(err).To(BeNil())
			})
		})

	})

	Context("NewCreateFn function", Label("unit"), func() {
		var (
			ctx           context.Context
			data          []byte
			fileManager   *mock_filesystem.MockFileManager
			fn            repository.CreateFn
			createFnParam repository.CreateFnParam
			existsParam   filesystem.IsFileExistsParam
			saveParam     filesystem.SaveFileParam
		)

		BeforeEach(func() {
			ctx = context.Background()
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			data = []byte{}
			fileManager = mock_filesystem.NewMockFileManager(ctrl)
			fn = file.NewCreateFn(data, fileManager)
			createFnParam = repository.CreateFnParam{
				FilePath: "mock/path/name.jpg",
			}
			existsParam = filesystem.IsFileExistsParam{
				Path: createFnParam.FilePath,
			}
			saveParam = filesystem.SaveFileParam{
				Name:       createFnParam.FilePath,
				Data:       data,
				Permission: 0644,
			}
		})

		When("failed check file existance", func() {
			It("should return error", func() {
				fileManager.
					EXPECT().
					IsFileExists(gomock.Eq(ctx), gomock.Eq(existsParam)).
					Return(false, fmt.Errorf("disk error")).
					Times(1)

				err := fn(ctx, createFnParam)

				Expect(err).To(Equal(fmt.Errorf("disk error")))
			})
		})

		When("file already exists", func() {
			It("should return error", func() {
				fileManager.
					EXPECT().
					IsFileExists(gomock.Eq(ctx), gomock.Eq(existsParam)).
					Return(true, nil).
					Times(1)

				err := fn(ctx, createFnParam)

				Expect(err).To(Equal(file.ErrorExists))
			})
		})

		When("failed save file", func() {
			It("should return error", func() {
				fileManager.
					EXPECT().
					IsFileExists(gomock.Eq(ctx), gomock.Eq(existsParam)).
					Return(false, nil).
					Times(1)

				fileManager.
					EXPECT().
					SaveFile(gomock.Eq(ctx), gomock.Eq(saveParam)).
					Return(nil, fmt.Errorf("disk error")).
					Times(1)

				err := fn(ctx, createFnParam)

				Expect(err).To(Equal(fmt.Errorf("disk error")))
			})
		})

		When("success save file", func() {
			It("should return nil", func() {
				fileManager.
					EXPECT().
					IsFileExists(gomock.Eq(ctx), gomock.Eq(existsParam)).
					Return(false, nil).
					Times(1)

				saveRes := filesystem.SaveFileResult{}
				fileManager.
					EXPECT().
					SaveFile(gomock.Eq(ctx), gomock.Eq(saveParam)).
					Return(&saveRes, nil).
					Times(1)

				err := fn(ctx, createFnParam)

				Expect(err).To(BeNil())
			})
		})
	})

})
