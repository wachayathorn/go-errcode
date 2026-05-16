package errcode

import (
	"fmt"
)

type Code string

type AppError struct {
	StatusCode int    `json:"status_code" example:"500"`
	ErrorCode  Code   `json:"error_code" example:"INTERNAL_SERVER_ERROR"`
	Message    string `json:"message" example:"Something went wrong"`
	Cause      error  `json:"-"`
}

func New(statusCode int, errorCode Code, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func (e *AppError) WithMessage(format string, a ...interface{}) *AppError {
	return &AppError{
		StatusCode: e.StatusCode,
		ErrorCode:  e.ErrorCode,
		Message:    fmt.Sprintf(format, a...),
		Cause:      e.Cause,
	}
}

func (e *AppError) WithCause(cause error) *AppError {
	return &AppError{
		StatusCode: e.StatusCode,
		ErrorCode:  e.ErrorCode,
		Message:    e.Message,
		Cause:      cause,
	}
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message == "" {
		return string(e.ErrorCode)
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.ErrorCode, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.ErrorCode, e.Message)
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func (e *AppError) Is(target error) bool {
	targetErr, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.ErrorCode == targetErr.ErrorCode
}

func (e *AppError) Status() int {
	return e.StatusCode
}

func (e *AppError) ErrCode() string {
	return string(e.ErrorCode)
}
