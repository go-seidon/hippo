package rest_app

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/go-seidon/local/internal/file"
	"github.com/go-seidon/local/internal/healthcheck"
	"github.com/go-seidon/local/internal/logging"
	"github.com/go-seidon/local/internal/serialization"
	"github.com/go-seidon/local/internal/status"
	"github.com/go-seidon/local/internal/validation"
	"github.com/gorilla/mux"
)

func NewNotFoundHandler(log logging.Logger, s serialization.Serializer) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		Response(
			WithWriterSerializer(w, s),
			WithHttpCode(http.StatusNotFound),
			WithCode(status.RESOURCE_NOTFOUND),
			WithMessage("resource not found"),
		)
	}
}

func NewMethodNotAllowedHandler(log logging.Logger, s serialization.Serializer) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		Response(
			WithWriterSerializer(w, s),
			WithHttpCode(http.StatusMethodNotAllowed),
			WithCode(status.ACTION_FAILED),
			WithMessage("method is not allowed"),
		)
	}
}

func NewRootHandler(log logging.Logger, s serialization.Serializer, config *RestAppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		d := struct {
			AppName    string `json:"app_name"`
			AppVersion string `json:"app_version"`
		}{
			AppName:    config.AppName,
			AppVersion: config.AppVersion,
		}
		Response(WithWriterSerializer(w, s), WithData(d))
	}
}

func NewHealthCheckHandler(log logging.Logger, s serialization.Serializer, healthService healthcheck.HealthCheck) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r, err := healthService.Check()
		if err != nil {

			Response(
				WithWriterSerializer(w, s),
				WithCode(status.ACTION_FAILED),
				WithMessage(err.Error()),
				WithHttpCode(http.StatusBadRequest),
			)
			return
		}

		jobs := map[string]struct {
			Name      string `json:"name"`
			Status    string `json:"status"`
			CheckedAt int64  `json:"checked_at"`
			Error     string `json:"error"`
		}{}
		for jobName, item := range r.Items {
			jobs[jobName] = struct {
				Name      string `json:"name"`
				Status    string `json:"status"`
				CheckedAt int64  `json:"checked_at"`
				Error     string `json:"error"`
			}{
				Name:      item.Name,
				Status:    item.Status,
				Error:     item.Error,
				CheckedAt: item.CheckedAt.UnixMilli(),
			}
		}

		d := struct {
			Status  string `json:"status"`
			Details map[string]struct {
				Name      string `json:"name"`
				Status    string `json:"status"`
				CheckedAt int64  `json:"checked_at"`
				Error     string `json:"error"`
			} `json:"details"`
		}{
			Status:  r.Status,
			Details: jobs,
		}

		Response(
			WithWriterSerializer(w, s),
			WithData(d),
			WithMessage("success check service health"),
		)
	}
}

func NewDeleteFileHandler(log logging.Logger, s serialization.Serializer, deleter file.File) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		r, err := deleter.DeleteFile(req.Context(), file.DeleteFileParam{
			FileId: vars["id"],
		})
		if err != nil {
			code := status.ACTION_FAILED
			httpCode := http.StatusBadRequest
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
				WithWriterSerializer(w, s),
				WithCode(code),
				WithMessage(message),
				WithHttpCode(httpCode),
			)
			return
		}

		d := struct {
			DeletedAt int64 `json:"deleted_at"`
		}{
			DeletedAt: r.DeletedAt.UnixMilli(),
		}

		Response(
			WithWriterSerializer(w, s),
			WithData(d),
			WithMessage("success delete file"),
		)
	}
}

func NewRetrieveFileHandler(log logging.Logger, s serialization.Serializer, retriever file.File) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		r, err := retriever.RetrieveFile(req.Context(), file.RetrieveFileParam{
			FileId: vars["id"],
		})
		if err != nil {
			code := status.ACTION_FAILED
			httpCode := http.StatusBadRequest
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
				WithWriterSerializer(w, s),
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
				WithWriterSerializer(w, s),
				WithCode(status.ACTION_FAILED),
				WithMessage(err.Error()),
				WithHttpCode(http.StatusBadRequest),
			)
			return
		}

		if r.MimeType != "" {
			w.Header().Set("Content-Type", r.MimeType)
		} else {
			w.Header().Del("Content-Type")
		}

		w.Write(data.Bytes())
	}
}

func NewUploadFileHandler(log logging.Logger, s serialization.Serializer, uploader file.File, config *RestAppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		// set form max size + add 1KB (non file size estimation if any)
		req.Body = http.MaxBytesReader(w, req.Body, config.UploadFormSize+1024)

		fileReader, fileHeader, err := req.FormFile("file")
		if err != nil {
			Response(
				WithWriterSerializer(w, s),
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
				WithWriterSerializer(w, s),
				WithCode(status.ACTION_FAILED),
				WithMessage(err.Error()),
				WithHttpCode(http.StatusBadRequest),
			)
			return
		}

		uploadRes, err := uploader.UploadFile(req.Context(),
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
			httpCode := http.StatusBadRequest
			message := err.Error()

			switch e := err.(type) {
			case *validation.ValidationError:
				code = status.INVALID_PARAM
				httpCode = http.StatusBadRequest
				message = e.Error()
			}

			Response(
				WithWriterSerializer(w, s),
				WithCode(code),
				WithMessage(message),
				WithHttpCode(httpCode),
			)
			return
		}

		d := struct {
			UniqueId   string `json:"id"`
			Name       string `json:"name"`
			Mimetype   string `json:"mimetype"`
			Extension  string `json:"extension"`
			Size       int64  `json:"size"`
			UploadedAt int64  `json:"uploaded_at"`
		}{
			UniqueId:   uploadRes.UniqueId,
			Name:       uploadRes.Name,
			Mimetype:   uploadRes.Mimetype,
			Extension:  uploadRes.Extension,
			Size:       uploadRes.Size,
			UploadedAt: uploadRes.UploadedAt.UnixMilli(),
		}

		Response(
			WithWriterSerializer(w, s),
			WithData(d),
			WithMessage("success upload file"),
		)
	}
}
