package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lazerphp/calculator/pkg/calculator"
)

func CalcHandler(w http.ResponseWriter, r *http.Request) {

	// получение expression
	request := &Request{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		panic(InvalidRequestBody)
	}

	// вызов калькулятора
	result, err := calculator.Calc(request.Expression)
	if err != nil {
		panic(err)
	}

	// подготовка ответа
	response := &ResponseSuccess{result}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	// ответ
	fmt.Fprint(w, string(jsonResponse))
}

func Validation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// перехват ошибок из валидации калькулятора
		defer func() {
			if err := recover(); err != nil {
				expectedErrType := false
				var expectedErr error
				_, expectedErrTypeRequest := err.(*RequestError)
				_, expectedErrTypeCalc := err.(*calculator.ErrorDetailed)

				switch {
				case expectedErrTypeRequest:
					expectedErr = err.(*RequestError)
					expectedErrType = true
				case expectedErrTypeCalc:
					expectedErr = err.(*calculator.ErrorDetailed)
					expectedErrType = true
				}

				if expectedErrType {

					GenerateErr(expectedErr.Error(), http.StatusUnprocessableEntity, &w)
				} else {
					panic(err)

				}
			}
		}()

		// валидация самого запроса
		if r.Method != http.MethodPost {
			GenerateErr("неверный метод запроса", http.StatusUnprocessableEntity, &w)

		} else if r.Header.Get("Content-Type") != "application/json" {
			GenerateErr("неверный формат запроса", http.StatusUnprocessableEntity, &w)

		} else {
			next(w, r)
		}
	}
}

func Outer(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// перехват внутренних ошибок
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", "неизвестная ошибка на стороне сервера")
				// log.Printf("Panic: %v", err)
				GenerateErr("Internal server Error", http.StatusInternalServerError, &w)
			}
		}()
		next(w, r)
	}
}
