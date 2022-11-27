package responses

// Maaping Client Error Responses https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#client_error_responses

func (hrsci *HttpResponseStatusCodesImpl) BadRequest(status string) *HttpResponseStatusCodesImpl {
	return &HttpResponseStatusCodesImpl{
		Code:   400,
		Status: SetStatus(status, "BAD_REQUEST"),
	}
}

func (hrsci *HttpResponseStatusCodesImpl) Unathorized(status string) *HttpResponseStatusCodesImpl {
	return &HttpResponseStatusCodesImpl{
		Code:   401,
		Status: SetStatus(status, "UNAUTHORIZED"),
	}
}

func (hrsci *HttpResponseStatusCodesImpl) Forbidden(status string) *HttpResponseStatusCodesImpl {
	return &HttpResponseStatusCodesImpl{
		Code:   403,
		Status: SetStatus(status, "FORBIDDEN"),
	}
}

func (hrsci *HttpResponseStatusCodesImpl) NotFound(status string) *HttpResponseStatusCodesImpl {
	return &HttpResponseStatusCodesImpl{
		Code:   404,
		Status: SetStatus(status, "NOT_FOUND"),
	}
}

func (hrsci *HttpResponseStatusCodesImpl) Conflict(status string) *HttpResponseStatusCodesImpl {
	return &HttpResponseStatusCodesImpl{
		Code:   409,
		Status: SetStatus(status, "CONFLICT"),
	}
}

func (hrsci *HttpResponseStatusCodesImpl) UnprocessableEntity(status string) *HttpResponseStatusCodesImpl {
	return &HttpResponseStatusCodesImpl{
		Code:   422,
		Status: SetStatus(status, "UNPROCESSABLE_ENTITY"),
	}
}
