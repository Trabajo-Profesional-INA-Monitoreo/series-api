package exceptions

type NotFound struct {
}

func NewNotFound() *NotFound {
	return &NotFound{}
}

func (m *NotFound) Error() string {
	return "404 - "
}

type BadRequest struct {
}

func NewBadRequest() *BadRequest {
	return &BadRequest{}
}

func (m *BadRequest) Error() string {
	return "400 - "
}

type InternalError struct {
}

func NewInternalError() *InternalError {
	return &InternalError{}
}

func (m *InternalError) Error() string {
	return "500 - "
}

func MapCodeToError(code int) error {
	switch code {
	case 404:
		return NewNotFound()
	case 400:
		return NewBadRequest()
	default:
		return NewInternalError()
	}
}
