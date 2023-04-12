package service_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-seidon/hippo/internal/file"
	mock_file "github.com/go-seidon/hippo/internal/file/mock"
	"github.com/go-seidon/hippo/internal/filesystem"
	mock_filesystem "github.com/go-seidon/hippo/internal/filesystem/mock"
	"github.com/go-seidon/hippo/internal/repository"
	mock_repository "github.com/go-seidon/hippo/internal/repository/mock"
	"github.com/go-seidon/hippo/internal/service"
	mock_datetime "github.com/go-seidon/provider/datetime/mock"
	mock_identifier "github.com/go-seidon/provider/identity/mock"
	mock_io "github.com/go-seidon/provider/io/mock"
	mock_logging "github.com/go-seidon/provider/logging/mock"
	"github.com/go-seidon/provider/system"
	"github.com/go-seidon/provider/typeconv"
	mock_validation "github.com/go-seidon/provider/validation/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("File Service", func() {
	Context("RetrieveFile function", Label("unit"), func() {
		var (
			ctx           context.Context
			currentTs     time.Time
			p             service.RetrieveFileParam
			r             *service.RetrieveFileResult
			fileRepo      *mock_repository.MockFile
			fileManager   *mock_filesystem.MockFileManager
			log           *mock_logging.MockLogger
			validator     *mock_validation.MockValidator
			s             service.File
			retrieveParam repository.RetrieveFileParam
			retrieveRes   *repository.RetrieveFileResult
			openParam     filesystem.OpenFileParam
			openRes       *filesystem.OpenFileResult
		)

		BeforeEach(func() {
			ctx = context.Background()
			currentTs = time.Now().UTC()
			p = service.RetrieveFileParam{
				FileId: "mock-file-id",
			}
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileRepo = mock_repository.NewMockFile(ctrl)
			fileManager = mock_filesystem.NewMockFileManager(ctrl)
			dirManager := mock_filesystem.NewMockDirectoryManager(ctrl)
			identifier := mock_identifier.NewMockIdentifier(ctrl)
			locator := mock_file.NewMockUploadLocation(ctrl)
			log = mock_logging.NewMockLogger(ctrl)
			validator = mock_validation.NewMockValidator(ctrl)
			s = service.NewFile(service.FileParam{
				FileRepo:    fileRepo,
				FileManager: fileManager,
				DirManager:  dirManager,
				Identifier:  identifier,
				Logger:      log,
				Locator:     locator,
				Validator:   validator,
				Config: &service.FileConfig{
					UploadDir: "temp",
				},
			})
			retrieveParam = repository.RetrieveFileParam{
				UniqueId: p.FileId,
			}
			retrieveRes = &repository.RetrieveFileResult{
				UniqueId:  p.FileId,
				Name:      "mock-name",
				Path:      "mock-path",
				Mimetype:  "mock-mimetype",
				Extension: "mock-extension",
			}
			openParam = filesystem.OpenFileParam{
				Path: retrieveRes.Path,
			}
			osFile := &os.File{}
			openRes = &filesystem.OpenFileResult{
				File: osFile,
			}
			r = &service.RetrieveFileResult{
				Success: system.Success{
					Code:    1000,
					Message: "success retrieve file",
				},
				Data:      osFile,
				UniqueId:  retrieveRes.UniqueId,
				Name:      retrieveRes.Name,
				Path:      retrieveRes.Path,
				MimeType:  retrieveRes.Mimetype,
				Extension: retrieveRes.Extension,
			}

			log.
				EXPECT().
				Debug("In function: RetrieveFile").
				Times(1)
			log.
				EXPECT().
				Debug("Returning function: RetrieveFile").
				Times(1)
		})

		When("parameter is not valid", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(fmt.Errorf("invalid data")).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Code).To(Equal(int32(1002)))
				Expect(err.Message).To(Equal("invalid data"))
			})
		})

		When("file record is not found", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

				fileRepo.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retrieveParam)).
					Return(nil, repository.ErrNotFound).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Code).To(Equal(int32(1004)))
				Expect(err.Message).To(Equal("file is not found"))
			})
		})

		When("file record is deleted", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

				retrieveRes := &repository.RetrieveFileResult{
					UniqueId:  p.FileId,
					Name:      "mock-name",
					Path:      "mock-path",
					Mimetype:  "mock-mimetype",
					Extension: "mock-extension",
					DeletedAt: typeconv.Time(currentTs),
				}
				fileRepo.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retrieveParam)).
					Return(retrieveRes, nil).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Code).To(Equal(int32(1004)))
				Expect(err.Message).To(Equal("file is deleted"))
			})
		})

		When("failed find file record", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

				fileRepo.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retrieveParam)).
					Return(nil, fmt.Errorf("db error")).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Code).To(Equal(int32(1001)))
				Expect(err.Message).To(Equal("db error"))
			})
		})

		When("file is not available in disk", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

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
				Expect(err.Code).To(Equal(int32(1004)))
				Expect(err.Message).To(Equal("file is not found"))
			})
		})

		When("failed open file in disk", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

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
				Expect(err.Code).To(Equal(int32(1001)))
				Expect(err.Message).To(Equal("disk error"))
			})
		})

		When("success retrieve file", func() {
			It("should return result", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

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

	Context("UploadFile function", Label("unit"), func() {
		var (
			ctx            context.Context
			currentTs      time.Time
			fileRepo       *mock_repository.MockFile
			fileManager    *mock_filesystem.MockFileManager
			dirManager     *mock_filesystem.MockDirectoryManager
			logger         *mock_logging.MockLogger
			reader         *mock_io.MockReader
			identifier     *mock_identifier.MockIdentifier
			clock          *mock_datetime.MockClock
			locator        *mock_file.MockUploadLocation
			validator      *mock_validation.MockValidator
			s              service.File
			dirExistsParam filesystem.IsDirectoryExistsParam
			createDirParam filesystem.CreateDirParam
			createFileRes  *repository.CreateFileResult
			opts           []service.UploadFileOption
			r              *service.UploadFileResult
		)

		BeforeEach(func() {
			currentTs = time.Now().UTC()
			ctx = context.Background()
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileRepo = mock_repository.NewMockFile(ctrl)
			fileManager = mock_filesystem.NewMockFileManager(ctrl)
			dirManager = mock_filesystem.NewMockDirectoryManager(ctrl)
			logger = mock_logging.NewMockLogger(ctrl)
			identifier = mock_identifier.NewMockIdentifier(ctrl)
			clock = mock_datetime.NewMockClock(ctrl)
			locator = mock_file.NewMockUploadLocation(ctrl)
			validator = mock_validation.NewMockValidator(ctrl)
			reader = mock_io.NewMockReader(ctrl)
			s = service.NewFile(service.FileParam{
				FileRepo:    fileRepo,
				FileManager: fileManager,
				DirManager:  dirManager,
				Logger:      logger,
				Identifier:  identifier,
				Clock:       clock,
				Locator:     locator,
				Validator:   validator,
				Config: &service.FileConfig{
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
				CreatedAt: currentTs,
			}
			dataOpt := service.WithReader(reader)
			infoOpt := service.WithFileInfo("mock-name", "image/jpeg", "jpg", 100)
			opts = append(opts, dataOpt)
			opts = append(opts, infoOpt)
			r = &service.UploadFileResult{
				Success: system.Success{
					Code:    1000,
					Message: "success upload file",
				},
				UniqueId:   "mock-unique-id",
				Name:       "mock-name",
				Path:       "mock-path",
				Mimetype:   "mock-mimetype",
				Extension:  "mock-extension",
				Size:       200,
				UploadedAt: currentTs,
			}

			logger.
				EXPECT().
				Debug("In function: UploadFile").
				Times(1)
			logger.
				EXPECT().
				Debug("Returning function: UploadFile").
				Times(1)
		})

		When("there is invalid data", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Any()).
					Return(fmt.Errorf("invalid data")).
					Times(1)

				res, err := s.UploadFile(ctx)

				Expect(res).To(BeNil())
				Expect(err.Code).To(Equal(int32(1002)))
				Expect(err.Message).To(Equal("invalid data"))
			})
		})

		When("file is not specified", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Any()).
					Return(nil).
					Times(1)

				res, err := s.UploadFile(ctx)

				Expect(res).To(BeNil())
				Expect(err.Code).To(Equal(int32(1002)))
				Expect(err.Message).To(Equal("file is not specified"))
			})
		})

		When("failed check directory existance", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Any()).
					Return(nil).
					Times(1)

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
				Expect(err.Code).To(Equal(int32(1001)))
				Expect(err.Message).To(Equal("disk error"))
			})
		})

		When("failed create upload directory", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Any()).
					Return(nil).
					Times(1)

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
					Return(nil, fmt.Errorf("i/o error")).
					Times(1)

				res, err := s.UploadFile(ctx, opts...)

				Expect(res).To(BeNil())
				Expect(err.Code).To(Equal(int32(1001)))
				Expect(err.Message).To(Equal("i/o error"))
			})
		})

		When("failed read from file reader", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Any()).
					Return(nil).
					Times(1)

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

				fwOpt := service.WithReader(reader)
				copts := opts
				copts = append(copts, fwOpt)

				res, err := s.UploadFile(ctx, copts...)

				Expect(res).To(BeNil())
				Expect(err.Code).To(Equal(int32(1001)))
				Expect(err.Message).To(Equal("disk error"))
			})
		})

		When("failed generate file id", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Any()).
					Return(nil).
					Times(1)

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

				fwOpt := service.WithReader(reader)
				copts := opts
				copts = append(copts, fwOpt)

				res, err := s.UploadFile(ctx, copts...)

				Expect(res).To(BeNil())
				Expect(err.Code).To(Equal(int32(1001)))
				Expect(err.Message).To(Equal("generate error"))
			})
		})

		When("failed write file record", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Any()).
					Return(nil).
					Times(1)

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
					Return("mock-unique-id", nil).
					Times(1)

				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				fileRepo.
					EXPECT().
					CreateFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, fmt.Errorf("db error")).
					Times(1)

				res, err := s.UploadFile(ctx, opts...)

				Expect(res).To(BeNil())
				Expect(err.Code).To(Equal(int32(1001)))
				Expect(err.Message).To(Equal("db error"))
			})
		})

		When("success upload file", func() {
			It("should return result", func() {
				validator.
					EXPECT().
					Validate(gomock.Any()).
					Return(nil).
					Times(1)

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
					Return("mock-unique-id", nil).
					Times(1)

				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				fileRepo.
					EXPECT().
					CreateFile(gomock.Eq(ctx), gomock.Any()).
					Return(createFileRes, nil).
					Times(1)

				res, err := s.UploadFile(ctx, opts...)

				Expect(res).To(Equal(r))
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
			fn = service.NewCreateFn(data, fileManager)
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

				Expect(err).To(Equal(file.ErrExists))
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

	Context("DeleteFile function", Label("unit"), func() {
		var (
			ctx         context.Context
			currentTs   time.Time
			p           service.DeleteFileParam
			fileRepo    *mock_repository.MockFile
			fileManager *mock_filesystem.MockFileManager
			clock       *mock_datetime.MockClock
			log         *mock_logging.MockLogger
			validator   *mock_validation.MockValidator
			s           service.File
			deleteRes   *repository.DeleteFileResult
			r           *service.DeleteFileResult
		)

		BeforeEach(func() {
			currentTs = time.Now().UTC()
			ctx = context.Background()
			p = service.DeleteFileParam{
				FileId: "mock-file-id",
			}
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fileRepo = mock_repository.NewMockFile(ctrl)
			fileManager = mock_filesystem.NewMockFileManager(ctrl)
			dirManager := mock_filesystem.NewMockDirectoryManager(ctrl)
			identifier := mock_identifier.NewMockIdentifier(ctrl)
			clock = mock_datetime.NewMockClock(ctrl)
			locator := mock_file.NewMockUploadLocation(ctrl)
			log = mock_logging.NewMockLogger(ctrl)
			validator = mock_validation.NewMockValidator(ctrl)
			s = service.NewFile(service.FileParam{
				FileRepo:    fileRepo,
				FileManager: fileManager,
				DirManager:  dirManager,
				Logger:      log,
				Identifier:  identifier,
				Clock:       clock,
				Locator:     locator,
				Validator:   validator,
				Config: &service.FileConfig{
					UploadDir: "temp",
				},
			})
			deleteRes = &repository.DeleteFileResult{
				DeletedAt: currentTs,
			}
			r = &service.DeleteFileResult{
				Success: system.Success{
					Code:    1000,
					Message: "success delete file",
				},
				DeletedAt: currentTs,
			}

			log.
				EXPECT().
				Debug("In function: DeleteFile").
				Times(1)
			log.
				EXPECT().
				Debug("Returning function: DeleteFile").
				Times(1)
		})

		When("parameter is not valid", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(fmt.Errorf("invalid data")).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Code).To(Equal(int32(1002)))
				Expect(err.Message).To(Equal("invalid data"))
			})
		})

		When("failed delete file", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				fileRepo.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, fmt.Errorf("network error")).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Code).To(Equal(int32(1001)))
				Expect(err.Message).To(Equal("network error"))
			})
		})

		When("file is not available", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				fileRepo.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, repository.ErrNotFound).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Code).To(Equal(int32(1004)))
				Expect(err.Message).To(Equal("file is not found"))
			})
		})

		When("file is deleted", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				fileRepo.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, repository.ErrDeleted).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Code).To(Equal(int32(1004)))
				Expect(err.Message).To(Equal("file is deleted"))
			})
		})

		When("success delete file", func() {
			It("should return result", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

				clock.
					EXPECT().
					Now().
					Return(currentTs).
					Times(1)

				fileRepo.
					EXPECT().
					DeleteFile(gomock.Eq(ctx), gomock.Any()).
					Return(deleteRes, nil).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(Equal(r))
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
			fn = service.NewDeleteFn(fileManager)
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

				Expect(err).To(Equal(file.ErrNotFound))
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
