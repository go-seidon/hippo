package file_test

import (
	"context"
	"fmt"
	"os"

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

var _ = Describe("Retriever", func() {

	Context("RetrieveFile function", Label("unit"), func() {
		var (
			ctx           context.Context
			p             file.RetrieveFileParam
			r             *file.RetrieveFileResult
			fileRepo      *mock_repository.MockFileRepository
			fileManager   *mock_filesystem.MockFileManager
			log           *mock_logging.MockLogger
			validator     *mock_validation.MockValidator
			s             file.File
			retrieveParam repository.RetrieveFileParam
			retrieveRes   *repository.RetrieveFileResult
			openParam     filesystem.OpenFileParam
			openRes       *filesystem.OpenFileResult
		)

		BeforeEach(func() {
			ctx = context.Background()
			p = file.RetrieveFileParam{
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
				Identifier:  identifier,
				Logger:      log,
				Locator:     locator,
				Validator:   validator,
				Config: &file.FileConfig{
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
				MimeType:  "mock-mimetype",
				Extension: "mock-extension",
			}
			openParam = filesystem.OpenFileParam{
				Path: retrieveRes.Path,
			}
			osFile := &os.File{}
			openRes = &filesystem.OpenFileResult{
				File: osFile,
			}
			r = &file.RetrieveFileResult{
				Data:      osFile,
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

		When("parameter are not valid", func() {
			It("should return error", func() {
				p.FileId = ""

				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(fmt.Errorf("invalid data")).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid data")))
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
				Expect(err).To(Equal(file.ErrorNotFound))
			})
		})

		When("file record is deleted", func() {
			It("should return error", func() {
				validator.
					EXPECT().
					Validate(gomock.Eq(p)).
					Return(nil).
					Times(1)

				fileRepo.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(retrieveParam)).
					Return(nil, repository.ErrDeleted).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(file.ErrorNotFound))
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
				Expect(err).To(Equal(fmt.Errorf("db error")))
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
				Expect(err).To(Equal(file.ErrorNotFound))
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
				Expect(err).To(Equal(fmt.Errorf("disk error")))
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
})
