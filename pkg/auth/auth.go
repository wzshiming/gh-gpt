package auth

import (
	"errors"
	"os"
)

var ErrNotExists = os.ErrNotExist

type Auth interface {
	GetToken() (string, error)
}

type Auths []Auth

func (a Auths) GetToken() (string, error) {
	for _, auth := range a {
		token, err := auth.GetToken()
		if err == nil {
			return token, nil
		}
		if !errors.Is(err, ErrNotExists) {
			return "", err
		}
	}
	return "", ErrNotExists
}
