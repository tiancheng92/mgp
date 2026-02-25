package errors

import (
	"net/http"
	"sync"
)

var (
	codes sync.Map
	keys  sync.Map
)

type Coder interface {
	HTTPStatus() int
	String() string
	Code() int
}

type errCode struct {
	ErrCode    int
	HttpStatus int
	Message    string
}

func (e *errCode) Code() int {
	return e.ErrCode
}

func (e *errCode) String() string {
	return e.Message
}

func (e *errCode) HTTPStatus() int {
	if e.HttpStatus == 0 {
		return http.StatusInternalServerError
	}
	return e.HttpStatus
}

func Register(key string, code int, httpStatus int, message string) {
	if code == 0 {
		panic("code '0' is ErrUnknown error code")
	}
	if _, ok := keys.Load(key); ok {
		keys.Delete(key)
	}
	keys.Store(key, code)

	if _, ok := codes.Load(code); ok {
		codes.Delete(code)
	}
	codes.Store(code, &errCode{
		ErrCode:    code,
		HttpStatus: httpStatus,
		Message:    message,
	})
}

func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}

	var withCodeErr *withCode

	if As(err, &withCodeErr) {
		if res, ok := codes.Load(withCodeErr.code); ok {
			return res.(Coder)
		}
	}
	return &errCode{
		ErrCode:    500000,
		HttpStatus: http.StatusInternalServerError,
		Message:    "error code is undefined",
	}
}
