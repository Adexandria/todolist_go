package services

import "net/http"

type ActionResultModel[T any] struct {
	Data T
	ActionResult
}

type ActionResult struct {
	StatusCode int
	IsSuccess  bool
	Error      []string
}

func NotFoundResult(error ...string) *ActionResult {
	return &ActionResult{
		StatusCode: http.StatusNotFound,
		IsSuccess:  false,
		Error:      error,
	}
}

func SuccessResult() *ActionResult {
	return &ActionResult{
		StatusCode: http.StatusOK,
		IsSuccess:  true,
	}
}

func SuccessDataResult[T any](data T) *ActionResultModel[T] {
	return &ActionResultModel[T]{
		Data: data,
		ActionResult: ActionResult{
			StatusCode: http.StatusOK,
			IsSuccess:  true,
		},
	}
}

func BadRequestResult(err ...string) *ActionResult {
	return &ActionResult{
		StatusCode: http.StatusBadRequest,
		IsSuccess:  false,
		Error:      err,
	}
}
