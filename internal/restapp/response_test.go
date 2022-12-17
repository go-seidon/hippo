package restapp_test

import (
	"fmt"
	"net/http"

	"github.com/go-seidon/hippo/internal/restapp"
	mock_restapp "github.com/go-seidon/hippo/internal/restapp/mock"
	mock_serialization "github.com/go-seidon/provider/serialization/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Response Package", func() {

	Context("WithWriterSerializer function", Label("unit"), func() {
		var (
			w *mock_restapp.MockResponseWriter
			s *mock_serialization.MockSerializer
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			w = mock_restapp.NewMockResponseWriter(ctrl)
			s = mock_serialization.NewMockSerializer(ctrl)
		})

		When("writer is not specified", func() {
			It("should not set writer", func() {
				p := restapp.ResponseParam{}
				opt := restapp.WithWriterSerializer(nil, s)

				opt(&p)

				Expect(p.Writer).To(BeNil())
				Expect(p.Serializer).To(Equal(s))
			})
		})

		When("serializer is not specified", func() {
			It("should not set serializer", func() {
				p := restapp.ResponseParam{}
				opt := restapp.WithWriterSerializer(w, nil)

				opt(&p)

				Expect(p.Writer).To(Equal(w))
				Expect(p.Serializer).To(BeNil())
			})
		})

		When("writer and serializer are specified", func() {
			It("should set writer and serializer", func() {
				p := restapp.ResponseParam{}
				opt := restapp.WithWriterSerializer(w, s)

				opt(&p)

				Expect(p.Writer).To(Equal(w))
				Expect(p.Serializer).To(Equal(s))
			})
		})
	})

	Context("WithHttpCode function", Label("unit"), func() {
		var (
			code int
		)

		BeforeEach(func() {
			code = http.StatusCreated
		})

		When("http code is specified", func() {
			It("should set http code", func() {
				p := restapp.ResponseParam{}
				opt := restapp.WithHttpCode(code)

				opt(&p)

				Expect(p.HttpCode).To(Equal(code))
			})
		})
	})

	Context("WithMessage function", Label("unit"), func() {
		var (
			message string
		)

		BeforeEach(func() {
			message = "success do something"
		})

		When("message is specified", func() {
			It("should set http", func() {
				p := restapp.ResponseParam{}
				opt := restapp.WithMessage(message)

				opt(&p)

				Expect(p.Message).To(Equal(message))
			})
		})
	})

	Context("WithCode function", Label("unit"), func() {
		var (
			code int32
		)

		BeforeEach(func() {
			code = 1001
		})

		When("code is specified", func() {
			It("should set http", func() {
				p := restapp.ResponseParam{}
				opt := restapp.WithCode(code)

				opt(&p)

				Expect(p.Code).To(Equal(code))
			})
		})
	})

	Context("WithData function", Label("unit"), func() {
		var (
			data interface{}
		)

		BeforeEach(func() {
			data = struct{ Key string }{Key: "key"}
		})

		When("data is specified", func() {
			It("should set http", func() {
				p := restapp.ResponseParam{}
				opt := restapp.WithData(data)

				opt(&p)

				Expect(p.Data).To(Equal(data))
			})
		})
	})

	Context("Response function", Label("unit"), func() {
		var (
			w *mock_restapp.MockResponseWriter
			s *mock_serialization.MockSerializer
		)

		BeforeEach(func() {
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			w = mock_restapp.NewMockResponseWriter(ctrl)
			s = mock_serialization.NewMockSerializer(ctrl)
		})

		When("writer is not specified", func() {
			It("should return error", func() {
				err := restapp.Response(
					restapp.WithWriterSerializer(nil, s),
				)

				Expect(err).To(Equal(fmt.Errorf("writer should be specified")))
			})
		})

		When("serializer is not specified", func() {
			It("should return error", func() {
				err := restapp.Response(
					restapp.WithWriterSerializer(w, nil),
				)

				Expect(err).To(Equal(fmt.Errorf("serializer should be specified")))
			})
		})

		When("http code is specified", func() {
			It("should return nil", func() {
				s.
					EXPECT().
					Marshal(gomock.Any()).
					Return([]byte("mock"), nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(201)).
					Times(1)

				w.
					EXPECT().
					Write(gomock.Eq([]byte("mock"))).
					Times(1)

				err := restapp.Response(
					restapp.WithWriterSerializer(w, s),
					restapp.WithHttpCode(201),
				)

				Expect(err).To(BeNil())
			})
		})

		When("message is specified", func() {
			It("should return nil", func() {
				b := restapp.ResponseBody{
					Code:    1000,
					Message: "success do something",
				}
				s.
					EXPECT().
					Marshal(gomock.Eq(b)).
					Return([]byte("mock"), nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(200)).
					Times(1)

				w.
					EXPECT().
					Write(gomock.Eq([]byte("mock"))).
					Times(1)

				err := restapp.Response(
					restapp.WithWriterSerializer(w, s),
					restapp.WithMessage("success do something"),
				)

				Expect(err).To(BeNil())
			})
		})

		When("failed encode data", func() {
			It("should return error", func() {
				s.
					EXPECT().
					Marshal(gomock.Any()).
					Return(nil, fmt.Errorf("failed encode")).
					Times(1)

				err := restapp.Response(
					restapp.WithWriterSerializer(w, s),
				)

				Expect(err).To(Equal(fmt.Errorf("failed encode")))
			})
		})

		When("success encode data", func() {
			It("should return nil", func() {
				s.
					EXPECT().
					Marshal(gomock.Any()).
					Return([]byte("mock"), nil).
					Times(1)

				w.
					EXPECT().
					WriteHeader(gomock.Eq(200)).
					Times(1)

				w.
					EXPECT().
					Write(gomock.Eq([]byte("mock"))).
					Times(1)

				err := restapp.Response(
					restapp.WithWriterSerializer(w, s),
				)

				Expect(err).To(BeNil())
			})
		})

	})

})