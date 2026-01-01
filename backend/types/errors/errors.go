package errors

type ErrNotFound struct {
	Msg string
}

func (e ErrNotFound) Error() string {
	return e.Msg
}

func IsNotFound(err error) bool {
	_, ok := err.(ErrNotFound)
	return ok
}

type ErrBadRequest struct {
	Msg string
}

func (e ErrBadRequest) Error() string {
	return e.Msg
}

func IsBadRequest(err error) bool {
	_, ok := err.(ErrBadRequest)
	return ok
}

type ErrFileTypeNotSupported struct {
	Msg string
}

func (e ErrFileTypeNotSupported) Error() string {
	return e.Msg
}

func IsFileTypeNotSupported(err error) bool {
	_, ok := err.(ErrFileTypeNotSupported)
	return ok
}

type ErrFileSizeTooBig struct {
	Msg string
}

func (e ErrFileSizeTooBig) Error() string {
	return e.Msg
}

func IsFileSizeTooBig(err error) bool {
	_, ok := err.(ErrFileSizeTooBig)
	return ok
}

type ErrUnauthorized struct {
	Msg string
}

func (e ErrUnauthorized) Error() string {
	return e.Msg
}

func IsUnauthorized(err error) bool {
	_, ok := err.(ErrUnauthorized)
	return ok
}

type ErrInternalServerError struct {
	Msg string
}

func (e ErrInternalServerError) Error() string {
	return e.Msg
}

func IsInternalServerError(err error) bool {
	_, ok := err.(ErrInternalServerError)
	return ok
}
