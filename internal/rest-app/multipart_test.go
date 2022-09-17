package rest_app_test

import (
	"fmt"
	"io"
	"mime/multipart"

	mock_io "github.com/go-seidon/local/internal/io/mock"
	rest_app "github.com/go-seidon/local/internal/rest-app"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Multipart Package", func() {
	Context("ParseFileName function", Label("unit"), func() {
		var (
			fh *multipart.FileHeader
		)

		BeforeEach(func() {
			fh = &multipart.FileHeader{
				Filename: "dolpin.jpeg",
			}
		})

		When("file name is empty", func() {
			It("should return empty", func() {
				fh.Filename = ""
				res := rest_app.ParseFileName(fh)

				Expect(res).To(Equal(""))
			})
		})

		When("name contain extension", func() {
			It("should return only name", func() {
				res := rest_app.ParseFileName(fh)

				Expect(res).To(Equal("dolpin"))
			})
		})

		When("name not contain extension", func() {
			It("should return only name", func() {
				fh.Filename = "dolpin"
				res := rest_app.ParseFileName(fh)

				Expect(res).To(Equal("dolpin"))
			})
		})

		When("name contain multiple dot", func() {
			It("should return only name", func() {
				fh.Filename = "dolpin.new.jpeg"
				res := rest_app.ParseFileName(fh)

				Expect(res).To(Equal("dolpin"))
			})
		})
	})

	Context("ParseFileExtension function", Label("unit"), func() {
		var (
			fh *multipart.FileHeader
		)

		BeforeEach(func() {
			fh = &multipart.FileHeader{
				Filename: "dolpin.jpeg",
			}
		})

		When("file name is empty", func() {
			It("should return empty", func() {
				fh.Filename = ""
				res := rest_app.ParseFileExtension(fh)

				Expect(res).To(Equal(""))
			})
		})

		When("name contain extension", func() {
			It("should return extension", func() {
				res := rest_app.ParseFileExtension(fh)

				Expect(res).To(Equal("jpeg"))
			})
		})

		When("name not contain extension", func() {
			It("should return empty", func() {
				fh.Filename = "dolpin"
				res := rest_app.ParseFileExtension(fh)

				Expect(res).To(Equal(""))
			})
		})

		When("name contain multiple dot", func() {
			It("should return only extension", func() {
				fh.Filename = "dolpin.new.jpeg"
				res := rest_app.ParseFileExtension(fh)

				Expect(res).To(Equal("jpeg"))
			})
		})
	})

	Context("ParseMultipartFile function", Label("unit"), func() {
		var (
			f  *mock_io.MockReadSeeker
			fh *multipart.FileHeader
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			f = mock_io.NewMockReadSeeker(ctrl)
			fh = &multipart.FileHeader{
				Filename: "dolpin.jpeg",
				Size:     200,
			}
		})

		When("failed read file", func() {
			It("should return error", func() {
				buff := make([]byte, 512)
				f.
					EXPECT().
					Read(gomock.Eq(buff)).
					Return(0, fmt.Errorf("disk error")).
					Times(1)

				res, err := rest_app.ParseMultipartFile(f, fh)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("disk error")))
			})
		})

		When("failed seek to the end of file", func() {
			It("should return error", func() {
				buff := make([]byte, 512)
				f.
					EXPECT().
					Read(gomock.Eq(buff)).
					Return(512, nil).
					Times(1)

				f.
					EXPECT().
					Seek(gomock.Eq(int64(0)), gomock.Eq(0)).
					Return(int64(0), fmt.Errorf("disk error")).
					Times(1)

				res, err := rest_app.ParseMultipartFile(f, fh)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("disk error")))
			})
		})

		When("reach end of file", func() {
			It("should return result", func() {
				buff := make([]byte, 512)
				f.
					EXPECT().
					Read(gomock.Eq(buff)).
					Return(200, io.EOF).
					Times(1)

				f.
					EXPECT().
					Seek(gomock.Eq(int64(0)), gomock.Eq(0)).
					Return(int64(1), nil).
					Times(1)

				res, err := rest_app.ParseMultipartFile(f, fh)

				expectedRes := &rest_app.FileInfo{
					Name:      "dolpin",
					Extension: "jpeg",
					Size:      200,
					Mimetype:  "application/octet-stream",
				}
				Expect(res).To(Equal(expectedRes))
				Expect(err).To(BeNil())
			})
		})

		When("success parse file", func() {
			It("should return result", func() {
				buff := make([]byte, 512)
				f.
					EXPECT().
					Read(gomock.Eq(buff)).
					Return(512, nil).
					Times(1)

				f.
					EXPECT().
					Seek(gomock.Eq(int64(0)), gomock.Eq(0)).
					Return(int64(1), nil).
					Times(1)

				res, err := rest_app.ParseMultipartFile(f, fh)

				expectedRes := &rest_app.FileInfo{
					Name:      "dolpin",
					Extension: "jpeg",
					Size:      200,
					Mimetype:  "application/octet-stream",
				}
				Expect(res).To(Equal(expectedRes))
				Expect(err).To(BeNil())
			})
		})

	})
})
