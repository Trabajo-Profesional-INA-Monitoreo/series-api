package exceptions

type NotFound struct {
}

func NewNotFound() *NotFound {
	return &NotFound{}
}

func (m *NotFound) Error() string {
	return "404 - "
}
