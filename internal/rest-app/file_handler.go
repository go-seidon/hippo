package rest_app

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	rest_v1 "github.com/go-seidon/hippo/generated/rest-v1"
	"github.com/go-seidon/hippo/internal/file"
	"github.com/go-seidon/hippo/internal/logging"
	"github.com/go-seidon/hippo/internal/serialization"
	"github.com/go-seidon/hippo/internal/status"
	"github.com/go-seidon/hippo/internal/validation"
	"github.com/gorilla/mux"
)

type fileHandler struct {
	logger      logging.Logger
	serializer  serialization.Serializer
	config      *RestAppConfig
	fileService file.File
}

func (h *fileHandler) UploadFile(w http.ResponseWriter, req *http.Request) {
	// set form max size + add 1KB (non file size estimation if any)
	req.Body = http.MaxBytesReader(w, req.Body, h.config.UploadFormSize+1024)

	fileReader, fileHeader, err := req.FormFile("file")
	if err != nil {
		Response(
			WithWriterSerializer(w, h.serializer),
			WithCode(status.ACTION_FAILED),
			WithMessage(err.Error()),
			WithHttpCode(http.StatusBadRequest),
		)
		return
	}
	defer fileReader.Close()

	fileInfo, err := ParseMultipartFile(fileReader, fileHeader)
	if err != nil {
		Response(
			WithWriterSerializer(w, h.serializer),
			WithCode(status.ACTION_FAILED),
			WithMessage(err.Error()),
			WithHttpCode(http.StatusBadRequest),
		)
		return
	}

	uploadRes, err := h.fileService.UploadFile(req.Context(),
		file.WithReader(fileReader),
		file.WithFileInfo(
			fileInfo.Name,
			fileInfo.Mimetype,
			fileInfo.Extension,
			fileInfo.Size,
		),
	)
	if err != nil {
		code := status.ACTION_FAILED
		httpCode := http.StatusInternalServerError
		message := err.Error()

		switch e := err.(type) {
		case *validation.ValidationError:
			code = status.INVALID_PARAM
			httpCode = http.StatusBadRequest
			message = e.Error()
		}

		Response(
			WithWriterSerializer(w, h.serializer),
			WithCode(code),
			WithMessage(message),
			WithHttpCode(httpCode),
		)
		return
	}

	d := &rest_v1.UploadFileData{
		Id:         uploadRes.UniqueId,
		Name:       uploadRes.Name,
		Mimetype:   uploadRes.Mimetype,
		Extension:  uploadRes.Extension,
		Size:       uploadRes.Size,
		UploadedAt: uploadRes.UploadedAt.UnixMilli(),
	}

	Response(
		WithWriterSerializer(w, h.serializer),
		WithData(d),
		WithMessage("success upload file"),
	)
}

func (h *fileHandler) RetrieveFileById(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	r, err := h.fileService.RetrieveFile(req.Context(), file.RetrieveFileParam{
		FileId: vars["id"],
	})
	if err != nil {
		code := status.ACTION_FAILED
		httpCode := http.StatusInternalServerError
		message := err.Error()

		switch e := err.(type) {
		case *validation.ValidationError:
			code = status.INVALID_PARAM
			httpCode = http.StatusBadRequest
			message = e.Error()
		}

		if errors.Is(err, file.ErrorNotFound) {
			code = status.RESOURCE_NOTFOUND
			httpCode = http.StatusNotFound
			message = err.Error()
		}

		Response(
			WithWriterSerializer(w, h.serializer),
			WithCode(code),
			WithMessage(message),
			WithHttpCode(httpCode),
		)
		return
	}
	defer r.Data.Close()

	data := bytes.NewBuffer([]byte{})
	_, err = io.Copy(data, r.Data)
	if err != nil {
		Response(
			WithWriterSerializer(w, h.serializer),
			WithCode(status.ACTION_FAILED),
			WithMessage(err.Error()),
			WithHttpCode(http.StatusInternalServerError),
		)
		return
	}

	if r.MimeType != "" {
		w.Header().Set("Content-Type", r.MimeType)
	} else {
		w.Header().Del("Content-Type")
	}

	w.Header().Set("X-File-Name", r.Name)
	w.Header().Set("X-File-Mimetype", r.MimeType)
	w.Header().Set("X-File-Extension", r.Extension)
	w.Header().Set("X-File-Size", fmt.Sprintf("%d", r.Size))
	w.Write(data.Bytes())
}

func (h *fileHandler) DeleteFileById(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	r, err := h.fileService.DeleteFile(req.Context(), file.DeleteFileParam{
		FileId: vars["id"],
	})
	if err != nil {
		code := status.ACTION_FAILED
		httpCode := http.StatusInternalServerError
		message := err.Error()

		if errors.Is(err, file.ErrorNotFound) {
			code = status.RESOURCE_NOTFOUND
			httpCode = http.StatusNotFound
			message = err.Error()
		}

		switch e := err.(type) {
		case *validation.ValidationError:
			code = status.INVALID_PARAM
			httpCode = http.StatusBadRequest
			message = e.Error()
		}

		Response(
			WithWriterSerializer(w, h.serializer),
			WithCode(code),
			WithMessage(message),
			WithHttpCode(httpCode),
		)
		return
	}

	d := &rest_v1.DeleteFileData{
		DeletedAt: r.DeletedAt.UnixMilli(),
	}

	Response(
		WithWriterSerializer(w, h.serializer),
		WithData(d),
		WithMessage("success delete file"),
	)
}

type FileHandlerParam struct {
	Logger      logging.Logger
	Serializer  serialization.Serializer
	Config      *RestAppConfig
	FileService file.File
}

func NewFileHandler(p FileHandlerParam) *fileHandler {
	return &fileHandler{
		logger:      p.Logger,
		serializer:  p.Serializer,
		config:      p.Config,
		fileService: p.FileService,
	}
}
