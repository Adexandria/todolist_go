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

func NotFoundDataResult[T any](err ...string) *ActionResultModel[T] {
	return &ActionResultModel[T]{
		ActionResult: ActionResult{
			StatusCode: http.StatusNotFound,
			IsSuccess:  false,
			Error:      err,
		},
	}
}

func NotAuthorizedResult(errors ...string) *ActionResult {
	return &ActionResult{
		StatusCode: http.StatusUnauthorized,
		IsSuccess:  false,
		Error:      errors,
	}
}

func ForbiddenResult(errors ...string) *ActionResult {
	return &ActionResult{
		StatusCode: http.StatusForbidden,
		IsSuccess:  false,
		Error:      errors,
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

func BadRequestDataResult[T any](err ...string) *ActionResultModel[T] {
	return &ActionResultModel[T]{
		ActionResult: ActionResult{
			StatusCode: http.StatusBadRequest,
			IsSuccess:  false,
			Error:      err,
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

func ServerErrorResult(err ...string) *ActionResult {
	return &ActionResult{
		StatusCode: http.StatusInternalServerError,
		IsSuccess:  false,
		Error:      err,
	}
}

func ServerErrorDataResult[T any](err ...string) *ActionResultModel[T] {
	return &ActionResultModel[T]{
		ActionResult: ActionResult{
			StatusCode: http.StatusInternalServerError,
			IsSuccess:  false,
			Error:      err,
		},
	}
}
