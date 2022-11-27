package exception

type ErrorResponsesImpl struct {
	err     error
	code    int
	status  string
	data    interface{}
	message string
}

func (eri *ErrorResponsesImpl) NewErrorResponses(err error, data any, message string) *ErrorResponsesImpl {
	return &ErrorResponsesImpl{
		err:     err,
		data:    data,
		code:    eri.code,
		status:  eri.status,
		message: message,
	}
}
