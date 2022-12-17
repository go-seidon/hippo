// Package restapp provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package restapp

import (
	"encoding/json"
	"fmt"
)

const (
	BasicAuthScopes = "basicAuth.Scopes"
)

// CheckHealthData defines model for CheckHealthData.
type CheckHealthData struct {
	Details CheckHealthData_Details `json:"details"`
	Status  string                  `json:"status"`
}

// CheckHealthData_Details defines model for CheckHealthData.Details.
type CheckHealthData_Details struct {
	AdditionalProperties map[string]CheckHealthDetail `json:"-"`
}

// CheckHealthDetail defines model for CheckHealthDetail.
type CheckHealthDetail struct {
	CheckedAt int64  `json:"checked_at"`
	Error     string `json:"error"`
	Name      string `json:"name"`
	Status    string `json:"status"`
}

// CheckHealthResponse defines model for CheckHealthResponse.
type CheckHealthResponse struct {
	Code    int32           `json:"code"`
	Data    CheckHealthData `json:"data"`
	Message string          `json:"message"`
}

// DeleteFileByIdData defines model for DeleteFileByIdData.
type DeleteFileByIdData struct {
	DeletedAt int64 `json:"deleted_at"`
}

// DeleteFileByIdResponse defines model for DeleteFileByIdResponse.
type DeleteFileByIdResponse struct {
	Code    int32              `json:"code"`
	Data    DeleteFileByIdData `json:"data"`
	Message string             `json:"message"`
}

// GetAppInfoData defines model for GetAppInfoData.
type GetAppInfoData struct {
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
}

// GetAppInfoResponse defines model for GetAppInfoResponse.
type GetAppInfoResponse struct {
	Code    int32          `json:"code"`
	Data    GetAppInfoData `json:"data"`
	Message string         `json:"message"`
}

// ResponseBodyInfo defines model for ResponseBodyInfo.
type ResponseBodyInfo struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// RetrieveFileByIdResponse defines model for RetrieveFileByIdResponse.
type RetrieveFileByIdResponse = string

// UploadFileData defines model for UploadFileData.
type UploadFileData struct {
	Extension  string `json:"extension"`
	Id         string `json:"id"`
	Mimetype   string `json:"mimetype"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	UploadedAt int64  `json:"uploaded_at"`
}

// UploadFileRequest defines model for UploadFileRequest.
type UploadFileRequest struct {
	File string `json:"file"`
}

// UploadFileResponse defines model for UploadFileResponse.
type UploadFileResponse struct {
	Code    int32          `json:"code"`
	Data    UploadFileData `json:"data"`
	Message string         `json:"message"`
}

// CorrelationId defines model for CorrelationId.
type CorrelationId = string

// ObjectId defines model for ObjectId.
type ObjectId = string

// BadRequest defines model for BadRequest.
type BadRequest = ResponseBodyInfo

// NotFound defines model for NotFound.
type NotFound = ResponseBodyInfo

// ServerError defines model for ServerError.
type ServerError = ResponseBodyInfo

// UnauthenticatedAccess defines model for UnauthenticatedAccess.
type UnauthenticatedAccess = ResponseBodyInfo

// GetAppInfoParams defines parameters for GetAppInfo.
type GetAppInfoParams struct {
	// correlation id for tracing purposes
	XCorrelationId *CorrelationId `json:"X-Correlation-Id,omitempty"`
}

// CheckHealthParams defines parameters for CheckHealth.
type CheckHealthParams struct {
	// correlation id for tracing purposes
	XCorrelationId *CorrelationId `json:"X-Correlation-Id,omitempty"`
}

// UploadFileParams defines parameters for UploadFile.
type UploadFileParams struct {
	// correlation id for tracing purposes
	XCorrelationId *CorrelationId `json:"X-Correlation-Id,omitempty"`
}

// DeleteFileByIdParams defines parameters for DeleteFileById.
type DeleteFileByIdParams struct {
	// correlation id for tracing purposes
	XCorrelationId *CorrelationId `json:"X-Correlation-Id,omitempty"`
}

// RetrieveFileByIdParams defines parameters for RetrieveFileById.
type RetrieveFileByIdParams struct {
	// correlation id for tracing purposes
	XCorrelationId *CorrelationId `json:"X-Correlation-Id,omitempty"`
}

// Getter for additional properties for CheckHealthData_Details. Returns the specified
// element and whether it was found
func (a CheckHealthData_Details) Get(fieldName string) (value CheckHealthDetail, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for CheckHealthData_Details
func (a *CheckHealthData_Details) Set(fieldName string, value CheckHealthDetail) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]CheckHealthDetail)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for CheckHealthData_Details to handle AdditionalProperties
func (a *CheckHealthData_Details) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]CheckHealthDetail)
		for fieldName, fieldBuf := range object {
			var fieldVal CheckHealthDetail
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for CheckHealthData_Details to handle AdditionalProperties
func (a CheckHealthData_Details) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}
