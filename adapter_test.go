package errcode

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestNewAppError(t *testing.T) {
	err := New(http.StatusTeapot, Code("TEAPOT"), "short and stout")

	if err.StatusCode != http.StatusTeapot {
		t.Fatalf("expected status %d, got %d", http.StatusTeapot, err.StatusCode)
	}
	if err.ErrorCode != Code("TEAPOT") {
		t.Fatalf("expected code %q, got %q", Code("TEAPOT"), err.ErrorCode)
	}
	if err.Message != "short and stout" {
		t.Fatalf("expected message %q, got %q", "short and stout", err.Message)
	}
	if err.Status() != http.StatusTeapot {
		t.Fatalf("expected Status() %d, got %d", http.StatusTeapot, err.Status())
	}
	if err.ErrCode() != "TEAPOT" {
		t.Fatalf("expected ErrCode() %q, got %q", "TEAPOT", err.ErrCode())
	}
	if got := err.Error(); got != "TEAPOT: short and stout" {
		t.Fatalf("expected Error() %q, got %q", "TEAPOT: short and stout", got)
	}
}

func TestWithMessage(t *testing.T) {
	err := NotFound.WithMessage("product %s", "not found")

	if err == NotFound {
		t.Fatal("expected WithMessage to return a new AppError")
	}
	if err.StatusCode != NotFound.StatusCode {
		t.Fatalf("expected status %d, got %d", NotFound.StatusCode, err.StatusCode)
	}
	if err.ErrorCode != NotFound.ErrorCode {
		t.Fatalf("expected code %q, got %q", NotFound.ErrorCode, err.ErrorCode)
	}
	if err.Message != "product not found" {
		t.Fatalf("expected message %q, got %q", "product not found", err.Message)
	}
	if NotFound.Message != "Not Found" {
		t.Fatalf("expected sentinel message to remain %q, got %q", "Not Found", NotFound.Message)
	}
	if !errors.Is(err, NotFound) {
		t.Fatal("expected customized error to match NotFound")
	}
	if errors.Is(err, BadRequest) {
		t.Fatal("expected customized error not to match BadRequest")
	}
}

func TestWithCause(t *testing.T) {
	err := InternalServerError.WithCause(sql.ErrConnDone)

	if err == InternalServerError {
		t.Fatal("expected WithCause to return a new AppError")
	}
	if err.Cause != sql.ErrConnDone {
		t.Fatalf("expected cause %v, got %v", sql.ErrConnDone, err.Cause)
	}
	if !errors.Is(err, sql.ErrConnDone) {
		t.Fatal("expected error to unwrap to sql.ErrConnDone")
	}
	if !errors.Is(err, InternalServerError) {
		t.Fatal("expected error to match InternalServerError")
	}
	if !strings.Contains(err.Error(), sql.ErrConnDone.Error()) {
		t.Fatalf("expected Error() to include cause, got %q", err.Error())
	}
}

func TestErrorsAs(t *testing.T) {
	err := BadRequest.WithMessage("name is required")

	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Fatal("expected errors.As to extract *AppError")
	}
	if appErr.ErrCode() != string(CodeBadRequest) {
		t.Fatalf("expected code %q, got %q", CodeBadRequest, appErr.ErrCode())
	}
}

func TestCauseExcludedFromJSON(t *testing.T) {
	err := InternalServerError.WithCause(sql.ErrConnDone)

	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("expected json marshal to succeed: %v", marshalErr)
	}

	jsonBody := string(data)
	if strings.Contains(jsonBody, "Cause") || strings.Contains(jsonBody, sql.ErrConnDone.Error()) {
		t.Fatalf("expected JSON to exclude cause, got %s", jsonBody)
	}
	if !strings.Contains(jsonBody, string(CodeInternalServerError)) {
		t.Fatalf("expected JSON to include error code, got %s", jsonBody)
	}
}

func TestNilAppError(t *testing.T) {
	var err *AppError

	if err.Error() != "" {
		t.Fatalf("expected nil Error() to return empty string, got %q", err.Error())
	}
	if err.Unwrap() != nil {
		t.Fatalf("expected nil Unwrap() to return nil, got %v", err.Unwrap())
	}
	if err.Status() != http.StatusInternalServerError {
		t.Fatalf("expected nil Status() to return 500, got %d", err.Status())
	}
	if err.ErrCode() != string(CodeInternalServerError) {
		t.Fatalf("expected nil ErrCode() to return INTERNAL_SERVER_ERROR, got %q", err.ErrCode())
	}
}
