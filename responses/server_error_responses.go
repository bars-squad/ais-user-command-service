package responses

// Mapping Server Error Responses https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#server_error_responses

func (hrsci *HttpResponseStatusCodesImpl) InternalServerError(status string) *HttpResponseStatusCodesImpl {
	return &HttpResponseStatusCodesImpl{
		Code:   500,
		Status: SetStatus(status, "INTERNAL_SERVER_ERROR"),
	}
}
