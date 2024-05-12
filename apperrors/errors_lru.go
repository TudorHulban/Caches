package apperrors

const msgErrRecordNotFound = "record not found"

type ErrRecordNotFound struct{}

func (ErrRecordNotFound) Error() string {
	return msgErrRecordNotFound
}

func (e ErrRecordNotFound) Is(err error) bool {
	_, couldCast := err.(ErrRecordNotFound)

	return couldCast
}
