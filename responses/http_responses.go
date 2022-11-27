package responses

// You can say this file is a main of responses folder

// HTTP response status codes indicate whether a specific HTTP request has been successfully completed. Responses are grouped in five classes. Please read these article https://developer.mozilla.org/en-US/docs/Web/HTTP/Status

type HttpResponseStatusCodes interface {
	Ok(status string) *HttpResponseStatusCodesImpl
	Created(status string) *HttpResponseStatusCodesImpl
	BadRequest(status string) *HttpResponseStatusCodesImpl
	Unauthorized(status string) *HttpResponseStatusCodesImpl
	Forbidden(status string) *HttpResponseStatusCodesImpl
	NotFound(status string) *HttpResponseStatusCodesImpl
	Conflict(status string) *HttpResponseStatusCodesImpl
	UnprocessableEntity(status string) *HttpResponseStatusCodesImpl
}

type HttpResponseStatusCodesImpl struct {
	Code   int
	Status string
}

func SetStatus(status string, defaultStatus string) string {
	if status != "" {
		return status
	}
	return defaultStatus
}
