package application

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	Expression string `json:"expression"`
}

type ResponseSuccess struct {
	Result float64 `json:"result"`
}
type ResponseError struct {
	Result string `json:"error"`
}

func GenerateErr(msg string, code int, w *http.ResponseWriter) {
	errResponse := &ResponseError{msg}
	jsonedMsg, _ := json.Marshal(errResponse)
	http.Error(*w, string(jsonedMsg), code)
}
