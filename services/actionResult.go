package services

type ActionResultModel[T any] struct {
	Data T
	ActionResult
}

type ActionResult struct {
	StatusCode int
	IsSuccess  bool
	Error      []string
}
