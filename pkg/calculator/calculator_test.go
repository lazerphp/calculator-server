package calculator_test

import (
	"testing"

	"github.com/lazerphp/calculator/pkg/calculator"
)

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "1 число",
			expression:     "12",
			expectedResult: 12,
		},

		{
			name:           "все операторы",
			expression:     "1 + 2 - 3 * 4 / 6",
			expectedResult: 1,
		},
		{
			name:           "различные числа",
			expression:     "10 + 200 - 3000 * 45678 / 678 + 0.7 - 0.89",
			expectedResult: -201905.23424778762,
		},
		{
			name:           "отрицательные числа",
			expression:     "-1 + 3",
			expectedResult: 2,
		},
		{
			name:           "скобки",
			expression:     "1+(1+1)+3",
			expectedResult: 6,
		},
		{
			name:           "отрицательные числа в скобках",
			expression:     "1 + (2 + 3) * 4 + (-5 + 6) + 7 + (8 + 9)",
			expectedResult: 46,
		},
		{
			name:           "несколько скобок",
			expression:     "1 + (2 + 3) + 4 + (5 - 6)",
			expectedResult: 9,
		},
		{
			name:           "вложенные скобки",
			expression:     "25 + 3 * -(188 - 3 * (25 - 1) + 2 * -(10 + 1)) + 6",
			expectedResult: -251,
		},
		{
			name:           "комбинации *- и /-",
			expression:     "1/-2",
			expectedResult: -0.5,
		},
		{
			name:           "- перед скобками",
			expression:     "1/-(2)",
			expectedResult: -0.5,
		},
		{
			name:           "скобки в начале",
			expression:     "(1+((1+1)+3)) + (6*(5-4)) + 10",
			expectedResult: 22,
		},
		{
			name:           "скобки в конце",
			expression:     "1+((1+1)+3)",
			expectedResult: 6,
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calculator.Calc(testCase.expression)
			if err != nil {
				t.Fatalf("successful case %s returns error", testCase.expression)
			}
			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}

	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:       "пустой запрос",
			expression: "",
		},
		{
			name:       "в запросе недопустимые лишние символы",
			expression: "1,2+3,4",
		},
		{
			name:       "в запросе нет цифр",
			expression: "+--+*//()-*",
		},
		{
			name:       "лишняя закрывающая скобка",
			expression: "(1 + 2 + 3) + 4) + 5",
		},
		{
			name:       "пустые скобки",
			expression: "1 + () + 2",
		},
		{
			name:       "нарушен порядок скобок",
			expression: "(1 + 2 + (3 + 4 +) + 5)",
		},
		{
			name:       "оператор в начале выражения",
			expression: "+ 1 - 2",
		},
		{
			name:       "оператор в конце выражения",
			expression: "1 - 2 -",
		},
		{
			name:       "оператор в начале скобок",
			expression: "1 + (+2)",
		},
		{
			name:       "оператор в конце скобок",
			expression: "1 + (2-)",
		},
		{
			name:       "ввод +-",
			expression: "1 + -(2 + 3)",
		},
		{
			name:       "цепочка из минусов",
			expression: "1*--2",
		},
		{
			name:       "точки не на своих местах",
			expression: "1. + 2",
		},
		{
			name:       "деление на 0",
			expression: "1 / 0",
		},
		{
			name:       "точка посреди пробелов",
			expression: "1 . 2 + 3   . 4",
		},
		{
			name:       "пробел посреди числа",
			expression: "1 2 + 3 4 - 5 6",
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calculator.Calc(testCase.expression)
			if err == nil {
				t.Fatalf("expression %s is invalid but result  %f was obtained", testCase.expression, val)
			}
		})
	}
}
