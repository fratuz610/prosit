package web

import (
	"github.com/martini-contrib/render"
	"log"
	"prosit/cerr"
)

type JSONError struct {
	HttpCode     int
	ErrorCode    string
	ErrorMessage string
}

func NewJSONError(httpCode int, err error) *JSONError {

	errorCode := ""

	switch {
	case httpCode == 500:
		errorCode = "SERVER_ERROR"
	case httpCode == 400:
		errorCode = "BAD_REQUEST"
	case httpCode == 403:
		errorCode = "UNAUTHORIZED"
	}
	return &JSONError{HttpCode: httpCode, ErrorCode: errorCode, ErrorMessage: err.Error()}
}

func outputError(outErr error, r render.Render) {

	var httpCode = 0

	switch outErr.(type) {
	case *cerr.BadRequestError:
		httpCode = 400
	case *cerr.NotFoundError:
		httpCode = 404
	default:
		httpCode = 500
	}

	// we output the error
	r.JSON(httpCode, NewJSONError(httpCode, outErr))

	// we also log it
	log.Printf("--- HTTP ERROR %d: '%v'", httpCode, outErr)
}
