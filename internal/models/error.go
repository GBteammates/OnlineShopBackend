package models

type ErrorNotFound struct {
}

func (e ErrorNotFound) Error() string {
	return ""
}

type ErrNotExist struct{}

func (e ErrNotExist) Error() string {
	return "cache is not exists"
}
