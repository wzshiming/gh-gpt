package auth

import (
	"context"
	"errors"
	"os"
)

var ErrNotExists = os.ErrNotExist

type Auth interface {
	GetToken(ctx context.Context) (string, error)
}

type Auths []Auth

func (a Auths) GetToken(ctx context.Context) (string, error) {
	var errs []error
	for _, auth := range a {
		token, err := auth.GetToken(ctx)
		if err == nil {
			if token != "" {
				return token, nil
			}
			errs = append(errs, ErrNotExists)
		} else {
			errs = append(errs, err)
		}
	}
	return "", errors.Join(errs...)
}
