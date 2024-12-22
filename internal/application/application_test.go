package application

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// заголовки запроса
func TestHandlerRequestHeaders(t *testing.T) {

	// запрос, отличный от POST
	name := "отправка запроса, отличного от POST"
	request, _ := http.NewRequest("GET", "/api/v1/calculate", nil)
	response := httptest.NewRecorder()
	handler := Outer(Validation(CalcHandler))
	handler(response, request)
	res := response.Result()

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("%v\n неожиданный код ответа %v", name, res.StatusCode)
	}

	// запрос без нужного Content-Type
	name = "отправка запроса без Content-Type: application/json"
	preparedBody := strings.NewReader(`{"expression": "1+2-3"}`)
	request, _ = http.NewRequest("POST", "/api/v1/calculate", preparedBody)
	response = httptest.NewRecorder()
	handler = Outer(Validation(CalcHandler))
	handler(response, request)
	res = response.Result()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("%v\n неожиданный код ответа %v", name, res.StatusCode)
	}
}

// тело запроса
func TestHandlerRequestBody(t *testing.T) {

	tt := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{"Valid", `{"expression": "1+2-3"}`, http.StatusOK},
		{"Valid with other fields", `{"expression": "1+2-3", "hello": "there"}`, http.StatusOK},
		{"Invalid in calculator", `{"expression": "(1+2-3"}`, http.StatusUnprocessableEntity},
		{"No key 'expression'", `{"hell": "0"}`, http.StatusUnprocessableEntity},
		{"Empty'", "", http.StatusUnprocessableEntity},
		{"JSON array'", `["люди","бананы"]`, http.StatusUnprocessableEntity},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			preparedBody := strings.NewReader(tc.body)
			request, _ := http.NewRequest("POST", "/api/v1/calculate", preparedBody)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			handler := Outer(Validation(CalcHandler))

			handler(response, request)
			res := response.Result()
			if res.StatusCode != tc.expectedStatus {
				// data, err := io.ReadAll(res.Body)
				// if err != nil {
				// 	fmt.Println(err)
				// }
				defer res.Body.Close()
				t.Errorf("неожиданный код ответа %v", res.StatusCode)
			}
		})
	}
}
