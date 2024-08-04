package entities

type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "Entity not found"
}

type ErrAlreadyExists struct{}

func (e *ErrAlreadyExists) Error() string {
	return "Entity already exists"
}

type ErrBadCredentials struct{}

func (e *ErrBadCredentials) Error() string {
	return "Bad credentials"
}
