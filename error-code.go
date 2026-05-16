package errcode

import "net/http"

const (
	CodeBadRequest          Code = "BAD_REQUEST"
	CodeDuplicate           Code = "DUPLICATE"
	CodeUnauthorized        Code = "UNAUTHORIZED"
	CodeForbidden           Code = "FORBIDDEN"
	CodeNotFound            Code = "NOT_FOUND"
	CodeAlreadyExists       Code = "ALREADY_EXISTS"
	CodeTooManyRequests     Code = "TOO_MANY_REQUESTS"
	CodeInternalServerError Code = "INTERNAL_SERVER_ERROR"
	CodeInvalidTokenError   Code = "INVALID_TOKEN_ERROR"
)

var (
	// 400
	BadRequest = New(http.StatusBadRequest, CodeBadRequest, "Bad Request")

	// 401
	Unauthorized = New(http.StatusUnauthorized, CodeUnauthorized, "Unauthorized")

	// 403
	Forbidden = New(http.StatusForbidden, CodeForbidden, "Forbidden")

	// 404
	NotFound = New(http.StatusNotFound, CodeNotFound, "Not Found")

	// 409
	Duplicate     = New(http.StatusConflict, CodeDuplicate, "Duplicate")
	AlreadyExists = New(http.StatusConflict, CodeAlreadyExists, "Already exists")

	// 429
	TooManyRequests = New(http.StatusTooManyRequests, CodeTooManyRequests, "Too many requests")

	// 498
	InvalidTokenError = New(498, CodeInvalidTokenError, "Invalid token")

	// 500
	InternalServerError = New(http.StatusInternalServerError, CodeInternalServerError, "Internal Server Error")
)
