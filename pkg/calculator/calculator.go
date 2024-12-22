package calculator

import (
	"regexp"
	"strconv"
	"strings"
)

// валидация с пробелами
func checkSpaceInNums(str string) error {
	var prevDigitPos int
	var prevNoDigitPos int
	for i, x := range str {
		if x == ' ' {
			continue
		}

		if strings.Contains("1234567890.", string(x)) {
			if i > 0 && prevDigitPos > prevNoDigitPos && prevDigitPos+1 != i {
				return InvalidInput.With("пробелы внутри числа")
			}
			prevDigitPos = i
		} else {
			prevNoDigitPos = i
		}
	}
	return nil
}

// более полная валидация
func checkInput(str string) error {

	// валидация на смысл
	if len(str) == 0 {
		return InvalidInput.With("пустой ввод")
	}
	validPattern := regexp.MustCompile(`^[0-9+\-*/\(\).]*$`)
	if !validPattern.MatchString(str) {
		return InvalidInput.With("невалидные символы")
	}
	if !strings.ContainsAny(str, "1234567890") {
		return InvalidInput.With("в выражении нет цифр")
	}

	// валидация скобок
	cntBrackets := 0
	for i, x := range str {
		if x == ')' {
			cntBrackets--
		}
		if x == '(' {
			cntBrackets++
		}
		if cntBrackets < 0 {
			return InvalidBrackets.With("лишняя закрывающая скобка")
		}
		if x == ')' && i > 0 && str[i-1] == '(' {
			return InvalidBrackets.With("замечены пустые скобки")
		}
	}
	if cntBrackets > 0 {
		return InvalidBrackets.With("не хватает закрывающих скобок")
	}

	// остальные случаи валидации
	for i, x := range str {

		// обработка знаков, кроме -
		if strings.Contains("+*/", string(x)) {
			if i == 0 {
				return InvalidOperators.With("'%c' в начале выражения", x)
			}
			if i == len(str)-1 {
				return InvalidOperators.With("'%c' в конце выражения", x)
			}
			if i > 0 && strings.Contains("+-*/", string(str[i-1])) {
				return InvalidOperators.With("подряд идущие знаки")
			}
			if i > 0 && str[i-1] == '(' {
				return InvalidOperators.With("'%c' стоит сразу после ( на позиции %v", x, i+1)
			}
			if i < len(str)-1 && str[i+1] == ')' {
				return InvalidOperators.With("'%c' стоит перед ) на позиции %v", x, i+1)
			}
		}

		// обработка знака -
		if x == '-' {
			if i == len(str)-1 {
				return InvalidOperators.With("минус в конце выражения")
			}
			if i < len(str)-1 && str[i+1] == ')' {
				return InvalidOperators.With("минус в конце скобок")
			}
			if i > 0 && str[i-1] == '+' {
				return InvalidOperators.With("недопустимый '+-'")
			}
			if i > 0 && str[i-1] == '-' {
				return InvalidOperators.With("недопустимый '--'")
			}
			// это уже не нужно
			// if i > 1 && str[i-1] == '-' && strings.Contains("*/", string(str[i-2])) {
			// 	return InvalidOperators.With("знак и два минуса подряд")
			// }
		}

		// обработка точки
		if x == '.' {
			if i == 0 {
				return InvalidOperators.With("точка в начале")
			}
			if i == len(str)-1 {
				return InvalidOperators.With("точка в конце")
			}
			if i > 0 && i < len(str)-1 && (!strings.Contains("1234567890", string(str[i-1])) || !strings.Contains("1234567890", string(str[i+1]))) {
				return InvalidOperators.With("точка стоит не на своём месте")
			}
		}
	}

	return nil
}

// подготовка массива
func prepareExp(str string) []interface{} {
	var res []interface{}
	var buf []rune

	for i, x := range str {
		if strings.Contains("1234567890.", string(x)) {
			buf = append(buf, x)
		} else if x == '-' &&
			i > 0 && i < len(str)-1 &&
			strings.Contains("*/", string(str[i-1])) && str[i+1] != '(' {
			buf = append(buf, x)
		} else {
			if len(buf) != 0 {
				bufNum, _ := strconv.ParseFloat(string(buf), 64)
				res = append(res, bufNum)
			}
			buf = nil
			res = append(res, x)
		}

		if i == len(str)-1 && len(buf) != 0 {
			bufNum, _ := strconv.ParseFloat(string(buf), 64)
			res = append(res, bufNum)
		}
	}

	return res
}

