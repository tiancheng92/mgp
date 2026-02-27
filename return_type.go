package mgp

import "net/http"

type ReturnType struct {
	StatusCode int
	Body       any
}

func RT[T any](statusCode ...int) *ReturnType {
	hc := http.StatusOK
	if len(statusCode) != 0 {
		hc = statusCode[0]
	}

	return &ReturnType{
		StatusCode: hc,
		Body:       new(Result[T]),
	}
}

func RTE(statusCode ...int) *ReturnType {
	hc := http.StatusOK
	if len(statusCode) != 0 {
		hc = statusCode[0]
	}

	return &ReturnType{
		StatusCode: hc,
	}
}

func PRT[T any](statusCode ...int) *ReturnType {
	hc := http.StatusOK
	if len(statusCode) != 0 {
		hc = statusCode[0]
	}

	return &ReturnType{
		StatusCode: hc,
		Body:       new(Result[PaginateData[[]T]]),
	}
}
