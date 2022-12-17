package restapp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/go-seidon/hippo/api/restapp"
	"github.com/go-seidon/hippo/internal/file"
	"github.com/go-seidon/provider/logging"
	"github.com/go-seidon/provider/serialization"
	"github.com/go-seidon/provider/status"
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

	uploadRes, uerr := h.fileService.UploadFile(req.Context(),
		file.WithReader(fileReader),
		file.WithFileInfo(
			fileInfo.Name,
			fileInfo.Mimetype,
			fileInfo.Extension,
			fileInfo.Size,
		),
	)
	if uerr != nil {
		httpCode := http.StatusInternalServerError
		code := uerr.Code
		message := uerr.Message

		switch uerr.Code {
		case status.INVALID_PARAM:
			httpCode = http.StatusBadRequest
		}

		Response(
			WithWriterSerializer(w, h.serializer),
			WithCode(code),
			WithMessage(message),
			WithHttpCode(httpCode),
		)
		return
	}

	d := &restapp.UploadFileData{
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

	retrieve, rerr := h.fileService.RetrieveFile(req.Context(), file.RetrieveFileParam{
		FileId: vars["id"],
	})
	if rerr != nil {
		httpCode := http.StatusInternalServerError
		code := rerr.Code
		message := rerr.Message

		switch rerr.Code {
		case status.INVALID_PARAM:
			httpCode = http.StatusBadRequest
		case status.RESOURCE_NOTFOUND:
			httpCode = http.StatusNotFound
		}

		Response(
			WithWriterSerializer(w, h.serializer),
			WithCode(code),
			WithMessage(message),
			WithHttpCode(httpCode),
		)
		return
	}
	defer retrieve.Data.Close()

	data := bytes.NewBuffer([]byte{})
	_, err := io.Copy(data, retrieve.Data)
	if err != nil {
		Response(
			WithWriterSerializer(w, h.serializer),
			WithCode(status.ACTION_FAILED),
			WithMessage(err.Error()),
			WithHttpCode(http.StatusInternalServerError),
		)
		return
	}

	if retrieve.MimeType != "" {
		w.Header().Set("Content-Type", retrieve.MimeType)
	} else {
		w.Header().Del("Content-Type")
	}

	w.Header().Set("X-File-Name", retrieve.Name)
	w.Header().Set("X-File-Mimetype", retrieve.MimeType)
	w.Header().Set("X-File-Extension", retrieve.Extension)
	w.Header().Set("X-File-Size", fmt.Sprintf("%d", retrieve.Size))
	w.Write(data.Bytes())
}

func (h *fileHandler) DeleteFileById(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	delete, derr := h.fileService.DeleteFile(req.Context(), file.DeleteFileParam{
		FileId: vars["id"],
	})
	if derr != nil {
		httpCode := http.StatusInternalServerError
		code := derr.Code
		message := derr.Message

		switch derr.Code {
		case status.INVALID_PARAM:
			httpCode = http.StatusBadRequest
		case status.RESOURCE_NOTFOUND:
			httpCode = http.StatusNotFound
		}

		Response(
			WithWriterSerializer(w, h.serializer),
			WithCode(code),
			WithMessage(message),
			WithHttpCode(httpCode),
		)
		return
	}

	d := &restapp.DeleteFileData{
		DeletedAt: delete.DeletedAt.UnixMilli(),
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