// поиск скобочного выражения
func findBrackets(arr []interface{}) (int, int, bool) {
	openPos := -1
	closePos := -1
	cnt := 0

	for i := range len(arr) {
		if arr[i] == '(' && openPos == -1 {
			openPos = i
		}

		if arr[i] == '(' {
			cnt++
		} else if arr[i] == ')' {
			cnt--
		}

		if cnt == 0 && openPos != -1 {
			closePos = i
			return openPos, closePos, true
		}
	}
	return 0, 0, false
}

// функция удаления из середины
func remove(slice []interface{}, i int) []interface{} {
	return append(slice[:i], slice[i+1:]...)
}

// выполнение операторов согласно приоритету
func execOps(arr []interface{}) (float64, error) {

	i := 0
	for i < len(arr)-1 && arr[i] != ')' {
		if arr[i] == '-' && i > 0 && (arr[i-1] == '*' || arr[i-1] == '/') {

			arr[i+1] = -arr[i+1].(float64)
			arr = remove(arr, i)
		} else if arr[i] == '-' && i == 0 {
			arr[i+1] = -arr[i+1].(float64)
			arr = remove(arr, i)
		} else {
			i++
		}
	}

	i = 1
	for i < (len(arr)-1) && arr[i] != ')' {
		if arr[i] == '*' || arr[i] == '/' {
			var buf float64
			if arr[i] == '*' {
				buf = arr[i-1].(float64) * arr[i+1].(float64)
			} else if arr[i] == '/' {
				if arr[i+1].(float64) == 0 {
					return 0, DivizionByZero
				}
				buf = arr[i-1].(float64) / arr[i+1].(float64)
			}
			arr[i-1] = buf
			arr = remove(arr, i)
			arr = remove(arr, i)
		} else {
			i++
		}
	}

	i = 1
	for i < (len(arr)-1) && arr[i] != ')' {
		if arr[i] == '+' || arr[i] == '-' {
			var buf float64
			if arr[i] == '+' {
				buf = arr[i-1].(float64) + arr[i+1].(float64)
			} else if arr[i] == '-' {
				buf = arr[i-1].(float64) - arr[i+1].(float64)
			}
			arr[i-1] = buf
			arr = remove(arr, i)
			arr = remove(arr, i)
		}
	}
	return arr[0].(float64), nil
}

// рекурсивный (из-за скобочных выражений) подсчет
func recCalc(arr []interface{}) (float64, error) {
	for range arr {
		openBracket, closeBracket, hasBrackets := findBrackets(arr)
		if hasBrackets {
			// подход
			copiedArr := make([]interface{}, closeBracket-openBracket-1)
			copy(copiedArr, arr[openBracket+1:closeBracket])

			res, err := recCalc(copiedArr)
			if err != nil {
				return 0, err
			}
			tail := make([]interface{}, len(arr)-closeBracket-1)
			copy(tail, arr[closeBracket+1:])
			arr = append(arr[:openBracket], res)
			arr = append(arr, tail...)
		}
	}
	calculated, err := execOps(arr)
	if err != nil {
		return 0, err
	}
	return calculated, nil
}

// кальк
func Calc(input string) (float64, error) {
	spaceErr := checkSpaceInNums(input)
	if spaceErr != nil {
		return 0, spaceErr
	}

	preparedStr := strings.ReplaceAll(input, " ", "")
	inputErr := checkInput(strings.ReplaceAll(input, " ", ""))
	if inputErr != nil {
		return 0, inputErr
	}

	preparedExp := prepareExp(preparedStr)
	// fmt.Println(preparedExp)
	res, calcErr := recCalc(preparedExp)
	if calcErr != nil {
		return 0, calcErr
	}

	return res, nil
}
