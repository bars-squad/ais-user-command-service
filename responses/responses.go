package responses

type Responses interface {
	DataProperty() interface{}
	StatusProperty() string
	CodeProperty() int
	MessageProperty() string
	PaginationProperty() interface{}
}

// We can use this internal response to wrap response
type ResponsesImpl struct {
	Code       int         `json:"code"`
	Status     string      `json:"status"`
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	Pagination interface{} `json:"pagination,omitempty"`
}

func (hrsci *HttpResponseStatusCodesImpl) NewResponses(data any, message string) Responses {
	return &ResponsesImpl{
		Data:    data,
		Code:    hrsci.Code,
		Status:  hrsci.Status,
		Message: message,
	}
}

// DataProperty returns data.
func (r *ResponsesImpl) DataProperty() interface{} {
	return r.Data
}

// StatusProperty returns HTTP status.
func (r *ResponsesImpl) StatusProperty() string {
	return r.Status
}

// CodeProperty returns http code.
func (r *ResponsesImpl) CodeProperty() int {
	return r.Code
}

// MessageProperty returns message.
func (r *ResponsesImpl) MessageProperty() string {
	return r.Message
}

// PaginationProperty returns pagination.
func (r *ResponsesImpl) PaginationProperty() interface{} {
	return r.Pagination
}
