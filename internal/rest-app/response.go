package rest_app

import (
	"fmt"
	"net/http"

	"github.com/go-seidon/hippo/internal/status"
	"github.com/go-seidon/provider/serialization"
)

type ResponseBody struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseParam struct {
	Writer     http.ResponseWriter
	Serializer serialization.Serializer

	Body            ResponseBody
	defaultHttpCode int
	HttpCode        int
	Message         string
	Code            int32
	Data            interface{}
}

type ResponseOption = func(*ResponseParam)

func WithWriterSerializer(w http.ResponseWriter, s serialization.Serializer) ResponseOption {
	return func(rp *ResponseParam) {
		rp.Writer = w
		rp.Serializer = s
	}
}

func WithHttpCode(c int) ResponseOption {
	return func(rp *ResponseParam) {
		rp.HttpCode = c
	}
}

func WithMessage(m string) ResponseOption {
	return func(rp *ResponseParam) {
		rp.Message = m
	}
}

func WithCode(c int32) ResponseOption {
	return func(rp *ResponseParam) {
		rp.Code = c
	}
}

func WithData(d interface{}) ResponseOption {
	return func(rp *ResponseParam) {
		rp.Data = d
	}
}

func Response(opts ...ResponseOption) error {
	p := ResponseParam{
		Body: ResponseBody{
			Code:    status.ACTION_SUCCESS,
			Message: "success",
		},
		defaultHttpCode: http.StatusOK,
	}
	for _, opt := range opts {
		opt(&p)
	}

	if p.Writer == nil {
		return fmt.Errorf("writer should be specified")
	}
	if p.Serializer == nil {
		return fmt.Errorf("serializer should be specified")
	}

	httpCode := p.defaultHttpCode
	if p.HttpCode != 0 {
		httpCode = p.HttpCode
	}

	if p.Message != "" {
		p.Body.Message = p.Message
	}

	if p.Code != 0 {
		p.Body.Code = p.Code
	}

	if p.Data != nil {
		p.Body.Data = p.Data
	}

	r, err := p.Serializer.Marshal(p.Body)
	if err != nil {
		return err
	}

	p.Writer.WriteHeader(httpCode)
	p.Writer.Write(r)
	return nil
}
