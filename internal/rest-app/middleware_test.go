package rest_app_test

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-seidon/hippo/internal/auth"
	mock_auth "github.com/go-seidon/hippo/internal/auth/mock"
	mock_logging "github.com/go-seidon/hippo/internal/logging/mock"
	rest_app "github.com/go-seidon/hippo/internal/rest-app"
	mock_restapp "github.com/go-seidon/hippo/internal/rest-app/mock"
	mock_serialization "github.com/go-seidon/hippo/internal/serialization/mock"
	mock_datetime "github.com/go-seidon/provider/datetime/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Middleware Package", func() {

	Context("DefaultHeaderMiddleware", Label("unit"), func() {
		var (
			r           *http.Request
			w           *mock_restapp.MockResponseWriter
			middleware  http.Handler
			httpHandler *mock_restapp.MockHandler
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)

			r = &http.Request{}
			w = mock_restapp.NewMockResponseWriter(ctrl)
			httpHandler = mock_restapp.NewMockHandler(ctrl)
			fn := rest_app.NewDefaultMiddleware(rest_app.DefaultMiddlewareParam{
				CorrelationIdHeaderKey: "X-Correlation-Id",
				CorrelationIdCtxKey:    rest_app.CorrelationIdCtxKey,
			})
			middleware = fn(httpHandler)
		})

		When("middleware is called", func() {
			It("should call serve http", func() {
				httpHandler.
					EXPECT().
					ServeHTTP(gomock.Eq(w), gomock.Any()).
					Times(1)

				w.EXPECT().
					Header().
					Return(http.Header{}).
					Times(1)

				middleware.ServeHTTP(w, r)
			})
		})
	})

	Context("NewBasicAuthMiddleware", Label("unit"), func() {
		var (
			a       *mock_auth.MockBasicAuth
			s       *mock_serialization.MockSerializer
			handler *mock_restapp.MockHandler
			m       http.Handler

			rw  *mock_restapp.MockResponseWriter
			req *http.Request
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			a = mock_auth.NewMockBasicAuth(ctrl)
			s = mock_serialization.NewMockSerializer(ctrl)
			handler = mock_restapp.NewMockHandler(ctrl)
			fn := rest_app.NewBasicAuthMiddleware(a, s)
			m = fn(handler)

			rw = mock_restapp.NewMockResponseWriter(ctrl)
			req = &http.Request{
				Header: http.Header{},
			}
			req.Header.Set("Authorization", "Basic basic-token")
		})

		When("basic auth is not specified", func() {
			It("should return error", func() {
				req.Header.Del("Authorization")

				b := rest_app.ResponseBody{
					Code:    1003,
					Message: "credential is not specified",
				}
				s.
					EXPECT().
					Marshal(gomock.Eq(b)).
					Return([]byte{}, nil).
					Times(1)
				rw.
					EXPECT().
					WriteHeader(401).
					Times(1)
				rw.
					EXPECT().
					Write(gomock.Eq([]byte{})).
					Times(1)

				m.ServeHTTP(rw, req)
			})
		})

		When("failed check credential", func() {
			It("should return error", func() {
				checkParam := auth.CheckCredentialParam{
					AuthToken: "basic-token",
				}
				a.
					EXPECT().
					CheckCredential(gomock.Any(), gomock.Eq(checkParam)).
					Return(nil, fmt.Errorf("db error")).
					Times(1)

				b := rest_app.ResponseBody{
					Code:    1003,
					Message: "failed check credential",
				}
				s.
					EXPECT().
					Marshal(gomock.Eq(b)).
					Return([]byte{}, nil).
					Times(1)
				rw.
					EXPECT().
					WriteHeader(401).
					Times(1)
				rw.
					EXPECT().
					Write(gomock.Eq([]byte{})).
					Times(1)

				m.ServeHTTP(rw, req)
			})
		})

		When("failed token is invalid", func() {
			It("should return error", func() {
				checkParam := auth.CheckCredentialParam{
					AuthToken: "basic-token",
				}
				checkRes := &auth.CheckCredentialResult{
					TokenValid: false,
				}
				a.
					EXPECT().
					CheckCredential(gomock.Any(), gomock.Eq(checkParam)).
					Return(checkRes, nil).
					Times(1)

				b := rest_app.ResponseBody{
					Code:    1003,
					Message: "credential is invalid",
				}
				s.
					EXPECT().
					Marshal(gomock.Eq(b)).
					Return([]byte{}, nil).
					Times(1)
				rw.
					EXPECT().
					WriteHeader(401).
					Times(1)
				rw.
					EXPECT().
					Write(gomock.Eq([]byte{})).
					Times(1)

				m.ServeHTTP(rw, req)
			})
		})

		When("token is valid", func() {
			It("should return result", func() {
				checkParam := auth.CheckCredentialParam{
					AuthToken: "basic-token",
				}
				checkRes := &auth.CheckCredentialResult{
					TokenValid: true,
				}
				a.
					EXPECT().
					CheckCredential(gomock.Any(), gomock.Eq(checkParam)).
					Return(checkRes, nil).
					Times(1)

				handler.
					EXPECT().
					ServeHTTP(gomock.Eq(rw), gomock.Eq(req)).
					Times(1)

				m.ServeHTTP(rw, req)
			})
		})
	})

	Context("NewRequestLogMiddleware", Label("unit"), func() {
		var (
			logger  *mock_logging.MockLogger
			clock   *mock_datetime.MockClock
			handler *mock_restapp.MockHandler
			m       http.Handler

			rw               *mock_restapp.MockResponseWriter
			req              *http.Request
			currentTimestamp time.Time
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			logger = mock_logging.NewMockLogger(ctrl)
			clock = mock_datetime.NewMockClock(ctrl)
			handler = mock_restapp.NewMockHandler(ctrl)
			fn, _ := rest_app.NewRequestLogMiddleware(rest_app.RequestLogMiddlewareParam{
				Logger: logger,
				Clock:  clock,
			})
			m = fn(handler)

			rw = mock_restapp.NewMockResponseWriter(ctrl)
			req = &http.Request{
				Header:     http.Header{},
				Method:     http.MethodPost,
				RequestURI: "/custom-endpoint",
				URL: &url.URL{
					Host: "localhost",
				},
			}
			currentTimestamp = time.Now()
		})

		When("logger is not specified", func() {
			It("should return error", func() {
				fn, err := rest_app.NewRequestLogMiddleware(rest_app.RequestLogMiddlewareParam{
					Logger: nil,
					Clock:  clock,
				})

				Expect(fn).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("logger is not specified")))
			})
		})

		When("ignore uri is specified", func() {
			It("should return result", func() {
				req.RequestURI = "uri"
				fn, _ := rest_app.NewRequestLogMiddleware(rest_app.RequestLogMiddlewareParam{
					Logger: logger,
					Clock:  clock,
					IgnoreURI: map[string]bool{
						"uri": true,
					},
				})
				m = fn(handler)

				handler.
					EXPECT().
					ServeHTTP(gomock.Eq(rw), gomock.Eq(req)).
					Times(1)

				m.ServeHTTP(rw, req)
			})
		})

		When("ignore uri is not specified", func() {
			It("should return result", func() {
				fn, _ := rest_app.NewRequestLogMiddleware(rest_app.RequestLogMiddlewareParam{
					Logger: logger,
					Clock:  clock,
				})
				m = fn(handler)

				clock.
					EXPECT().
					Now().
					Return(currentTimestamp).
					Times(1)

				handler.
					EXPECT().
					ServeHTTP(gomock.Any(), gomock.Eq(req)).
					Times(1)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Error(gomock.Eq("request: POST /custom-endpoint")).
					Times(1)

				m.ServeHTTP(rw, req)
			})
		})

		When("header is specified", func() {
			It("should return result", func() {
				fn, _ := rest_app.NewRequestLogMiddleware(rest_app.RequestLogMiddlewareParam{
					Logger: logger,
					Clock:  clock,
					Header: map[string]string{
						"Authorization": "auth",
					},
				})
				m = fn(handler)

				clock.
					EXPECT().
					Now().
					Return(currentTimestamp).
					Times(1)

				handler.
					EXPECT().
					ServeHTTP(gomock.Any(), gomock.Eq(req)).
					Times(1)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Error(gomock.Eq("request: POST /custom-endpoint")).
					Times(1)

				m.ServeHTTP(rw, req)
			})
		})

		When("clock is not specified", func() {
			It("should return result", func() {
				fn, _ := rest_app.NewRequestLogMiddleware(rest_app.RequestLogMiddlewareParam{
					Logger: logger,
					Clock:  nil,
				})
				m = fn(handler)

				clock.
					EXPECT().
					Now().
					Times(0)

				rw.
					EXPECT().
					Write(gomock.Eq([]byte{})).
					Return(1, nil).
					Times(1)

				handler.
					EXPECT().
					ServeHTTP(gomock.Any(), gomock.Eq(req)).
					Do(func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte{})
					}).
					Times(1)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Info(gomock.Eq("request: POST /custom-endpoint")).
					Times(1)

				m.ServeHTTP(rw, req)
			})
		})

		When("request is success", func() {
			It("should return result", func() {
				fn, _ := rest_app.NewRequestLogMiddleware(rest_app.RequestLogMiddlewareParam{
					Logger: logger,
					Clock:  clock,
				})
				m = fn(handler)

				clock.
					EXPECT().
					Now().
					Return(currentTimestamp).
					Times(1)

				rw.
					EXPECT().
					WriteHeader(gomock.Eq(200)).
					Times(1)

				rw.
					EXPECT().
					Write(gomock.Eq([]byte{})).
					Return(1, nil).
					Times(1)

				handler.
					EXPECT().
					ServeHTTP(gomock.Any(), gomock.Eq(req)).
					Do(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(200)
						w.Write([]byte{})
					}).
					Times(1)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Info(gomock.Eq("request: POST /custom-endpoint")).
					Times(1)

				m.ServeHTTP(rw, req)
			})
		})

		When("client request error", func() {
			It("should return result", func() {
				fn, _ := rest_app.NewRequestLogMiddleware(rest_app.RequestLogMiddlewareParam{
					Logger: logger,
					Clock:  clock,
				})
				m = fn(handler)

				clock.
					EXPECT().
					Now().
					Return(currentTimestamp).
					Times(1)

				rw.
					EXPECT().
					WriteHeader(gomock.Eq(400)).
					Times(1)

				handler.
					EXPECT().
					ServeHTTP(gomock.Any(), gomock.Eq(req)).
					Do(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(400)
					}).
					Times(1)

				logger.
					EXPECT().
					WithFields(gomock.Any()).
					Return(logger).
					Times(1)

				logger.
					EXPECT().
					Warn(gomock.Eq("request: POST /custom-endpoint")).
					Times(1)

				m.ServeHTTP(rw, req)
			})
		})
	})

})
