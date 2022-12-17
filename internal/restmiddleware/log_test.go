package restmiddleware_test

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-seidon/hippo/internal/restmiddleware"
	mock_datetime "github.com/go-seidon/provider/datetime/mock"
	mock_http "github.com/go-seidon/provider/http/mock"
	mock_logging "github.com/go-seidon/provider/logging/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Log Package", func() {

	Context("Handle function", Label("unit"), func() {
		var (
			logger  *mock_logging.MockLogger
			clock   *mock_datetime.MockClock
			handler *mock_http.MockHandler
			m       http.Handler

			rw               *mock_http.MockResponseWriter
			req              *http.Request
			currentTimestamp time.Time
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			logger = mock_logging.NewMockLogger(ctrl)
			clock = mock_datetime.NewMockClock(ctrl)
			handler = mock_http.NewMockHandler(ctrl)
			fn := restmiddleware.NewRequestLog(restmiddleware.RequestLogParam{
				Logger: logger,
				Clock:  clock,
			})
			m = fn.Handle(handler)

			rw = mock_http.NewMockResponseWriter(ctrl)
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

		When("ignore uri is specified", func() {
			It("should return result", func() {
				req.RequestURI = "uri"
				fn := restmiddleware.NewRequestLog(restmiddleware.RequestLogParam{
					Logger: logger,
					Clock:  clock,
					IgnoreURI: map[string]bool{
						"uri": true,
					},
				})
				m = fn.Handle(handler)

				handler.
					EXPECT().
					ServeHTTP(gomock.Eq(rw), gomock.Eq(req)).
					Times(1)

				m.ServeHTTP(rw, req)
			})
		})

		When("ignore uri is not specified", func() {
			It("should return result", func() {
				fn := restmiddleware.NewRequestLog(restmiddleware.RequestLogParam{
					Logger: logger,
					Clock:  clock,
				})
				m = fn.Handle(handler)

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
				fn := restmiddleware.NewRequestLog(restmiddleware.RequestLogParam{
					Logger: logger,
					Clock:  clock,
					Header: map[string]string{
						"Authorization": "auth",
					},
				})
				m = fn.Handle(handler)

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
				fn := restmiddleware.NewRequestLog(restmiddleware.RequestLogParam{
					Logger: logger,
					Clock:  nil,
				})
				m = fn.Handle(handler)

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
				fn := restmiddleware.NewRequestLog(restmiddleware.RequestLogParam{
					Logger: logger,
					Clock:  clock,
				})
				m = fn.Handle(handler)

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
				fn := restmiddleware.NewRequestLog(restmiddleware.RequestLogParam{
					Logger: logger,
					Clock:  clock,
				})
				m = fn.Handle(handler)

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
