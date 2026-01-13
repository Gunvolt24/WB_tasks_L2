package domain

// ErrorWithCode позволяет каждой ошибке возвращать HTTP код
type ErrorWithCode interface {
	error
	Code() int
}

// ValidationError - возвращает 400
type ValidationError struct{ Msg string }

func (e ValidationError) Error() string { return e.Msg }
func (e ValidationError) Code() int     { return 400 }

// BusinessError - возвращает 503
type BusinessError struct{ Msg string }

func (e BusinessError) Error() string { return e.Msg }
func (e BusinessError) Code() int     { return 503 }

// NotFoundError - возвращает 503
type NotFoundError struct{ Msg string }

// Имплементация интерфейса ErrorWithCode
func (e NotFoundError) Error() string { return e.Msg }
func (e NotFoundError) Code() int     { return 503 }

// Конструкторы
func NewValidationError(msg string) error { return ValidationError{Msg: msg} }
func NewBusinessError(msg string) error   { return BusinessError{Msg: msg} }
func NewNotFoundError(msg string) error   { return NotFoundError{Msg: msg} }
