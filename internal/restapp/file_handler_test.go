package restapp_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"time"

	api "github.com/go-seidon/hippo/api/restapp"
	"github.com/go-seidon/hippo/internal/file"
	mock_file "github.com/go-seidon/hippo/internal/file/mock"
	"github.com/go-seidon/hippo/internal/restapp"
	mock_restapp "github.com/go-seidon/hippo/internal/restapp/mock"
	mock_io "github.com/go-seidon/provider/io/mock"
	mock_logging "github.com/go-seidon/provider/logging/mock"
	"github.com/go-seidon/provider/serialization"
	"github.com/go-seidon/provider/serialization/json"
	mock_serialization "github.com/go-seidon/provider/serialization/mock"
	"github.com/go-seidon/provider/system"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("File Handler", func() {

	Context("DeleteFile Handler", Label("unit"), func() {
		var (
			handler     http.HandlerFunc
			r           *http.Request
			w           *mock_restapp.MockResponseWriter
			log         *mock_logging.MockLogger
			serializer  *mock_serialization.MockSerializer
			fileService *mock_file.MockFile
			p           file.DeleteFileParam
		)

		BeforeEach(func() {
			t := GinkgoT()
			r = mux.SetURLVars(&http.Request{}, map[string]string{
				"id": "mock-file-id",
			})
			ctrl := gomock.NewController(t)
			w = mock_restapp.NewMockResponseWriter(ctrl)
			log = mock_logging.NewMockLogger(ctrl)
			serializer = mock_serialization.NewMockSerializer(ctrl)
			fileService = mock_file.NewMockFile(ctrl)
			fileHandler := restapp.NewFileHandler(restapp.FileHandlerParam{
				Logger:      log,
				Serializer:  serializer,
				Config:      &restapp.RestAppConfig{},
				FileService: fileService,
			})
			handler = fileHandler.DeleteFileById
			p = file.DeleteFileParam{
				FileId: "mock-file-id",
			}
		})

		When("failed delete file", func() {
			It("should write response", func() {
				fileService.
					EXPECT().
					DeleteFile(gomock.Any(), gomock.Eq(p)).
					Return(nil, &system.Error{
						Code:    1001,
						Message: "network error",
					}).
					Times(1)

				b := restapp.ResponseBody{
					Code:    1001,
					Message: "network error",
				}
				serializer.
					EXPECT().
					Marshal(b).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(500)).
					Times(1)

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})

		When("file is not found", func() {
			It("should write response", func() {
				fileService.
					EXPECT().
					DeleteFile(gomock.Any(), gomock.Eq(p)).
					Return(nil, &system.Error{
						Code:    1004,
						Message: "file is not found",
					}).
					Times(1)

				b := restapp.ResponseBody{
					Code:    1004,
					Message: "file is not found",
				}
				serializer.
					EXPECT().
					Marshal(b).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(404)).
					Times(1)

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})

		When("there is invalid data", func() {
			It("should write response", func() {
				fileService.
					EXPECT().
					DeleteFile(gomock.Any(), gomock.Eq(p)).
					Return(nil, &system.Error{
						Code:    1002,
						Message: "invalid data",
					}).
					Times(1)

				b := restapp.ResponseBody{
					Code:    1002,
					Message: "invalid data",
				}
				serializer.
					EXPECT().
					Marshal(b).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(400)).
					Times(1)

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})

		When("success delete file", func() {
			It("should write response", func() {
				res := &file.DeleteFileResult{
					Success: system.Success{
						Code:    1000,
						Message: "success delete file",
					},
					DeletedAt: time.Now(),
				}
				b := restapp.ResponseBody{
					Code:    1000,
					Message: "success delete file",
					Data: &api.DeleteFileByIdData{
						DeletedAt: res.DeletedAt.UnixMilli(),
					},
				}

				fileService.
					EXPECT().
					DeleteFile(gomock.Any(), gomock.Eq(p)).
					Return(res, nil).
					Times(1)

				serializer.
					EXPECT().
					Marshal(b).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(200)).
					Times(1)

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})
	})

	Context("RetrieveFile Handler", Label("unit"), func() {
		var (
			handler     http.HandlerFunc
			r           *http.Request
			w           *mock_restapp.MockResponseWriter
			log         *mock_logging.MockLogger
			serializer  *mock_serialization.MockSerializer
			fileService *mock_file.MockFile
			fileData    *mock_io.MockReadCloser
			p           file.RetrieveFileParam
		)

		BeforeEach(func() {
			t := GinkgoT()
			r = mux.SetURLVars(&http.Request{}, map[string]string{
				"id": "mock-file-id",
			})
			ctrl := gomock.NewController(t)
			w = mock_restapp.NewMockResponseWriter(ctrl)
			log = mock_logging.NewMockLogger(ctrl)
			serializer = mock_serialization.NewMockSerializer(ctrl)
			fileService = mock_file.NewMockFile(ctrl)
			fileData = mock_io.NewMockReadCloser(ctrl)
			fileHandler := restapp.NewFileHandler(restapp.FileHandlerParam{
				Logger:      log,
				Serializer:  serializer,
				Config:      &restapp.RestAppConfig{},
				FileService: fileService,
			})
			handler = fileHandler.RetrieveFileById
			p = file.RetrieveFileParam{
				FileId: "mock-file-id",
			}
		})

		When("failed retrieve file", func() {
			It("should write response", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Any(), gomock.Eq(p)).
					Return(nil, &system.Error{
						Code:    1001,
						Message: "network error",
					}).
					Times(1)

				b := restapp.ResponseBody{
					Code:    1001,
					Message: "network error",
				}
				serializer.
					EXPECT().
					Marshal(b).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(500)).
					Times(1)

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})

		When("file is not found", func() {
			It("should write response", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Any(), gomock.Eq(p)).
					Return(nil, &system.Error{
						Code:    1004,
						Message: "file is not found",
					}).
					Times(1)

				b := restapp.ResponseBody{
					Code:    1004,
					Message: "file is not found",
				}
				serializer.
					EXPECT().
					Marshal(b).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(404)).
					Times(1)

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})

		When("there is invalid data", func() {
			It("should write response", func() {
				fileService.
					EXPECT().
					RetrieveFile(gomock.Any(), gomock.Eq(p)).
					Return(nil, &system.Error{
						Code:    1002,
						Message: "invalid data",
					}).
					Times(1)

				b := restapp.ResponseBody{
					Code:    1002,
					Message: "invalid data",
				}
				serializer.
					EXPECT().
					Marshal(b).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(400)).
					Times(1)

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})

		When("failed read file", func() {
			It("should write response", func() {
				fileData.
					EXPECT().
					Close().
					Times(1)

				fileData.
					EXPECT().
					Read(gomock.Any()).
					Return(0, fmt.Errorf("read error")).
					Times(1)

				res := &file.RetrieveFileResult{
					Success: system.Success{
						Code:    1000,
						Message: "success retrieve file",
					},
					Data: fileData,
				}
				fileService.
					EXPECT().
					RetrieveFile(gomock.Any(), gomock.Eq(p)).
					Return(res, nil).
					Times(1)

				b := restapp.ResponseBody{
					Code:    1001,
					Message: "read error",
				}
				serializer.
					EXPECT().
					Marshal(b).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(500)).
					Times(1)

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})

		When("mimetype is empty", func() {
			It("should write response", func() {
				fileData.
					EXPECT().
					Close().
					Times(1)

				fileData.
					EXPECT().
					Read(gomock.Any()).
					Return(0, io.EOF).
					Times(1)

				res := &file.RetrieveFileResult{
					Success: system.Success{
						Code:    1000,
						Message: "success retrieve file",
					},
					Data:      fileData,
					UniqueId:  "mock-unique-id",
					Name:      "mock-name",
					Path:      "mock-path",
					MimeType:  "",
					Extension: "mock-extension",
					DeletedAt: nil,
				}
				fileService.
					EXPECT().
					RetrieveFile(gomock.Any(), gomock.Eq(p)).
					Return(res, nil).
					Times(1)

				w.
					EXPECT().
					Header().
					Return(http.Header{}).
					Times(5)

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})

		When("mimetype is not empty", func() {
			It("should write response", func() {
				fileData.
					EXPECT().
					Close().
					Times(1)

				fileData.
					EXPECT().
					Read(gomock.Any()).
					Return(0, io.EOF).
					Times(1)

				res := &file.RetrieveFileResult{
					Success: system.Success{
						Code:    1000,
						Message: "success retrieve file",
					},
					Data:      fileData,
					UniqueId:  "mock-unique-id",
					Name:      "mock-name",
					Path:      "mock-path",
					MimeType:  "text/plain",
					Extension: "mock-extension",
					DeletedAt: nil,
				}
				fileService.
					EXPECT().
					RetrieveFile(gomock.Any(), gomock.Eq(p)).
					Return(res, nil).
					Times(1)

				w.
					EXPECT().
					Header().
					Return(http.Header{}).
					Times(5)

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})
	})

	Context("UploadFile Handler", Label("integration"), Ordered, func() {
		var (
			currentTimestamp time.Time
			ctx              context.Context
			ctrl             *gomock.Controller
			r                *http.Request
			body             *bytes.Buffer
			writer           *multipart.Writer
			handler          http.HandlerFunc
			log              *mock_logging.MockLogger
			serializer       serialization.Serializer
			uploadService    *mock_file.MockFile
		)

		BeforeEach(func() {
			currentTimestamp = time.Now()
			t := GinkgoT()
			ctx = context.Background()
			ctrl = gomock.NewController(t)

			body = new(bytes.Buffer)
			writer = multipart.NewWriter(body)
			_, err := writer.CreateFormFile("file", "app.go")
			if err != nil {
				AbortSuite("failed create file mock: " + err.Error())
			}
			writer.Close()

			r, _ = http.NewRequest(http.MethodPost, "/v1/file", body)
			r.Header.Add("Content-Type", writer.FormDataContentType())

			log = mock_logging.NewMockLogger(ctrl)
			serializer = json.NewSerializer()
			uploadService = mock_file.NewMockFile(ctrl)
			fileHandler := restapp.NewFileHandler(restapp.FileHandlerParam{
				Logger:      log,
				Serializer:  serializer,
				Config:      &restapp.RestAppConfig{},
				FileService: uploadService,
			})
			handler = fileHandler.UploadFile
		})

		When("failed parse form file", func() {
			It("should return error", func() {

				r, _ := http.NewRequest(http.MethodPost, "/v1/file", nil)
				w := httptest.NewRecorder()

				handler.ServeHTTP(w, r)

				resBody := restapp.ResponseBody{}
				serializer.Unmarshal(w.Body.Bytes(), &resBody)

				Expect(w.Code).To(Equal(400))
				Expect(resBody.Code).To(Equal(int32(1001)))
				Expect(resBody.Message).To(Equal("request Content-Type isn't multipart/form-data"))
				Expect(resBody.Data).To(BeNil())
			})
		})

		When("failed upload file", func() {
			It("should return error", func() {
				uploadService.
					EXPECT().
					UploadFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, &system.Error{
						Code:    1001,
						Message: "disk error",
					}).
					Times(1)

				w := httptest.NewRecorder()

				handler.ServeHTTP(w, r)

				resBody := restapp.ResponseBody{}
				serializer.Unmarshal(w.Body.Bytes(), &resBody)

				Expect(w.Code).To(Equal(500))
				Expect(resBody.Code).To(Equal(int32(1001)))
				Expect(resBody.Message).To(Equal("disk error"))
				Expect(resBody.Data).To(BeNil())
			})
		})

		When("there are invalid data", func() {
			It("should return error", func() {
				uploadService.
					EXPECT().
					UploadFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, &system.Error{
						Code:    1002,
						Message: "invalid data",
					}).
					Times(1)

				w := httptest.NewRecorder()

				handler.ServeHTTP(w, r)

				resBody := restapp.ResponseBody{}
				serializer.Unmarshal(w.Body.Bytes(), &resBody)

				Expect(w.Code).To(Equal(400))
				Expect(resBody.Code).To(Equal(int32(1002)))
				Expect(resBody.Message).To(Equal("invalid data"))
				Expect(resBody.Data).To(BeNil())
			})
		})

		When("success upload file", func() {
			It("should return result", func() {
				uploadRes := &file.UploadFileResult{
					UniqueId:   "mock-unique-id",
					Name:       "dolpin.jpg",
					Path:       "mock/location/mock-unique-id.jpg",
					Mimetype:   "image/jpeg",
					Extension:  "jpg",
					Size:       200,
					UploadedAt: currentTimestamp,
				}
				uploadService.
					EXPECT().
					UploadFile(gomock.Eq(ctx), gomock.Any()).
					Return(uploadRes, nil).
					Times(1)

				w := httptest.NewRecorder()

				handler.ServeHTTP(w, r)

				resBody := restapp.ResponseBody{}
				serializer.Unmarshal(w.Body.Bytes(), &resBody)
				data := map[string]interface{}{
					"id":          uploadRes.UniqueId,
					"name":        uploadRes.Name,
					"extension":   uploadRes.Extension,
					"mimetype":    uploadRes.Mimetype,
					"size":        float64(200),
					"uploaded_at": float64(uploadRes.UploadedAt.UnixMilli()),
				}

				Expect(w.Code).To(Equal(200))
				Expect(resBody.Code).To(Equal(int32(1000)))
				Expect(resBody.Message).To(Equal("success upload file"))
				Expect(resBody.Data).To(Equal(data))
			})
		})
	})
})