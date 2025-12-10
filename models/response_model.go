package models

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type ErrorCode string

const (
	ErrInsufficientStock ErrorCode = "INSUFFICIENT_STOCK"
	ErrItemNotFound      ErrorCode = "ITEM_NOT_FOUND"
	ErrValidationError   ErrorCode = "VALIDATION_ERROR"
	ErrInternalError     ErrorCode = "INTERNAL_ERROR"
	ErrUnauthorized      ErrorCode = "UNAUTHORIZED"
	ErrForbidden         ErrorCode = "FORBIDDEN"
)

type ErrorResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Code    ErrorCode `json:"code,omitempty"`
}

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func SuccessResponseWithMeta(message string, data interface{}, meta *Meta) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

func ErrorResponseMsg(message string, code ErrorCode) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Message: message,
		Code:    code,
	}
}
