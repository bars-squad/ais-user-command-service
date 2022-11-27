package responses

// Mapping Successful Responses https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#successful_responses

func (hrsci *HttpResponseStatusCodesImpl) Ok(status string) *HttpResponseStatusCodesImpl {
	return &HttpResponseStatusCodesImpl{
		Code:   200,
		Status: SetStatus(status, "OK"),
	}
}

func (hrsci *HttpResponseStatusCodesImpl) Created(status string) *HttpResponseStatusCodesImpl {
	return &HttpResponseStatusCodesImpl{
		Code:   201,
		Status: SetStatus(status, "CREATED"),
	}
}
