package guard

import "fmt"

func NewNotFoundClientError(entity string) error {
	return NewClientError(fmt.Errorf("%s not found", entity))
}

func NewAlreadyExistsClientError(entity string) error {
	return NewClientError(fmt.Errorf("%s already exists", entity))
}
