package utils

import "github.com/codecrafters-io/redis-starter-go/app/globals"

func ErrorResponse(request *globals.Request, error_code globals.ErrorCode) *globals.Response {
	return &globals.Response{
		Request: request,
		Body:    error_code.Bytes(),
	}
}

func SuccessResponse(request *globals.Request, body []byte) *globals.Response {
	return &globals.Response{
		Request: request,
		Body:    append(globals.Success.Bytes(), body...),
	}
}
