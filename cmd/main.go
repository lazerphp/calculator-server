package main

import "github.com/lazerphp/calculator/internal/application"

func main() {

	app := application.New()
	// app.Run()
	app.RunServer()
}

// func main() {
// ввод из консоли
// in := bufio.NewReader(os.Stdin)
// line, _ := in.ReadString('\n')
// line = strings.TrimSpace(line)

// ввод тествый
// line := "1 + (2 + 3) + 4 + (5 - 6)"
// line := "1 + 2 + 3 + (4 * (5 + 6)) + 7 - 8 + (9 / 10)"
// line := "1 + (2 + 3) * 4 + (-5 + 6) + 7 + (8 + 9)"
// line := "1 + (-2)"
// line := ""
// line := "1 + 2 / 3 * 5"
// line := "10 + (2 + 3) / (5 + 5) + 6"
// line := "1+((1+1)+3) + (6*(5-4))"
// line := "1+(1+1)+3"
// line := "1+((1+1)+3)"

// не должно работать
// line := "25 + 3 * -(188 - 3 * (25 - 1) + 2 * -(10 + 1)) + 6"
// line := "1 + +2"
/*
	test := 1 + -(2)
	fmt.Println(test)
*/

// res, _ := Calc(line)
// fmt.Println(res)

// красивый вывод для тестов
// возможно ли без такого цикла?
// for _, x := range res {
// 	if _, ok := x.(float64); ok {
// 		fmt.Print(x.(float64), " ")
// 	} else {
// 		fmt.Print(string(x.(int32)), " ")
// 	}
// }
// }

/*
достаточно ли функционала в выводах
нормально ли так копировать слайсы
написать тесты
доработать вывод ошибок

функционал с сервером, логгером, файлами, json
тесты для App
поддержка всех возможных обращений
*/
