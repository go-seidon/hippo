package resthandler

import (
	"fmt"
	"net/http"

	"github.com/go-seidon/hippo/api/restapp"
	"github.com/go-seidon/hippo/internal/file"
	"github.com/go-seidon/hippo/internal/storage/multipart"
	"github.com/go-seidon/provider/status"
	"github.com/labstack/echo/v4"
)

type fileHandler struct {
	fileClient file.File
	fileParser multipart.Parser
}

func (h *fileHandler) UploadFile(ctx echo.Context) error {
	fileHeader, ferr := ctx.FormFile("file")
	if ferr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &restapp.ResponseBodyInfo{
			Code:    status.INVALID_PARAM,
			Message: ferr.Error(),
		})
	}

	fileInfo, ferr := h.fileParser(fileHeader)
	if ferr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &restapp.ResponseBodyInfo{
			Code:    status.INVALID_PARAM,
			Message: ferr.Error(),
		})
	}

	uploadFile, err := h.fileClient.UploadFile(
		ctx.Request().Context(),
		file.WithReader(fileInfo.Data),
		file.WithFileInfo(
			fileInfo.Name,
			fileInfo.Mimetype,
			fileInfo.Extension,
			fileInfo.Size,
		),
	)
	if err != nil {
		httpCode := http.StatusInternalServerError
		switch err.Code {
		case status.INVALID_PARAM:
			httpCode = http.StatusBadRequest
		}
		return echo.NewHTTPError(httpCode, &restapp.ResponseBodyInfo{
			Code:    err.Code,
			Message: err.Message,
		})
	}

	return ctx.JSON(http.StatusOK, &restapp.UploadFileResponse{
		Code:    uploadFile.Success.Code,
		Message: uploadFile.Success.Message,
		Data: restapp.UploadFileData{
			Id:         uploadFile.UniqueId,
			Name:       uploadFile.Name,
			Mimetype:   uploadFile.Mimetype,
			Extension:  uploadFile.Extension,
			Size:       uploadFile.Size,
			UploadedAt: uploadFile.UploadedAt.UnixMilli(),
		},
	})
}

func (h *fileHandler) RetrieveFileById(ctx echo.Context) error {
	findFile, err := h.fileClient.RetrieveFile(ctx.Request().Context(), file.RetrieveFileParam{
		FileId: ctx.Param("id"),
	})
	if err != nil {
		httpCode := http.StatusInternalServerError
		switch err.Code {
		case status.INVALID_PARAM:
			httpCode = http.StatusBadRequest
		case status.RESOURCE_NOTFOUND:
			httpCode = http.StatusNotFound
		}
		return echo.NewHTTPError(httpCode, &restapp.ResponseBodyInfo{
			Code:    err.Code,
			Message: err.Message,
		})
	}

	header := ctx.Response().Header()
	header.Set("X-File-Name", findFile.Name)
	header.Set("X-File-Mimetype", findFile.MimeType)
	header.Set("X-File-Extension", findFile.Extension)
	header.Set("X-File-Size", fmt.Sprintf("%d", findFile.Size))
	return ctx.Stream(http.StatusOK, findFile.MimeType, findFile.Data)
}

func (h *fileHandler) DeleteFileById(ctx echo.Context) error {
	deleteFile, err := h.fileClient.DeleteFile(ctx.Request().Context(), file.DeleteFileParam{
		FileId: ctx.Param("id"),
	})
	if err != nil {
		httpCode := http.StatusInternalServerError
		switch err.Code {
		case status.INVALID_PARAM:
			httpCode = http.StatusBadRequest
		case status.RESOURCE_NOTFOUND:
			httpCode = http.StatusNotFound
		}
		return echo.NewHTTPError(httpCode, &restapp.ResponseBodyInfo{
			Code:    err.Code,
			Message: err.Message,
		})
	}

	return ctx.JSON(http.StatusOK, &restapp.DeleteFileByIdResponse{
		Code:    deleteFile.Success.Code,
		Message: deleteFile.Success.Message,
		Data: restapp.DeleteFileByIdData{
			DeletedAt: deleteFile.DeletedAt.UnixMilli(),
		},
	})
}

type FileParam struct {
	FileClient file.File
	FileParser multipart.Parser
}

func NewFile(p FileParam) *fileHandler {
	return &fileHandler{
		fileClient: p.FileClient,
		fileParser: p.FileParser,
	}
}
