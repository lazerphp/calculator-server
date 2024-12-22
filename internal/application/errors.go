package application

const (
	InvalidRequestMethodMsg = "неправильный метод запроса"
	InvalidRequestHeaderMsg = "неправильный заголовк запроса"
	InvalidRequestBodyMsg   = "неправильное тело запроса"
)

type RequestError struct {
	Message string
}

func (e *RequestError) Error() string {
	return e.Message
}

func NewError(msg string) *RequestError {
	return &RequestError{msg}
}

var (
	InvalidRequestMethod = NewError(InvalidRequestMethodMsg)
	InvalidRequestHeader = NewError(InvalidRequestHeaderMsg)
	InvalidRequestBody   = NewError(InvalidRequestBodyMsg)
)
