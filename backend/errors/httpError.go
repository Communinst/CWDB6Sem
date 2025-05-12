package customErrors

type ErrorWithStatusCode struct {
	HTTPStatus int
	Msg        string
}

func (e *ErrorWithStatusCode) Error() string {
	return e.Msg
}
