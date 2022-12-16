package rest_app_test

import (
	"net/http"

	"github.com/go-seidon/hippo/generated/restapp"
	rest_app "github.com/go-seidon/hippo/internal/rest-app"
	mock_restapp "github.com/go-seidon/hippo/internal/rest-app/mock"
	mock_logging "github.com/go-seidon/provider/logging/mock"
	mock_serialization "github.com/go-seidon/provider/serialization/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Basic Handler", func() {

	Context("NotFound Handler", Label("unit"), func() {
		var (
			handler    http.HandlerFunc
			r          *http.Request
			w          *mock_restapp.MockResponseWriter
			log        *mock_logging.MockLogger
			serializer *mock_serialization.MockSerializer
		)

		BeforeEach(func() {
			t := GinkgoT()
			r = &http.Request{}
			ctrl := gomock.NewController(t)
			w = mock_restapp.NewMockResponseWriter(ctrl)
			log = mock_logging.NewMockLogger(ctrl)
			serializer = mock_serialization.NewMockSerializer(ctrl)
			basicHandler := rest_app.NewBasicHandler(rest_app.BasicHandlerParam{
				Logger:     log,
				Serializer: serializer,
				Config:     &rest_app.RestAppConfig{},
			})
			handler = http.HandlerFunc(basicHandler.NotFound)
		})

		When("success call the function", func() {
			It("should write response", func() {

				b := rest_app.ResponseBody{
					Code:    1004,
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

	Context("MethodNowAllowed Handler", Label("unit"), func() {
		var (
			handler    http.HandlerFunc
			r          *http.Request
			w          *mock_restapp.MockResponseWriter
			log        *mock_logging.MockLogger
			serializer *mock_serialization.MockSerializer
		)

		BeforeEach(func() {
			t := GinkgoT()
			r = &http.Request{}
			ctrl := gomock.NewController(t)
			w = mock_restapp.NewMockResponseWriter(ctrl)
			log = mock_logging.NewMockLogger(ctrl)
			serializer = mock_serialization.NewMockSerializer(ctrl)
			basicHandler := rest_app.NewBasicHandler(rest_app.BasicHandlerParam{
				Logger:     log,
				Serializer: serializer,
				Config:     &rest_app.RestAppConfig{},
			})
			handler = http.HandlerFunc(basicHandler.MethodNotAllowed)
		})

		When("success call the function", func() {
			It("should write response", func() {

				b := rest_app.ResponseBody{
					Code:    1001,
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

	Context("GetAppInfi Handler", Label("unit"), func() {
		var (
			handler    http.HandlerFunc
			r          *http.Request
			w          *mock_restapp.MockResponseWriter
			log        *mock_logging.MockLogger
			serializer *mock_serialization.MockSerializer
		)

		BeforeEach(func() {
			t := GinkgoT()
			r = &http.Request{}
			ctrl := gomock.NewController(t)
			w = mock_restapp.NewMockResponseWriter(ctrl)
			log = mock_logging.NewMockLogger(ctrl)
			serializer = mock_serialization.NewMockSerializer(ctrl)
			cfg := &rest_app.RestAppConfig{
				AppName:    "mock-name",
				AppVersion: "mock-version",
			}
			basicHandler := rest_app.NewBasicHandler(rest_app.BasicHandlerParam{
				Logger:     log,
				Serializer: serializer,
				Config:     cfg,
			})
			handler = basicHandler.GetAppInfo
		})

		When("success call the function", func() {
			It("should write response", func() {

				b := rest_app.ResponseBody{
					Code:    1000,
					Message: "success",
					Data: &restapp.GetAppInfoData{
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
})
