package main

//go:generate stringer -type=StatusCode
type StatusCode int

const (
	OK                    StatusCode = 200
	CREATED               StatusCode = 201
	ACCEPTED              StatusCode = 202
	NO_CONTENT            StatusCode = 204
	FOUND                 StatusCode = 302
	BAD_REQUEST           StatusCode = 400
	UNAUTHORIZED          StatusCode = 401
	FORBIDDEN             StatusCode = 403
	NOT_FOUND             StatusCode = 404
	METHOD_NOT_ALLOWED    StatusCode = 405
	CONTENT_TOO_LARGE     StatusCode = 413
	URI_TOO_LONG          StatusCode = 414
	INTERNAL_SERVER_ERROR StatusCode = 500
)
