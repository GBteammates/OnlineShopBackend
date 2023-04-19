package models

type ErrorNotFound struct {

}

func (e ErrorNotFound) Error() string {
	return ""
}
