package Utilities

type ValidatorResult struct {
	IsSuccess bool
	Errors    []string
}

func Success() *ValidatorResult {
	return &ValidatorResult{
		IsSuccess: true,
	}
}

func Fail(errors []string) *ValidatorResult {
	return &ValidatorResult{
		IsSuccess: false,
		Errors:    errors,
	}
}
