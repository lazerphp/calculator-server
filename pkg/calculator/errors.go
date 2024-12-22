package calculator

import "fmt"

const (
	InvalidInputMsg     = "некорректный ввод"
	InvalidBracketsMsg  = "некорректные скобки"
	InvalidOperatorsMsg = "некорректная расстановка знаков"
	DivizionByZeroMsg   = "деление на 0"
)

type ErrorDetailed struct {
	Message string
	Details string
}

func (e *ErrorDetailed) Error() string {
	return e.Message + " -- " + e.Details
}
func (e *ErrorDetailed) With(details string, variables ...interface{}) error {
	// или через указатели
	// проверить этот метод
	if len(variables) != 0 {
		e.Details = fmt.Sprintf(details, variables...)
	} else {
		e.Details = details
	}
	return e
}

func NewError(msg string) *ErrorDetailed {
	return &ErrorDetailed{msg, ""}
}

var (
	InvalidInput     = NewError(InvalidInputMsg)
	InvalidOperators = NewError(InvalidOperatorsMsg)
	InvalidBrackets  = NewError(InvalidBracketsMsg)
	DivizionByZero   = NewError(DivizionByZeroMsg)
)
