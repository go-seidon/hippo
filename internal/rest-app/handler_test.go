package rest_app_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-seidon/local/internal/deleting"
	mock_deleting "github.com/go-seidon/local/internal/deleting/mock"
	"github.com/go-seidon/local/internal/healthcheck"
	mock_healthcheck "github.com/go-seidon/local/internal/healthcheck/mock"
	"github.com/go-seidon/local/internal/mock"
	rest_app "github.com/go-seidon/local/internal/rest-app"
	"github.com/go-seidon/local/internal/retrieving"
	"github.com/go-seidon/local/internal/serialization"
	"github.com/go-seidon/local/internal/uploading"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler Package", func() {

	Context("NotFoundHandler", Label("unit"), func() {
		var (
			handler    http.HandlerFunc
			r          *http.Request
			w          *mock.MockResponseWriter
			log        *mock.MockLogger
			serializer *mock.MockSerializer
		)

		BeforeEach(func() {
			t := GinkgoT()
			r = &http.Request{}
			ctrl := gomock.NewController(t)
			w = mock.NewMockResponseWriter(ctrl)
			log = mock.NewMockLogger(ctrl)
			serializer = mock.NewMockSerializer(ctrl)
			handler = rest_app.NewNotFoundHandler(log, serializer)
		})

		When("success call the function", func() {
			It("should write response", func() {

				b := rest_app.ResponseBody{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				}

				serializer.
					EXPECT().
					Marshal(b).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					Header().
					Return(http.Header{}).
					Times(1)

				w.
					EXPECT().
					WriteHeader(http.StatusNotFound).
					Times(1)

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})
	})

	Context("MethodNowAllowedHandler", Label("unit"), func() {
		var (
			handler    http.HandlerFunc
			r          *http.Request
			w          *mock.MockResponseWriter
			log        *mock.MockLogger
			serializer *mock.MockSerializer
		)

		BeforeEach(func() {
			t := GinkgoT()
			r = &http.Request{}
			ctrl := gomock.NewController(t)
			w = mock.NewMockResponseWriter(ctrl)
			log = mock.NewMockLogger(ctrl)
			serializer = mock.NewMockSerializer(ctrl)
			handler = rest_app.NewMethodNotAllowedHandler(log, serializer)
		})

		When("success call the function", func() {
			It("should write response", func() {

				b := rest_app.ResponseBody{
					Code:    "ERROR",
					Message: "method is not allowed",
				}

				serializer.
					EXPECT().
					Marshal(b).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					Header().
					Return(http.Header{}).
					Times(1)

				w.
					EXPECT().
					WriteHeader(http.StatusMethodNotAllowed).
					Times(1)

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})
	})

	Context("RootHandler", Label("unit"), func() {
		var (
			handler    http.HandlerFunc
			r          *http.Request
			w          *mock.MockResponseWriter
			log        *mock.MockLogger
			serializer *mock.MockSerializer
		)

		BeforeEach(func() {
			t := GinkgoT()
			r = &http.Request{}
			ctrl := gomock.NewController(t)
			w = mock.NewMockResponseWriter(ctrl)
			log = mock.NewMockLogger(ctrl)
			serializer = mock.NewMockSerializer(ctrl)
			cfg := &rest_app.RestAppConfig{
				AppName:    "mock-name",
				AppVersion: "mock-version",
			}
			handler = rest_app.NewRootHandler(log, serializer, cfg)
		})

		When("success call the function", func() {
			It("should write response", func() {

				b := rest_app.ResponseBody{
					Code:    "SUCCESS",
					Message: "success",
					Data: struct {
						AppName    string `json:"app_name"`
						AppVersion string `json:"app_version"`
					}{
						AppName:    "mock-name",
						AppVersion: "mock-version",
					},
				}

				serializer.
					EXPECT().
					Marshal(b).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(200))

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})
	})

	Context("HealthCheckHandler", Label("unit"), func() {
		var (
			handler       http.HandlerFunc
			r             *http.Request
			w             *mock.MockResponseWriter
			log           *mock.MockLogger
			serializer    *mock.MockSerializer
			healthService *mock_healthcheck.MockHealthCheck
		)

		BeforeEach(func() {
			t := GinkgoT()
			r = &http.Request{}
			ctrl := gomock.NewController(t)
			w = mock.NewMockResponseWriter(ctrl)
			log = mock.NewMockLogger(ctrl)
			serializer = mock.NewMockSerializer(ctrl)
			healthService = mock_healthcheck.NewMockHealthCheck(ctrl)
			handler = rest_app.NewHealthCheckHandler(log, serializer, healthService)
		})

		When("failed check service health", func() {
			It("should write response", func() {

				err := fmt.Errorf("failed check health")

				b := rest_app.ResponseBody{
					Code:    "ERROR",
					Message: err.Error(),
				}

				healthService.
					EXPECT().
					Check().
					Return(nil, err).
					Times(1)

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

		When("success check service health", func() {
			It("should write response", func() {

				currentTimestamp := time.Now()
				res := &healthcheck.CheckResult{
					Status: "WARNING",
					Items: map[string]healthcheck.CheckResultItem{
						"app-disk": {
							Name:      "app-disk",
							Status:    "FAILED",
							Error:     "Critical: disk usage too high 96.71 percent",
							CheckedAt: currentTimestamp,
							Metadata:  nil,
						},
						"internet-connection": {
							Name:      "internet-connection",
							Status:    "OK",
							Error:     "",
							CheckedAt: currentTimestamp,
							Metadata:  nil,
						},
					},
				}
				jobs := map[string]struct {
					Name      string      `json:"name"`
					Status    string      `json:"status"`
					CheckedAt time.Time   `json:"checked_at"`
					Error     string      `json:"error"`
					Metadata  interface{} `json:"metadata"`
				}{
					"app-disk": {
						Name:      "app-disk",
						Status:    "FAILED",
						Error:     "Critical: disk usage too high 96.71 percent",
						CheckedAt: currentTimestamp,
						Metadata:  nil,
					},
					"internet-connection": {
						Name:      "internet-connection",
						Status:    "OK",
						Error:     "",
						CheckedAt: currentTimestamp,
						Metadata:  nil,
					},
				}

				b := rest_app.ResponseBody{
					Data: struct {
						Status  string `json:"status"`
						Details map[string]struct {
							Name      string      `json:"name"`
							Status    string      `json:"status"`
							CheckedAt time.Time   `json:"checked_at"`
							Error     string      `json:"error"`
							Metadata  interface{} `json:"metadata"`
						} `json:"details"`
					}{
						Status:  "WARNING",
						Details: jobs,
					},
					Code:    "SUCCESS",
					Message: "success check service health",
				}

				healthService.
					EXPECT().
					Check().
					Return(res, nil).
					Times(1)

				serializer.
					EXPECT().
					Marshal(b).
					Return([]byte{}, nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(200))

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})
	})

	Context("NewDeleteFileHandler", Label("unit"), func() {
		var (
			handler       http.HandlerFunc
			r             *http.Request
			w             *mock.MockResponseWriter
			log           *mock.MockLogger
			serializer    *mock.MockSerializer
			deleteService *mock_deleting.MockDeleter
			p             deleting.DeleteFileParam
		)

		BeforeEach(func() {
			t := GinkgoT()
			r = mux.SetURLVars(&http.Request{}, map[string]string{
				"id": "mock-file-id",
			})
			ctrl := gomock.NewController(t)
			w = mock.NewMockResponseWriter(ctrl)
			log = mock.NewMockLogger(ctrl)
			serializer = mock.NewMockSerializer(ctrl)
			deleteService = mock_deleting.NewMockDeleter(ctrl)
			handler = rest_app.NewDeleteFileHandler(log, serializer, deleteService)
			p = deleting.DeleteFileParam{
				FileId: "mock-file-id",
			}
		})

		When("failed delete file", func() {
			It("should write response", func() {

				err := fmt.Errorf("failed delete file")

				b := rest_app.ResponseBody{
					Code:    "ERROR",
					Message: err.Error(),
				}

				deleteService.
					EXPECT().
					DeleteFile(gomock.Any(), gomock.Eq(p)).
					Return(nil, err).
					Times(1)

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

		When("file is not found", func() {
			It("should write response", func() {

				err := deleting.ErrorResourceNotFound

				b := rest_app.ResponseBody{
					Code:    "NOT_FOUND",
					Message: err.Error(),
				}

				deleteService.
					EXPECT().
					DeleteFile(gomock.Any(), gomock.Eq(p)).
					Return(nil, err).
					Times(1)

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

		When("success delete file", func() {
			It("should write response", func() {
				res := &deleting.DeleteFileResult{
					DeletedAt: time.Now(),
				}
				b := rest_app.ResponseBody{
					Code:    "SUCCESS",
					Message: "success delete file",
					Data: struct {
						DeletedAt int64 `json:"deleted_at"`
					}{
						DeletedAt: res.DeletedAt.UnixMilli(),
					},
				}

				deleteService.
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

	Context("NewRetrieveFileHandler", Label("unit"), func() {
		var (
			ctx             context.Context
			handler         http.HandlerFunc
			r               *http.Request
			w               *mock.MockResponseWriter
			log             *mock.MockLogger
			serializer      *mock.MockSerializer
			retrieveService *mock.MockRetriever
			fileData        *mock.MockReadCloser
			p               retrieving.RetrieveFileParam
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctx = context.Background()
			r = mux.SetURLVars(&http.Request{}, map[string]string{
				"id": "mock-file-id",
			})
			ctrl := gomock.NewController(t)
			w = mock.NewMockResponseWriter(ctrl)
			log = mock.NewMockLogger(ctrl)
			serializer = mock.NewMockSerializer(ctrl)
			retrieveService = mock.NewMockRetriever(ctrl)
			fileData = mock.NewMockReadCloser(ctrl)
			handler = rest_app.NewRetrieveFileHandler(log, serializer, retrieveService)
			p = retrieving.RetrieveFileParam{
				FileId: "mock-file-id",
			}
		})

		When("failed retrieve file", func() {
			It("should write response", func() {

				err := fmt.Errorf("failed retrieve file")

				b := rest_app.ResponseBody{
					Code:    "ERROR",
					Message: err.Error(),
				}

				retrieveService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(p)).
					Return(nil, err).
					Times(1)

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

		When("file is not found", func() {
			It("should write response", func() {

				err := retrieving.ErrorResourceNotFound

				b := rest_app.ResponseBody{
					Code:    "NOT_FOUND",
					Message: err.Error(),
				}

				retrieveService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(p)).
					Return(nil, err).
					Times(1)

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

				res := &retrieving.RetrieveFileResult{
					Data: fileData,
				}

				b := rest_app.ResponseBody{
					Code:    "ERROR",
					Message: "read error",
				}

				retrieveService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(p)).
					Return(res, nil).
					Times(1)

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

				res := &retrieving.RetrieveFileResult{
					Data:      fileData,
					UniqueId:  "mock-unique-id",
					Name:      "mock-name",
					Path:      "mock-path",
					MimeType:  "",
					Extension: "mock-extension",
					DeletedAt: nil,
				}

				retrieveService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(p)).
					Return(res, nil).
					Times(1)

				w.EXPECT().
					Header().
					Return(http.Header{}).
					Times(1)

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

				res := &retrieving.RetrieveFileResult{
					Data:      fileData,
					UniqueId:  "mock-unique-id",
					Name:      "mock-name",
					Path:      "mock-path",
					MimeType:  "text/plain",
					Extension: "mock-extension",
					DeletedAt: nil,
				}

				retrieveService.
					EXPECT().
					RetrieveFile(gomock.Eq(ctx), gomock.Eq(p)).
					Return(res, nil).
					Times(1)

				w.EXPECT().
					Header().
					Return(http.Header{}).
					Times(1)

				w.
					EXPECT().
					Write([]byte{}).
					Times(1)

				handler.ServeHTTP(w, r)
			})
		})
	})

	Context("NewUploadFileHandler", Label("integration"), Ordered, func() {
		var (
			currentTimestamp time.Time
			ctx              context.Context
			ctrl             *gomock.Controller
			r                *http.Request
			body             *bytes.Buffer
			writer           *multipart.Writer
			handler          http.HandlerFunc
			log              *mock.MockLogger
			serializer       serialization.Serializer
			uploadService    *mock.MockUploader
			locator          *mock.MockUploadLocation
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

			log = mock.NewMockLogger(ctrl)
			serializer = serialization.NewJsonSerializer()
			uploadService = mock.NewMockUploader(ctrl)
			locator = mock.NewMockUploadLocation(ctrl)
			cfg := &rest_app.RestAppConfig{}
			handler = rest_app.NewUploadFileHandler(
				log, serializer, uploadService,
				locator, cfg,
			)
		})

		When("failed parse form file", func() {
			It("should return error", func() {

				r, _ := http.NewRequest(http.MethodPost, "/v1/file", nil)
				w := httptest.NewRecorder()

				handler.ServeHTTP(w, r)

				resBody := rest_app.ResponseBody{}
				serializer.Unmarshal(w.Body.Bytes(), &resBody)

				Expect(w.Code).To(Equal(400))
				Expect(resBody.Code).To(Equal("ERROR"))
				Expect(resBody.Message).To(Equal("request Content-Type isn't multipart/form-data"))
				Expect(resBody.Data).To(BeNil())
			})
		})

		When("failed upload file", func() {
			It("should return error", func() {

				locator.
					EXPECT().
					GetLocation().
					Return("mock/location").
					Times(1)

				uploadService.
					EXPECT().
					UploadFile(gomock.Eq(ctx), gomock.Any()).
					Return(nil, fmt.Errorf("disk error")).
					Times(1)

				w := httptest.NewRecorder()

				handler.ServeHTTP(w, r)

				resBody := rest_app.ResponseBody{}
				serializer.Unmarshal(w.Body.Bytes(), &resBody)

				Expect(w.Code).To(Equal(400))
				Expect(resBody.Code).To(Equal("ERROR"))
				Expect(resBody.Message).To(Equal("disk error"))
				Expect(resBody.Data).To(BeNil())
			})
		})

		When("success upload file", func() {
			It("should return result", func() {

				locator.
					EXPECT().
					GetLocation().
					Return("mock/location").
					Times(1)

				uploadRes := &uploading.UploadFileResult{
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

				resBody := rest_app.ResponseBody{}
				serializer.Unmarshal(w.Body.Bytes(), &resBody)
				data := map[string]interface{}{
					"id":          uploadRes.UniqueId,
					"name":        uploadRes.Name,
					"mimetype":    uploadRes.Mimetype,
					"extension":   uploadRes.Extension,
					"size":        float64(200),
					"uploaded_at": float64(uploadRes.UploadedAt.UnixMilli()),
				}

				Expect(w.Code).To(Equal(200))
				Expect(resBody.Code).To(Equal("SUCCESS"))
				Expect(resBody.Message).To(Equal("success upload file"))
				Expect(resBody.Data).To(Equal(data))
			})
		})
	})
})
