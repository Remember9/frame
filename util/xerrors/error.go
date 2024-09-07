package xerrors

const (
	UnknownErrorCode = 1
)

func New(code int64, text string) error {
	return &CodeError{code, text}
}

type CodeError struct {
	code int64
	s    string
}

func (e *CodeError) Error() string {
	return e.s
}

func (e *CodeError) Code() int64 {
	return e.code
}

func ToCodeError(err error) *CodeError {
	if err == nil {
		return nil
	}

	errorCode, ok := err.(*CodeError)
	if ok {
		return errorCode
	}

	return New(UnknownErrorCode, err.Error()).(*CodeError)
}
